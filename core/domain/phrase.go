package domain

import "time"

type Phrase struct {
	ID          *int64
	Owner       string
	State       string
	Phrase      string
	Author      string
	CreatedAt   time.Time
	PublishedAt time.Time
}
