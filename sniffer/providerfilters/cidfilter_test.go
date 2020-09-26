package providerfilters

import (
	"github.com/ipfs-search/ipfs-search/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var filter = NewCidFilter()

func makeProvider(resource *types.Resource) *types.Provider {
	if resource == nil {
		resource = &types.Resource{
			Protocol: "ipfs",
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
		Protocol: "ipfs",
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
		Protocol: "ipfs",
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
		Protocol: "ipfs",
		ID:       "z43AaGEvwdfzjrCZ3Sq7DKxdDHrwoaPQDtqF4jfdkNEVTiqGVFW",
	}

	p := makeProvider(r)

	_, err := filter.Filter(*p)

	assert.Error(err)
	assert.Contains(err.Error(), "Unsupported codec")
}

func TestInvalid(t *testing.T) {
	assert := assert.New(t)

	invalidResource := &types.Resource{
		Protocol: "ipfs",
		ID:       "invalid",
	}

	p := makeProvider(invalidResource)

	_, err := filter.Filter(*p)

	assert.Error(err)
	assert.Contains(err.Error(), "decoding CID")
}

func TestNonIpfs(t *testing.T) {
	assert := assert.New(t)

	invalidResource := &types.Resource{
		Protocol: "bananafs",
		ID:       "dontcare",
	}

	p := makeProvider(invalidResource)

	_, err := filter.Filter(*p)

	assert.Error(err)
	assert.Contains(err.Error(), "Unsupported protocol")

}
