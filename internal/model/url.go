package model

import "time"

type URL struct {
	ID        int
	Original  string
	ShortCode string
	CreatedAt time.Time
	ExpiresAt time.Time
	Hits      int
}

type AccessLog struct {
	ID         int
	URLID      int
	AccessedAt time.Time
	IP         string
	Action     string // "shorten" or "redirect"
}
