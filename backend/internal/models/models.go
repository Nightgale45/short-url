package models

import "time"

type ShortKey struct {
	ShortUrl    string    `json:"short_url"`
	OriginalUrl string    `json:"original_url"`
	CreatedAt   time.Time `json:"create_at"`
	Counter     string    `json:"counter"`
	Passcode    string    `json:"passcode"`
}

type ShortenRequest struct {
	OriginalUrl string  `json:"original_url" binding:"required"`
	Passcode    *string `json:"passcode"` // needs to be a point for nil if not password and if no binding it becomes optional
}

type ShortenResponse struct {
	OriginalUrl string  `json:"original_url" binding:"required"`
	ShortenUrl  *string `json:"shorten_url"`
	Error       *string `json:"error"`
}
