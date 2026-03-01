package models

import "time"

/*
########## Data base models ############
*/
type UrlData struct {
	OriginalUrl string    `json:"original_url"`
	CreatedAt   time.Time `json:"create_at"`
	Salt        int64     `json:"salt"`
}

type CacheData struct {
	ShortenKey string
	Data       UrlData
}

/*
########## Customer Models ############
*/

type ShortenRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
}

type ShortenResponse struct {
	OriginalUrl string  `json:"original_url" binding:"required"`
	ShortenUrl  *string `json:"shorten_url"`
}
