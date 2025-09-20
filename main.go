package main

import (
	"discord-lookup/internal/discord"
	api "discord-lookup/internal/http"
	"discord-lookup/internal/types"
	"encoding/json"
	"log"
	"net/http"
	"os"

	logger "github.com/yourpov/logrite"
)

func main() {
	logger.SetConfig(logger.Config{
		ShowIcons:    true,
		UppercaseTag: true,
		UseColors:    true,
	})
	cfg, err := loadConfig("config/config.json")
	if err != nil {
		logger.Warn("config: %v", err)
	}
	if cfg.Token == "" {
		logger.Warn("missing bot token in config")
	}

	disc := discord.New(cfg.Token)
	api := &api.Server{Discord: disc}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("public")))
	api.Routes(mux)

	logger.Success("Listening on port %s", cfg.Port)
	logger.Info("API Endpoint: http://localhost:%s/lookup?id={id}", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}

func loadConfig(path string) (types.Config, error) {
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
		logger.Error("Port is required in config")
	}
	return cfg, nil
}
