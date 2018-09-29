package models

import (
	"time"
)

// Storm Model struct
// cannot put ID here because storm doesn't like it ¯\_(ツ)_/¯
type Model struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
