package index

import (
	"context"
)

// Indexer allows indexing of documents.
type Indexer interface {
	Index(ctx context.Context, id string, properties map[string]interface{}) error
}

// Updater allows updating of documents.
type Updater interface {
	Update(ctx context.Context, id string, properties map[string]interface{}) error
}

// Getter allows getting of documents.
type Getter interface {
	Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error)
}

// Index represents an index which stores and retrieves document properties.
type Index interface {
	Indexer
	Updater
	Getter
}
