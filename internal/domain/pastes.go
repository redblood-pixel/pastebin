package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
)

const DefaultTTL = time.Hour * 24 * 365 * 5            // 5 years
const DefaultLastVisitedTTL = time.Hour * 24 * 365 * 2 // 2 years

const (
	PublicType  = "public"
	PrivateType = "private"
)

type Paste struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title" validate:"required"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	Visibility     string    `json:"visibility" validate:"oneof='public private'"`
	LastVisited    time.Time `json:"last_visited"`
	BurnAfterRead  bool
	UserID         int `json:"-"`
	PasswordHashed pgtype.Text
}

type PasteParameters struct {
	Password string `json:"password"`
}
