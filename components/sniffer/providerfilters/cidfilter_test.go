package providerfilters

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ipfs-search/ipfs-search/types"
)

var filter = NewCidFilter()

func makeProvider(resource *types.Resource) *types.Provider {
	if resource == nil {
		resource = &types.Resource{
			Protocol: types.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		}
	}

	return &types.Provider{
		Resource: resource,
		Date:     time.Now(),
		Provider: "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd",
	}
}

func TestCid0(t *testing.T) {
	assert := assert.New(t)

	p := makeProvider(nil)

	result, err := filter.Filter(*p)

	assert.Empty(err)
	assert.True(result)
}

func TestCid1Rawleave(t *testing.T) {
	assert := assert.New(t)

	r := &types.Resource{
		Protocol: types.IPFSProtocol,
		ID:       "bafkreiblvqc3q73ygovlzaxz4iilm5fopppcdc3uzkrtepjsgkvyev3kgy",
	}

	p := makeProvider(r)

	result, err := filter.Filter(*p)

	assert.Empty(err)
	assert.True(result)
}

func TestCid1Protobuf(t *testing.T) {
	assert := assert.New(t)

	r := &types.Resource{
		Protocol: types.IPFSProtocol,
		ID:       "bafybeihpsvpelgck42nikpiiuvgbf3ob3ydjkzkq5267mnp5jq5uhzatcy",
	}

	p := makeProvider(r)

	result, err := filter.Filter(*p)

	assert.Empty(err)
	assert.True(result)
}

func TestUnsupported(t *testing.T) {
	assert := assert.New(t)

	// Ethereum block
	r := &types.Resource{
		Protocol: types.IPFSProtocol,
		ID:       "z43AaGEvwdfzjrCZ3Sq7DKxdDHrwoaPQDtqF4jfdkNEVTiqGVFW",
	}

	p := makeProvider(r)

	_, err := filter.Filter(*p)

	assert.True(errors.Is(err, errUnsupportedCodec))
}

func TestInvalid(t *testing.T) {
	assert := assert.New(t)

	invalidResource := &types.Resource{
		Protocol: types.IPFSProtocol,
		ID:       "invalid",
	}

	p := makeProvider(invalidResource)

	_, err := filter.Filter(*p)

	assert.True(errors.Is(err, errDecodingCID))
}

func TestNonIpfs(t *testing.T) {
	assert := assert.New(t)

	invalidResource := &types.Resource{
		Protocol: types.IPFSProtocol + 1,
		ID:       "dontcare",
	}

	p := makeProvider(invalidResource)

	_, err := filter.Filter(*p)

	assert.True(errors.Is(err, errUnsupportedProtocol))
}
