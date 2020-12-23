package sniffer

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ipfs-search/ipfs-search/queue"
	t "github.com/ipfs-search/ipfs-search/types"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-kad-dht/providers"
	"github.com/multiformats/go-base32"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SnifferTestSuite struct {
	suite.Suite
	ctx    context.Context
	cancel func()
	f      *queue.MockFactory
	ds     datastore.Batching
}

func (s *SnifferTestSuite) SetupTest() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.f = &queue.MockFactory{}
	s.f.Test(s.T())
	s.ds = datastore.NewMapDatastore()
}

func (s *SnifferTestSuite) TearDownTest() {
	s.cancel()
	s.ds.Close()
}

// TestNew does a burn test for New()
func (s *SnifferTestSuite) TestNew() {
	cfg := DefaultConfig()
	sniffy, e := New(cfg, s.ds, s.f)

	s.NotEmpty(sniffy)
	s.NoError(e)
}

// TestSniffCancel tests whether running Sniff() with a cancelled context returns with a context error.
func (s *SnifferTestSuite) TestSniffCancel() {
	cfg := DefaultConfig()
	sniffy, e := New(cfg, s.ds, s.f)
	s.NoError(e)

	// Cancel context
	s.cancel()

	// Setup Mock Publisher Factory
	qMock := &queue.Mock{}
	s.f.On("NewPublisher", mock.AnythingOfType("*context.cancelCtx")).Return(qMock, nil)

	err := sniffy.Sniff(s.ctx)
	s.Contains(err.Error(), "context canceled")

	s.f.AssertExpectations(s.T())
}

func timeToVal(t time.Time) []byte {
	// Ref: https://github.com/libp2p/go-libp2p-kad-dht/blob/master/providers/providers_manager.go#L239

	buf := make([]byte, 16)
	n := binary.PutVarint(buf, t.UnixNano())

	return buf[:n]
}

func makeKey(cidStr, provStr string) (datastore.Key, error) {
	testCid, err := cid.Decode(cidStr)
	if err != nil {
		return datastore.Key{}, err
	}

	testProv, err := peer.Decode(provStr)
	if err != nil {
		return datastore.Key{}, err
	}

	testCidBinary, err := testCid.MarshalBinary()
	if err != nil {
		return datastore.Key{}, err
	}

	testProvBinary, err := testProv.MarshalBinary()
	if err != nil {
		return datastore.Key{}, err
	}

	// Note; we need to explicitly encode to base32 here!
	encodedCid := base32.RawStdEncoding.EncodeToString(testCidBinary)
	encodedProv := base32.RawStdEncoding.EncodeToString(testProvBinary)

	return datastore.NewKey(fmt.Sprintf("%s%s/%s", providers.ProvidersKeyPrefix, encodedCid, encodedProv)), nil
}

// TestHandleToPublish tests the full chain from a yielded event to a publish.
func (s *SnifferTestSuite) TestHandleToPublish() {
	// Prepare provider key and value
	cidStr := "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp"
	provStr := "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd"

	key, err := makeKey(cidStr, provStr)
	s.NoError(err)

	now := time.Now()
	value := timeToVal(now)

	// Create sniffer
	cfg := DefaultConfig()
	sniffy, e := New(cfg, s.ds, s.f)
	s.NoError(e)

	// Get wrapped Datastore
	wrappedDs := sniffy.Batching()

	// Setup Mock Queue
	qMock := &queue.Mock{}
	qMock.On("Publish", mock.AnythingOfType("*context.valueCtx"), mock.MatchedBy(func(providerIf interface{}) bool {
		p := providerIf.(*t.Provider)
		s.Equal(p.Resource, &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       cidStr,
		})
		s.WithinDuration(p.Date, now, time.Second)
		s.Equal(p.Provider, provStr)
		return true
	}), uint8(9)).
		Return(nil).
		Run(func(args mock.Arguments) {
			fmt.Println("Publish() called, closing context")
			s.cancel()
		})

	// Setup Mock Queue Factory
	s.f.On("NewPublisher", mock.AnythingOfType("*context.cancelCtx")).Return(qMock, nil)

	// Start sniffing in goroutine
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		err := sniffy.Sniff(s.ctx)
		s.Contains(err.Error(), "context canceled")
		wg.Done()
	}()

	// Give the sniffer some time to start.
	// TODO: Create signal from all the way down that we are now sniffing
	time.Sleep(10 * time.Millisecond)

	// Put() test data to datastore, this *should* trigger a Publish event on the qMock
	s.NoError(wrappedDs.Put(key, value))

	wg.Wait()

	// Assert All The Things
	s.f.AssertExpectations(s.T())
	qMock.AssertExpectations(s.T())
}

// // TestLogToPublish tests the full chain from a log to a publish
// func (s *SnifferTestSuite) TestLogToPublish() {
// 	// Create queue and channels to retreive published messages and priorities
// 	pubs := make(chan interface{})
// 	priorities := make(chan uint8)
// 	q := &mockQueue{
// 		pubs:       pubs,
// 		priorities: priorities,
// 	}

// 	// Create sniffer
// 	cfg := DefaultConfig()
// 	s, e := New(cfg)
// 	s.NotEmpty(s)
// 	s.Empty(e)

// 	// Create buffered message channel and send mock message
// 	msgs := make(chan map[string]interface{}, 1)
// 	mockMsg := map[string]interface{}{
// 		"Duration":     33190,
// 		"Logs":         []string{},
// 		"Operation":    "handleAddProvider",
// 		"ParentSpanID": 0,
// 		"SpanID":       6.999711555735423e+18,
// 		"Start":        "2020-01-21T17:28:02.501941007Z",
// 		"Tags": map[string]interface{}{
// 			"key":    "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
// 			"peer":   "QmeTtFXm42Jb2todcKR538j6qHYxXt6suUzpF3rtT9FPSd",
// 			"system": "dht",
// 		},
// 		"TraceID": 4.483443946463055e+18,
// 	}
// 	msgs <- mockMsg

// 	// Create mock logger with associated messages
// 	l := mockLogger{
// 		msgs: msgs,
// 	}

// 	// Create cancelable context for sniffer to work with
// 	ctx, cancel := context.WithCancel(context.Background())

// 	// Run sniffer in goroutine
// 	go s.Sniff(ctx, l, q)

// 	// Retreive publication
// 	pub := <-pubs
// 	priority := <-priorities

// 	s.Equal(pub.(*crawler.Args).Hash, "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp")
// 	s.Equal(priority, uint8(9))

// 	// Cleanup
// 	cancel()
// }

func TestSnifferTestSuite(t *testing.T) {
	suite.Run(t, new(SnifferTestSuite))
}
