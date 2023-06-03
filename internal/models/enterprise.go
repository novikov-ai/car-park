package models

import "time"

type Enterprise struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	City        string    `json:"city"`
	Established time.Time `json:"established"`
}
