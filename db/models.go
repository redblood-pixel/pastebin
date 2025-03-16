// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type AccessType string

const (
	AccessTypePublic  AccessType = "public"
	AccessTypePrivate AccessType = "private"
)

func (e *AccessType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AccessType(s)
	case string:
		*e = AccessType(s)
	default:
		return fmt.Errorf("unsupported scan type for AccessType: %T", src)
	}
	return nil
}

type NullAccessType struct {
	AccessType AccessType
	Valid      bool // Valid is true if AccessType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAccessType) Scan(value interface{}) error {
	if value == nil {
		ns.AccessType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AccessType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAccessType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AccessType), nil
}

type Paste struct {
	ID              int32
	Title           string
	ContentLocation string
	CreatedAt       pgtype.Timestamp
	ExpiresAt       pgtype.Timestamp
	Visibility      AccessType
	LastVisited     pgtype.Timestamp
	UserID          pgtype.Int4
}

type User struct {
	ID             int32
	Name           string
	Email          string
	CreatedAt      pgtype.Timestamp
	LastLogin      pgtype.Timestamp
	PasswordHashed string
}
