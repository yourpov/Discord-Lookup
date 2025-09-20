package discord

import (
	"context"
	"discord-lookup/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	apiBase      = "https://discord.com/api/v10"
	discordEpoch = int64(1420070400000) // 2015-01-01
	maxBody      = 1 << 20              // 1 mb
	reqTimeout   = 10 * time.Second
	UserAgent    = "discord-lookup (https://github.com/yourpov/discord-lookup, 1.0)"
)

type Client struct {
	http  *http.Client
	token string
}

// New creates a new Discord API client
func New(token string) *Client {
	return &Client{
		http:  &http.Client{Timeout: reqTimeout},
		token: token,
	}
}

// FetchUser fetches a user by ID
func (c *Client) FetchUser(ctx context.Context, id string) (types.RawUser, int, error) {
	if c.token == "" {
		return types.RawUser{}, http.StatusInternalServerError, errors.New("missing bot token")
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, apiBase+"/users/"+id, nil)
	req.Header.Set("Authorization", "Bot "+c.token)
	req.Header.Set("User-Agent", UserAgent)

	res, err := c.http.Do(req)
	if err != nil {
		return types.RawUser{}, http.StatusBadGateway, err
	}
	defer res.Body.Close()

	dec := json.NewDecoder(http.MaxBytesReader(nil, res.Body, maxBody))
	if res.StatusCode != http.StatusOK {
		return types.RawUser{}, res.StatusCode, fmt.Errorf("invalid snowflake")
	}

	var u types.RawUser
	if err := dec.Decode(&u); err != nil {
		return types.RawUser{}, http.StatusBadGateway, err
	}
	return u, http.StatusOK, nil
}

// DecodeBadges converts the public flags into badge names
func DecodeBadges(flags int64) []string {
	badgeBits := []struct {
		Bit   uint
		Label string
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
	for _, b := range badgeBits {
		if (flags & (1 << b.Bit)) != 0 {
			out = append(out, b.Label)
		}
	}
	return out
}

// CreatedAt takes the account creation timestamp from a snowflake ID
func CreatedAt(id string) string {
	sf, ok := new(big.Int).SetString(id, 10)
	if !ok {
		return ""
	}
	ms := new(big.Int).Rsh(sf, 22).Int64() + discordEpoch
	return time.UnixMilli(ms).UTC().Format("01-02-2006")
}

// Avatar builds a avatar URL if present, otherwise a default avatar URL
func Avatar(u types.RawUser) string {
	if u.Avatar != "" && u.Avatar != "null" {
		ext := "png"
		if strings.HasPrefix(u.Avatar, "a_") {
			ext = "gif"
		}
		return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.%s?size=1024", u.ID, u.Avatar, ext)
	}

	// Default avatar based on discriminator or ID
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

// Banner builds a banner URL
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
