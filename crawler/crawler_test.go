package crawler

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/ipfs-search/ipfs-search/extractor"
	"github.com/ipfs-search/ipfs-search/index"
	indexTypes "github.com/ipfs-search/ipfs-search/index/types"
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

	protocol  *protocol.Mock
	extractor *extractor.Mock

	fileIdx    *index.Mock
	dirIdx     *index.Mock
	invalidIdx *index.Mock

	dirQ  *queue.Mock
	fileQ *queue.Mock
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

	s.fileQ, s.dirQ = &queue.Mock{}, &queue.Mock{}

	s.queues = Queues{
		Directories: s.dirQ,
		Files:       s.fileQ,
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
	s.Panics(func() { _ = s.c.Crawl(s.ctx, r) })
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
	s.protocol.
		On("Stat", mock.Anything, r).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*t.AnnotatedResource)
			r.Stat = t.Stat{
				Type: t.FileType,
			}
		}).
		Return(nil).
		Once()

	s.extractor.
		On("Extract", mock.Anything, r, mock.Anything).
		Return(nil).
		Once()

	s.fileIdx.
		On("Index", mock.Anything, r.Resource.ID, mock.IsType(&indexTypes.File{})).
		Return(nil).
		Once()

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

	// Mock assertions
	s.invalidIdx.
		On("Index", mock.Anything, r.Resource.ID, mock.MatchedBy(func(f *indexTypes.Invalid) bool {
			return s.Equal(f.Error, "unsupported type")
		})).
		Return(nil).
		Once()

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Undefined types should index as invalid.
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlPartialType() {
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

	// Note how nothing is indexed here!

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Unreferenced partial should be skipped.
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
		Stat: t.Stat{
			Type: t.FileType,
			Size: 15,
		},
	}

	// Mock assertions
	testMetadata := indexTypes.Metadata{"TestField": "TestValue"}

	s.extractor.
		On("Extract", mock.Anything, r, mock.Anything).
		Run(func(args mock.Arguments) {
			f := args.Get(2).(*indexTypes.File)
			f.Content = "testContent"
			f.Metadata = testMetadata
		}).
		Return(nil).
		Once()

	s.fileIdx.
		On("Index", mock.Anything, r.Resource.ID, mock.MatchedBy(func(f *indexTypes.File) bool {
			return s.Equal(f.Metadata, testMetadata) &&
				s.Equal(f.Content, "testContent") &&
				s.Equal(f.Size, uint64(15))
		})).
		Return(nil).
		Once()

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
		Stat: t.Stat{
			Type: t.DirectoryType,
			Size: 23,
		},
	}

	// Mock assertions
	fileEntry := t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmafrLBfzRLV4XSH1XcaMMeaXEUhDJjmtDfsYU95TrWG87",
		},
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
			},
			Name: "fileName.pdf",
		},
		Stat: t.Stat{
			Type: t.FileType,
			Size: 3431,
		},
	}

	dirEntry := t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv",
		},
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
			},
			Name: "dirName",
		},
		Stat: t.Stat{
			Type: t.DirectoryType,
			Size: 4534543,
		},
	}

	s.protocol.
		On("Ls", mock.Anything, r, mock.AnythingOfType("chan<- *types.AnnotatedResource")).
		Run(func(args mock.Arguments) {
			// Write bogus directory entry
			entryChan := args.Get(2).(chan<- *t.AnnotatedResource)
			entryChan <- &fileEntry
			entryChan <- &dirEntry
		}).
		Return(nil).
		Once()

	s.dirIdx.
		On("Index", mock.Anything, r.Resource.ID, mock.MatchedBy(func(f *indexTypes.Directory) bool {
			return s.Equal(f.Size, r.Size) &&
				s.Equal(f.Links, indexTypes.Links{
					indexTypes.Link{
						Hash: fileEntry.ID,
						Name: fileEntry.Reference.Name,
						Size: fileEntry.Size,
						Type: indexTypes.FileLinkType,
					},
					indexTypes.Link{
						Hash: dirEntry.ID,
						Name: dirEntry.Reference.Name,
						Size: dirEntry.Size,
						Type: indexTypes.DirectoryLinkType,
					},
				})
		})).
		Return(nil).
		Once()

	s.fileQ.
		On("Publish", mock.Anything, mock.MatchedBy(func(f *t.AnnotatedResource) bool {
			return s.Equal(*f, fileEntry)
		}), mock.AnythingOfType("uint8")).
		Return(nil).
		Once()

	s.dirQ.
		On("Publish", mock.Anything, mock.MatchedBy(func(f *t.AnnotatedResource) bool {
			return s.Equal(*f, dirEntry)
		}), mock.AnythingOfType("uint8")).
		Return(nil).
		Once()

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.NoError(err)
	s.assertExpectations()
}

func TestCrawlerTestSuite(t *testing.T) {
	suite.Run(t, new(CrawlerTestSuite))
}
