package crawler

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/ipfs-search/ipfs-search/extractor"
	"github.com/ipfs-search/ipfs-search/index"
	// indexTypes "github.com/ipfs-search/ipfs-search/index/types"
	"github.com/ipfs-search/ipfs-search/protocol"
	"github.com/ipfs-search/ipfs-search/queue"
	t "github.com/ipfs-search/ipfs-search/types"
)

type CrawlerTestSuite struct {
	suite.Suite

	ctx     context.Context
	indexes Indexes
	queues  Queues
	c       *Crawler

	protocol   *protocol.Mock
	extractor  *extractor.Mock
	fileIdx    *index.Mock
	dirIdx     *index.Mock
	invalidIdx *index.Mock
	hashQ      *queue.Mock
	dirQ       *queue.Mock
	fileQ      *queue.Mock
}

func (s *CrawlerTestSuite) SetupTest() {
	s.ctx = context.Background()

	// Creat a crawler with mocked dependencies
	s.fileIdx, s.dirIdx, s.invalidIdx = &index.Mock{}, &index.Mock{}, &index.Mock{}

	s.indexes = Indexes{
		Files:       s.fileIdx,
		Directories: s.dirIdx,
		Invalid:     s.invalidIdx,
	}

	s.hashQ, s.fileQ, s.dirQ = &queue.Mock{}, &queue.Mock{}, &queue.Mock{}

	s.queues = Queues{
		Hashes:      s.hashQ,
		Directories: s.fileQ,
		Files:       s.dirQ,
	}
	s.protocol = &protocol.Mock{}
	s.extractor = &extractor.Mock{}

	s.c = New(s.indexes, s.queues, s.protocol, s.extractor)
}

func (s *CrawlerTestSuite) assertExpectations() {
	mock.AssertExpectationsForObjects(s.T(),
		s.fileIdx,
		s.dirIdx,
		s.invalidIdx,
		s.hashQ,
		s.fileQ,
		s.dirQ,
		s.protocol,
		s.extractor,
	)
}

func (s *CrawlerTestSuite) TestCrawlInvalidProtocol() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.InvalidProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
	}

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// This is a programming error, should fail hard.
	s.Error(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlUndefinedType() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Stat: t.Stat{
			Type: t.UndefinedType,
		},
	}

	// Mock assertions
	s.protocol.On("Stat", mock.Anything, r).Run(func(args mock.Arguments) {
		r := args.Get(1).(*t.AnnotatedResource)
		r.Stat = t.Stat{
			Type: t.FileType,
		}
	}).Return(nil)

	s.extractor.On("Extract", mock.Anything, r, mock.Anything).Run(func(args mock.Arguments) {
		// m := args.Get(2)
		// Set metadata on this interface{}
	}).Return(nil)

	// s.fileIdx.On("Index", mock.Anything, r.Resource.ID, indexTypes.File{
	// FirstSeen  time.Time  `json:"first-seen"`
	// LastSeen   time.Time  `json:"last-seen"`
	// References References `json:"references"`
	// Size       uint64     `json:"size"`

	// Content         string   `json:"content"`
	// IpfsTikaVersion string   `json:"ipfs_tika_version"`
	// Language        Language `json:"language"`
	// Metadata        Metadata `json:"metadata"`
	// Urls            []string `json:"urls"`
	// }).Return(nil)

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Undefined types should index based on Stat'ed type.
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlUnsupportedType() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Stat: t.Stat{
			Type: t.UnsupportedType,
		},
	}

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Undefined types should index as invalid.
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlUnreferencedPartialType() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Stat: t.Stat{
			Type: t.PartialType,
		},
	}

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Unreferenced partial should be skipped.
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlReferencedPartialType() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmafrLBfzRLV4XSH1XcaMMeaXEUhDJjmtDfsYU95TrWG87",
			},
			Name: "fileName.pdf",
		},
		Stat: t.Stat{
			Type: t.PartialType,
		},
	}

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Referenced partials should be indexed.
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlFileType() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmafrLBfzRLV4XSH1XcaMMeaXEUhDJjmtDfsYU95TrWG87",
			},
			Name: "fileName.pdf",
		},
		Stat: t.Stat{
			Type: t.FileType,
		},
	}

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlDirectoryType() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmafrLBfzRLV4XSH1XcaMMeaXEUhDJjmtDfsYU95TrWG87",
			},
			Name: "directoryName",
		},
		Stat: t.Stat{
			Type: t.DirectoryType,
		},
	}

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.NoError(err)
	s.assertExpectations()
}

func TestCrawlerTestSuite(t *testing.T) {
	suite.Run(t, new(CrawlerTestSuite))
}
