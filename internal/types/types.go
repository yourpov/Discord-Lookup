package types

import "time"

// JSONTime is a time type for formatting time in JSON responses
type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).UTC().Format("2006-01-02 15:04 MST")
	return []byte(`"` + formatted + `"`), nil
}

type Config struct {
	Port  string `json:"port"`
	Token string `json:"token"`
}

// RawUser is the shape returned by discords api
type RawUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	GlobalName    string `json:"global_name"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Banner        string `json:"banner"`
	Bot           bool   `json:"bot"`
	System        bool   `json:"system"`
	PublicFlags   int64  `json:"public_flags"`
}

// User is what our API returns
type User struct {
	ID            string   `json:"id"`
	Username      string   `json:"username"`
	DisplayName   string   `json:"display_name"`
	Discriminator string   `json:"discriminator"`
	Bot           bool     `json:"bot"`
	System        bool     `json:"system"`
	Flags         int64    `json:"flags"`
	Badges        []string `json:"badges"`
	Avatar        string   `json:"avatar"`
	Banner        string   `json:"banner"`
	CreatedAt     string   `json:"created_at"`
	SearchedAt    JSONTime `json:"searched_at"` // custom time formatting

}
