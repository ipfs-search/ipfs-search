package aliasresolver

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	opensearch "github.com/opensearch-project/opensearch-go/v2"
)

var (
	ErrNotFound = errors.New("not found")
)

type AliasResolver interface {
	GetIndex(ctx context.Context, aliasName string) (string, error)
	GetAlias(ctx context.Context, indexName string) (string, error)
}

type defaultAliasResolver struct {
	mu              sync.RWMutex
	client          *opensearch.Client
	lastRefresh     time.Time
	refreshDuration time.Duration
	indexToAlias    map[string]string // Map from index name to alias name
	aliasToIndex    map[string]string // Map from alias name to index name
}

func NewAliasResolver(client *opensearch.Client) AliasResolver {
	return &defaultAliasResolver{
		client:          client,
		refreshDuration: 5 * time.Minute,
	}
}

type aliasesResponse map[string]struct {
	Aliases map[string]struct{} `json:"aliases"`
}

func (r *defaultAliasResolver) refreshAliases(ctx context.Context) error {
	// Get all aliases
	// TODO: Only get configured aliases/indexes.
	resp, err := r.client.API.Indices.GetAlias(
		r.client.API.Indices.GetAlias.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.IsError() {
		return errors.New(resp.String())
	}

	// Decode response
	var result aliasesResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	log.Printf("%s", result)

	r.mu.Lock()

	// Assign indexToAlias
	r.indexToAlias = make(map[string]string, len(result))
	for indexName, v := range result {
		if len(v.Aliases) != 1 {
			panic("unexpected amount of aliases for index '" + indexName + "'")
		}

		var aliasName string
		for aliasName = range v.Aliases {
			// Use the first alias name
			break
		}
		r.indexToAlias[indexName] = aliasName
	}

	// Assign aliasToIndex
	r.aliasToIndex = make(map[string]string, len(result))
	for indexName, aliasName := range r.indexToAlias {
		r.aliasToIndex[aliasName] = indexName
	}

	r.lastRefresh = time.Now()

	r.mu.Unlock()

	return nil
}

func (r *defaultAliasResolver) conditionalRefresh(ctx context.Context) bool {
	r.mu.RLock()
	if r.lastRefresh.Add(r.refreshDuration).Before(time.Now()) {
		r.mu.RUnlock()
		r.refreshAliases(ctx)
		return true
	}

	r.mu.RUnlock()
	return false
}

func (r *defaultAliasResolver) getCachedIndex(aliasName string) (string, bool) {
	r.mu.RLock()
	indexName, found := r.aliasToIndex[aliasName]
	r.mu.RUnlock()

	return indexName, found
}

func (r *defaultAliasResolver) GetIndex(ctx context.Context, aliasName string) (string, error) {
	refreshed := r.conditionalRefresh(ctx)
	indexName, found := r.getCachedIndex(aliasName)

	if found {
		return indexName, nil
	}

	if refreshed {
		// Not found, but we've already refreshed - just return.
		return "", ErrNotFound
	}

	// Not refreshed, try refreshing.
	err := r.refreshAliases(ctx)
	if err != nil {
		return "", err
	}

	indexName, found = r.getCachedIndex(aliasName)

	if found {
		return indexName, nil
	}

	return "", ErrNotFound
}

func (r *defaultAliasResolver) getCachedAlias(indexName string) (string, bool) {
	r.mu.RLock()
	aliasName, found := r.indexToAlias[indexName]
	r.mu.RUnlock()

	return aliasName, found
}

func (r *defaultAliasResolver) GetAlias(ctx context.Context, indexName string) (string, error) {
	refreshed := r.conditionalRefresh(ctx)
	aliasName, found := r.getCachedAlias(indexName)

	if found {
		return aliasName, nil
	}

	if refreshed {
		// Not found, but we've already refreshed - just return.
		return "", ErrNotFound
	}

	// Not refreshed, try refreshing.
	err := r.refreshAliases(ctx)
	if err != nil {
		return "", err
	}

	aliasName, found = r.getCachedAlias(indexName)

	if found {
		return aliasName, nil
	}

	return "", ErrNotFound
}
