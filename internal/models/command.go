package models

import (
	"time"
)

type Command struct {
	ID        uint      `json:"id"`
	Command   string    `json:"command"`
	Status    string    `json:"status"`
	Output    string    `json:"output"`
	CreatedAt time.Time `json:"created_at"`
}
