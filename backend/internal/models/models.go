package models

import "time"

type ShortKey struct {
	ShortUrl    string    `json:"short_url"`
	OriginalUrl string    `json:"original_url"`
	CreatedAt   time.Time `json:"create_at"`
}
