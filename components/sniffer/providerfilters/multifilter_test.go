package providerfilters

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ipfs-search/ipfs-search/types"
)

func TestPassSingle(t *testing.T) {
	assert := assert.New(t)

	f := MockFilter{
		R: true,
	}

	m := NewMultiFilter(&f)
	r, err := m.Filter(types.Provider{})

	assert.True(r)
	assert.Empty(err)
}

func TestRejectSingle(t *testing.T) {
	assert := assert.New(t)

	f := MockFilter{
		R: false,
	}

	m := NewMultiFilter(&f)
	r, err := m.Filter(types.Provider{})

	assert.False(r)
	assert.Empty(err)
}

func TestPassTwo(t *testing.T) {
	assert := assert.New(t)

	f := MockFilter{
		R: true,
	}

	m := NewMultiFilter(&f, &f)
	r, err := m.Filter(types.Provider{})

	assert.True(r)
	assert.Empty(err)

	assert.Equal(2, f.Calls)
}

func TestFailFirst(t *testing.T) {
	assert := assert.New(t)

	passFilter := MockFilter{
		R: true,
	}

	rejectFilter := MockFilter{
		R: false,
	}

	m := NewMultiFilter(&rejectFilter, &passFilter)
	r, err := m.Filter(types.Provider{})

	assert.False(r)
	assert.Empty(err)

	assert.Equal(0, passFilter.Calls)
	assert.Equal(1, rejectFilter.Calls)
}

func TestFailSecond(t *testing.T) {
	assert := assert.New(t)

	passFilter := MockFilter{
		R: true,
	}

	rejectFilter := MockFilter{
		R: false,
	}

	m := NewMultiFilter(&passFilter, &rejectFilter)
	r, err := m.Filter(types.Provider{})

	assert.False(r)
	assert.Empty(err)

	assert.Equal(1, passFilter.Calls)
	assert.Equal(1, rejectFilter.Calls)
}

func TestError(t *testing.T) {
	assert := assert.New(t)

	mockErr := errors.New("test")

	f := MockFilter{
		Err: mockErr,
	}

	m := NewMultiFilter(&f)
	_, err := m.Filter(types.Provider{})

	assert.True(errors.Is(err, errFilter))
	assert.Contains(err.Error(), "filter error: test")
}
