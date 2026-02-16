package models

import "time"

/*
########## Data base models ############
*/
type UrlData struct {
	OriginalUrl string    `json:"original_url"`
	CreatedAt   time.Time `json:"create_at"`
	Counter     int       `json:"counter"`
	Passcode    string    `json:"passcode"`
	Salt        int64     `json:"salt"`
}

type RedirectData struct {
	OriginalUrl string `json:"original_url"`
	Counter     int    `json:"counter"`
	Salt        int64  `json:"salt"`
}

type CacheData struct {
	ShortenKey string
	Data       UrlData
}

/*
########## Customer Models ############
*/

type ShortenRequest struct {
	OriginalUrl string  `json:"original_url" binding:"required"`
	Passcode    *string `json:"passcode"` // needs to be a point for nil if not password and if no binding it becomes optional
}

type ShortenResponse struct {
	OriginalUrl string  `json:"original_url" binding:"required"`
	ShortenUrl  *string `json:"shorten_url"`
}
