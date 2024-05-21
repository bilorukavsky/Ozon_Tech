package graph

import (
	"PostCommentService/db"
)

type Resolver struct {
	store db.Store
}

func NewResolver(useMemory bool) *Resolver {
	store := db.NewStore(useMemory)
	return &Resolver{store: store}
}
