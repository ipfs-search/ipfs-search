package factory

import (
	"context"
	// "encoding/binary"
	// "fmt"
	// "sync"
	"testing"
	// "time"

	// "github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	// "github.com/libp2p/go-libp2p-core/peer"
	// "github.com/libp2p/go-libp2p-kad-dht/providers"
	// "github.com/multiformats/go-base32"
	// "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/ipfs-search/ipfs-search/components/queue"
	// "github.com/ipfs-search/ipfs-search/instr"
	// t "github.com/ipfs-search/ipfs-search/types"
)

type FactoryTestSuite struct {
	suite.Suite
	ctx    context.Context
	cancel func()
	f      *queue.MockFactory
	ds     datastore.Batching
}

func (s *FactoryTestSuite) SetupTest() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.f = &queue.MockFactory{}
	s.f.Test(s.T())
	s.ds = datastore.NewMapDatastore()
}

func (s *FactoryTestSuite) TearDownTest() {
	s.cancel()
	s.ds.Close()
}

// TestStartBurn performs a burn test for Start().
func (s *FactoryTestSuite) TestStartBurn() {
	ctx, ds, err := Start(s.ctx, s.ds)
	s.NoError(err)

	s.NotEqual(s.ds, ds)
	s.NotEqual(s.ctx, ctx)
}

func TestFactoryTestSuite(t *testing.T) {
	suite.Run(t, new(FactoryTestSuite))
}
