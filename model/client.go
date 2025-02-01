package model

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ClientID  uuid.UUID `json:"client_id"`
	Email     string    `json:"email"`
	Suspended bool      `json:"suspended"`
	Plan      string    `json:"plan"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

