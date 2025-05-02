package domain

import (
	"time"

	"github.com/google/uuid"
)

// 5 years
const DefaultTTL = time.Hour * 24 * 365 * 5

const (
	PublicType  = "public"
	PrivateType = "private"
)

type Paste struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title" validate:"required"`
	CreatedAt   time.Time     `json:"created_at"`
	ExpiresAt   time.Time     `json:"expires_at"`
	TTL         time.Duration `json:"ttl"`
	Visibility  string        `json:"visibility" validate:"oneof='public private'"`
	LastVisited time.Time     `json:"last_visited"`
	UserID      int           `json:"-"`
}
