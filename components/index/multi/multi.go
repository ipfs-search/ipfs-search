package multi

import (
	"context"

	"github.com/ipfs-search/ipfs-search/components/index"
)

type Multi struct {
	indexMap  map[string]index.Index
	indexList []index.Index
	selector  Selector
}

func New(factory index.Factory, selector Selector) index.Index {
	m := &Multi{
		selector: selector,
	}

	m.setIndexes(factory)

	return m
}

func (m *Multi) setIndexes(factory index.Factory) {
	if m.selector == nil {
		panic("selector not specified.")
	}

	indexNames := m.selector.ListIndexes()
	indexCount := len(indexNames)

	m.indexList = make([]index.Index, indexCount)
	m.indexMap = make(map[string]index.Index, indexCount)

	for i, name := range indexNames {
		index := factory.NewIndex(name)
		m.indexList[i] = index
		m.indexMap[name] = index
	}
}

func (m *Multi) getIndex(id string, properties interface{}) index.Index {
	p, ok := properties.(Properties)
	if !ok {
		panic("Cannot assert Properties.")
	}

	indexName := m.selector.GetIndex(id, p)
	return m.indexMap[indexName]
}

func (m *Multi) Index(ctx context.Context, id string, properties interface{}) error {
	index := m.getIndex(id, properties)
	return index.Index(ctx, id, properties)
}

func (m *Multi) Update(ctx context.Context, id string, properties interface{}) error {
	index := m.getIndex(id, properties)
	return index.Index(ctx, id, properties)
}

func (m *Multi) Get(ctx context.Context, id string, dst interface{}, fields ...string) (bool, error) {
	index, err := index.MultiGet(ctx, m.indexList, id, dst, fields...)
	if err != nil {
		return false, err
	}

	if index != nil {
		return true, nil
	}

	return false, nil
}

func (m *Multi) Delete(ctx context.Context, id string) error {
	panic("Not implemented.")
}

// Compile-time assurance that implementation satisfies interface.
var _ index.Index = &Multi{}
