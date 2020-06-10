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
	Create(ctx context.Context, settings interface{}, mapping interface{}) error
}

// SettingsGetter returns the settings for an index.
type SettingsGetter interface {
	GetSettings(ctx context.Context) (settings interface{}, err error)
}

// SettingsSetter allows to update the settings of the index
type SettingsSetter interface {
	SetSettings(ctx context.Context, settings interface{}) error
}

// MappingGetter returns the mapping for an index.
type MappingGetter interface {
	GetMapping(ctx context.Context) (mapping interface{}, err error)
}

// MappingSetter sets the mapping for an index.
type MappingSetter interface {
	SetMapping(ctx context.Context, mapping interface{}) error
}

// Manager allows management of indexes
type Manager interface {
	Exister
	Creator
	SettingsGetter
	SettingsSetter
	MappingGetter
	MappingSetter
}

// ManagedIndex is an index which also allows management
type ManagedIndex interface {
	Index
	Manager
}
