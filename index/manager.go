package index

import (
	"context"
)

// Exister allows to check for existence of the index.
type Exister interface {
	Exists(ctx context.Context) (bool, error)
}

// Creator allows to create the index.
type Creator interface {
	Create(ctx context.Context) error
}

// ConfigUpdater represents an index with configuration.
type ConfigUpdater interface {
	ConfigUpToDate(context.Context) (bool, error)
	ConfigUpdate(context.Context) error
}

// ManagedIndex is an index which allows management
// TODO: Factor this into the index initialiser (automate it away) or add as specific -single- Bootstrap() or Prepare() method.
type ManagedIndex interface {
	Index
	Exister
	Creator
	ConfigUpdater
}
