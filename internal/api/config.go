package api

import (
	"github.com/jexlor/votingapp/db/store"
)

type Config struct {
	DB *store.Queries
}
