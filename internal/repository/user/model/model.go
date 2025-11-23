package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	IsActive  bool         `json:"is_active"`
	TeamName  string       `json:"team_name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
