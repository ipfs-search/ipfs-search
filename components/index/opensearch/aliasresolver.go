package opensearch

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"

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
	mu           sync.RWMutex
	client       *opensearch.Client
	indexToAlias map[string]string // Map from index name to alias name
	aliasToIndex map[string]string // Map from alias name to index name
}

func NewAliasResolver(client *opensearch.Client) AliasResolver {
	return &defaultAliasResolver{
		client: client,
	}
}

type aliasesResponse map[string]struct {
	Aliases map[string]struct{} `json:"aliases"`
}

func (r *defaultAliasResolver) refreshAliases(ctx context.Context) error {
	// Get all aliases
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

	r.mu.Unlock()

	return nil
}

func (r *defaultAliasResolver) GetIndex(ctx context.Context, aliasName string) (string, error) {
	r.mu.RLock()
	indexName, ok := r.aliasToIndex[aliasName]
	r.mu.RUnlock()

	if ok {
		return indexName, nil
	}

	err := r.refreshAliases(ctx)
	if err != nil {
		return "", err
	}

	r.mu.RLock()
	indexName, ok = r.aliasToIndex[aliasName]
	r.mu.RUnlock()

	if !ok {
		return "", ErrNotFound
	}

	return indexName, nil
}

func (r *defaultAliasResolver) GetAlias(ctx context.Context, indexName string) (string, error) {
	r.mu.RLock()
	aliasName, ok := r.indexToAlias[indexName]
	r.mu.RUnlock()

	if ok {
		return aliasName, nil
	}

	err := r.refreshAliases(ctx)
	if err != nil {
		return "", err
	}

	r.mu.RLock()
	aliasName, ok = r.indexToAlias[indexName]
	r.mu.RUnlock()

	if !ok {
		return "", ErrNotFound
	}

	return aliasName, nil
}
