package discord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"discord-lookup/internal/types"
)

const (
	apiBase = "https://discord.com/api/v10"
	epoch   = int64(1420070400000)
	maxBody = 1 << 20
	timeout = 10 * time.Second
	uagent  = "discord-lookup (https://github.com/yourpov/discord-lookup, 1.0)"
)

type Client struct {
	http  *http.Client
	token string
}

// New creates a new Discord API client.
func New(token string) *Client {
	return &Client{
		http:  &http.Client{Timeout: timeout},
		token: token,
	}
}

// FetchUser retrieves a Discord user by ID.
func (c *Client) FetchUser(ctx context.Context, id string) (types.RawUser, int, error) {
	if c.token == "" {
		return types.RawUser{}, http.StatusInternalServerError, errors.New("no token")
	}

	if !validSnowflake(id) {
		return types.RawUser{}, http.StatusBadRequest, errors.New("bad id")
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, apiBase+"/users/"+id, nil)
	req.Header.Set("Authorization", "Bot "+c.token)
	req.Header.Set("User-Agent", uagent)

	res, err := c.http.Do(req)
	if err != nil {
		return types.RawUser{}, http.StatusBadGateway, err
	}
	defer res.Body.Close()

	dec := json.NewDecoder(http.MaxBytesReader(nil, res.Body, maxBody))
	if res.StatusCode != http.StatusOK {
		var msg string
		switch res.StatusCode {
		case http.StatusNotFound:
			msg = "user not found"
		case http.StatusUnauthorized, http.StatusForbidden:
			msg = "invalid token"
		case http.StatusTooManyRequests:
			msg = "rate limited"
		default:
			msg = "api error"
		}
		return types.RawUser{}, res.StatusCode, fmt.Errorf("%s", msg)
	}

	var u types.RawUser
	if err := dec.Decode(&u); err != nil {
		return types.RawUser{}, http.StatusBadGateway, err
	}
	return u, http.StatusOK, nil
}

// DecodeBadges converts public flags into badge labels.
func DecodeBadges(flags int64) []string {
	badges := []struct {
		bit uint
		lbl string
	}{
		{0, "Discord Staff"},
		{1, "Partnered Owner"},
		{2, "HypeSquad Events"},
		{3, "Bug Hunter 1"},
		{6, "House Bravery"},
		{7, "House Brilliance"},
		{8, "House Balance"},
		{9, "Early Supporter"},
		{14, "Bug Hunter 2"},
		{16, "Verified Bot"},
		{17, "Early Bot Dev"},
		{18, "Moderator Alumni"},
		{22, "Active Developer"},
	}

	var out []string
	for _, b := range badges {
		if (flags & (1 << b.bit)) != 0 {
			out = append(out, b.lbl)
		}
	}
	return out
}

// CreatedAt returns the creation date of a Discord snowflake ID.
func CreatedAt(id string) string {
	sf, ok := new(big.Int).SetString(id, 10)
	if !ok {
		return ""
	}
	ms := new(big.Int).Rsh(sf, 22).Int64() + epoch
	return time.UnixMilli(ms).UTC().Format("01-02-2006")
}

// Avatar returns the URL of a user's avatar.
func Avatar(u types.RawUser) string {
	if u.Avatar != "" && u.Avatar != "null" {
		ext := "png"
		if strings.HasPrefix(u.Avatar, "a_") {
			ext = "gif"
		}
		return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.%s?size=1024", u.ID, u.Avatar, ext)
	}

	idx := 0
	if u.Discriminator != "" && u.Discriminator != "0" {
		if n, err := strconv.Atoi(u.Discriminator); err == nil {
			idx = n % 5
		}
	} else if bi, ok := new(big.Int).SetString(u.ID, 10); ok {
		idx = int(new(big.Int).Mod(bi, big.NewInt(5)).Int64())
	}
	return fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png", idx)
}

// Banner returns the URL of a user's banner.
func Banner(u types.RawUser) string {
	if u.Banner == "" || u.Banner == "null" {
		return ""
	}
	ext := "png"
	if strings.HasPrefix(u.Banner, "a_") {
		ext = "gif"
	}
	return fmt.Sprintf("https://cdn.discordapp.com/banners/%s/%s.%s?size=1024", u.ID, u.Banner, ext)
}

// validSnowflake checks if a string is a valid Discord snowflake ID.
func validSnowflake(id string) bool {
	if len(id) < 17 || len(id) > 19 {
		return false
	}

	for _, c := range id {
		if c < '0' || c > '9' {
			return false
		}
	}

	if _, ok := new(big.Int).SetString(id, 10); !ok {
		return false
	}

	return true
}
