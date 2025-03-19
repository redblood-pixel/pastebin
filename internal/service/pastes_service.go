package service

import (
	"github.com/redblood-pixel/pastebin/pkg/postgres_queries"
)

type PastesService struct {
	db *postgres_queries.Queries
}
