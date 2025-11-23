package model

import (
	"time"
)

type Team struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
