package todo

import "time"

type TODO struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
}
