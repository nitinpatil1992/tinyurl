package model

import "time"

type Form struct {
	LongURL string
}

type TinyUrl struct {
	ShortURL  string    `json:"short_url"`
	LongURL   string    `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
}

type UrlVisitsModel struct {
	ShortURL  string `json:"short_url"`
	VisitedAt string `json:"visited_at"`
}
