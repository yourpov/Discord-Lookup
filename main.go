package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	logger "github.com/yourpov/logrite"

	"discord-lookup/internal/discord"
	api "discord-lookup/internal/http"
	"discord-lookup/internal/types"
)

func main() {
	logger.SetConfig(logger.Config{
		ShowIcons:    true,
		UppercaseTag: true,
		UseColors:    true,
	})
	cfg, err := loadCfg("config/config.json")
	if err != nil {
		logger.Warn("config: %v", err)
	}
	if cfg.Token == "" {
		logger.Warn("no bot token found")
	}

	disc := discord.New(cfg.Token)
	srv := &api.Server{Discord: disc}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("public")))
	srv.Routes(mux)

	logger.Success("running on port %s", cfg.Port)
	logger.Info("endpoint: http://localhost:%s/lookup?id={id}", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}

func loadCfg(path string) (types.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return types.Config{}, err
	}
	defer f.Close()

	var cfg types.Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return types.Config{}, err
	}
	if cfg.Port == "" {
		logger.Error("port required")
	}
	return cfg, nil
}
