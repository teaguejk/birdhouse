package ai

import "time"

type Config struct {
	APIKey  string        `json:"-"`
	APIURL  string        `json:"api_url"`
	Model   string        `json:"model"`
	Timeout time.Duration `json:"timeout"`
}
