package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"discord-lookup/internal/discord"
	"discord-lookup/internal/types"
)

type Server struct {
	Discord *discord.Client
}

// Routes sets up the HTTP routes for the server.
func (s *Server) Routes(mux *http.ServeMux) {
	mux.Handle("/lookup", s.cors(http.HandlerFunc(s.lookup)))
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
}

// lookup handles user lookup requests.
func (s *Server) lookup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	u, code, err := s.Discord.FetchUser(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	resp := types.User{
		ID:            u.ID,
		Username:      u.Username,
		DisplayName:   u.GlobalName,
		Discriminator: u.Discriminator,
		Bot:           u.Bot,
		System:        u.System,
		Flags:         u.PublicFlags,
		Badges:        discord.DecodeBadges(u.PublicFlags),
		Avatar:        discord.Avatar(u),
		Banner:        discord.Banner(u),
		CreatedAt:     discord.CreatedAt(u.ID),
		SearchedAt:    types.JSONTime(time.Now().UTC()),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// cors is a middleware that adds CORS headers to responses.
func (s *Server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
