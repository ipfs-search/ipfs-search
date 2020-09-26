package eventsource

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-eventbus"
	"github.com/libp2p/go-libp2p-core/event"

	"github.com/stretchr/testify/suite"
)

type EventSourceTestSuite struct {
	suite.Suite
	ctx    context.Context
	cancel func()
	bus    event.Bus
	ds     datastore.Batching
}

func (s *EventSourceTestSuite) SetupTest() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.bus = eventbus.NewBus()
	s.ds = datastore.NewMapDatastore()
}

func (s *EventSourceTestSuite) TearDownTest() {
	s.ds.Close()
	s.cancel()
}

func (s *EventSourceTestSuite) TestNew() {
	es, err := New(s.bus, s.ds)
	s.NoError(err)
	s.NotEmpty(es)
}

func timeToVal(t time.Time) []byte {
	// Ref: https://github.com/libp2p/go-libp2p-kad-dht/blob/master/providers/providers_manager.go#L239

	buf := make([]byte, 16)
	n := binary.PutVarint(buf, t.UnixNano())

	return buf[:n]
}

func (s *EventSourceTestSuite) TestSubscribePut() {
	// Create an EventSource
	es, _ := New(s.bus, s.ds)

	// Create a waitgroup to wait for the goroutine
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Setup a listener
	go es.Subscribe(s.ctx, func(ctx context.Context, e EvtProviderPut) error {
		s.Equal(e.CID.String(), "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp")
		s.Equal(e.PeerID.String(), "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd")

		wg.Done()

		return nil
	})

	// Give the goroutine some time to start
	time.Sleep(10 * time.Millisecond)

	// Get the wrapped datastore
	proxyDs := es.Batching()

	// Perform the Put()
	k := datastore.NewKey("/providers/CIQDWKPBHXLJ3XVELRJZA2SYY7OGCSX6FRSIZS2VQQPVKOA2Z4VXN2I/CIQO7FK6IWMEVZU2QU6QRJKMCLW4DXQGSVSVB3V56Y272TB3IPSBGFQ")
	v := timeToVal(time.Now())
	err := proxyDs.Put(k, v)
	s.NoError(err)

	fmt.Printf("Waiting for Put event.")
	wg.Wait()
}

func TestEventSourceTestSuite(t *testing.T) {
	suite.Run(t, new(EventSourceTestSuite))
}
