package crawler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/ipfs-search/ipfs-search/components/extractor"
	"github.com/ipfs-search/ipfs-search/components/index"
	indexTypes "github.com/ipfs-search/ipfs-search/components/index/types"
	"github.com/ipfs-search/ipfs-search/components/protocol"
	"github.com/ipfs-search/ipfs-search/components/queue"

	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
)

type CrawlerTestSuite struct {
	suite.Suite

	ctx     context.Context
	cfg     *Config
	indexes *Indexes
	queues  *Queues
	c       *Crawler
	instr   *instr.Instrumentation

	protocol  *protocol.Mock
	extractor *extractor.Mock

	fileIdx    *index.Mock
	dirIdx     *index.Mock
	invalidIdx *index.Mock
	partialIdx *index.Mock

	dirQ  *queue.Mock
	fileQ *queue.Mock
	hashQ *queue.Mock
}

func (s *CrawlerTestSuite) SetupTest() {
	s.ctx = context.Background()

	// Creat a crawler with mocked dependencies
	s.fileIdx, s.dirIdx, s.invalidIdx, s.partialIdx = &index.Mock{}, &index.Mock{}, &index.Mock{}, &index.Mock{}

	s.indexes = &Indexes{
		Files:       s.fileIdx,
		Directories: s.dirIdx,
		Invalids:    s.invalidIdx,
		Partials:    s.partialIdx,
	}

	s.fileQ, s.dirQ, s.hashQ = &queue.Mock{}, &queue.Mock{}, &queue.Mock{}

	s.queues = &Queues{
		Directories: s.dirQ,
		Files:       s.fileQ,
		Hashes:      s.hashQ,
	}
	s.protocol = &protocol.Mock{}
	s.extractor = &extractor.Mock{}

	s.instr = instr.New()

	s.cfg = DefaultConfig()

	s.c = New(s.cfg, s.indexes, s.queues, s.protocol, s.extractor, s.instr)
}

func (s *CrawlerTestSuite) assertExpectations() {
	mock.AssertExpectationsForObjects(s.T(),
		s.fileIdx,
		s.dirIdx,
		s.invalidIdx,
		s.fileQ,
		s.dirQ,
		s.hashQ,
		s.protocol,
		s.extractor,
	)
}

func (s *CrawlerTestSuite) assertNotExists(rID string) {
	s.fileIdx.
		On("Get", mock.Anything, rID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Once()

	s.dirIdx.
		On("Get", mock.Anything, rID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Once()

	s.invalidIdx.
		On("Get", mock.Anything, rID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Once()

	s.partialIdx.
		On("Get", mock.Anything, rID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Once()
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

	s.assertNotExists(r.Resource.ID)

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
			Type: t.UndefinedType,
		},
	}

	// Mock assertions
	s.protocol.
		On("Stat", mock.Anything, r).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*t.AnnotatedResource)
			r.Stat = t.Stat{
				Type: t.UnsupportedType,
			}
		}).
		Return(nil).
		Once()

	// Mock assertions
	s.invalidIdx.
		On("Index", mock.Anything, r.Resource.ID, mock.MatchedBy(func(f *indexTypes.Invalid) bool {
			return s.Equal(f.Error, "unsupported type")
		})).
		Return(nil).
		Once()

	s.assertNotExists(r.Resource.ID)

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Undefined types should index as invalid.
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlUnexpectedType() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Stat: t.Stat{
			Type: 254,
		},
	}

	// Unexpected types yield a panic; undefined behaviour
	s.Panics(func() { _ = s.c.Crawl(s.ctx, r) })
}

func (s *CrawlerTestSuite) TestCrawlInvalid() {
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

	invalidErr := fmt.Errorf("%w: %s", t.ErrInvalidResource, "test error")

	// Mock assertions
	s.protocol.
		On("Stat", mock.Anything, r).
		Return(invalidErr).
		Once()

	// Mock assertions
	s.invalidIdx.
		On("Index", mock.Anything, r.Resource.ID, mock.MatchedBy(func(f *indexTypes.Invalid) bool {
			return s.Equal(f.Error, invalidErr.Error())
		})).
		Return(nil).
		Once()

	s.assertNotExists(r.Resource.ID)

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
			Type: t.UndefinedType,
		},
	}

	s.protocol.
		On("Stat", mock.Anything, r).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*t.AnnotatedResource)
			r.Stat = t.Stat{
				Type: t.PartialType,
			}
		}).
		Return(nil).
		Once()

	s.assertNotExists(r.Resource.ID)

	// Index as partial
	s.partialIdx.
		On("Index", mock.Anything, r.Resource.ID, mock.MatchedBy(func(f *indexTypes.Partial) bool {
			return true
		})).
		Return(nil).
		Once()

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Unreferenced partial should be skipped.
	s.NoError(err)
	s.assertExpectations()
}

// TestCrawlReferencedPartialType crawls a cached partial, now with a reference.
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
			Name: "referencedPartial.pdf",
		},
		Stat: t.Stat{
			Type: t.FileType,
		},
	}

	s.extractor.
		On("Extract", mock.Anything, r, mock.Anything).
		Return(nil).
		Once()

	s.fileIdx.
		On("Index", mock.Anything, r.Resource.ID, mock.MatchedBy(func(f *indexTypes.File) bool {
			return true
		})).
		Return(nil).
		Once()

	s.fileIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.dirIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.invalidIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.partialIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(true, nil).
		Once()

	s.partialIdx.
		On("Delete", mock.Anything, r.Resource.ID).
		Return(nil).
		Once()

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

	s.assertNotExists(r.Resource.ID)

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlLargeFile() {
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

	largeFileErr := fmt.Errorf("blabla %w", extractor.ErrFileTooLarge)

	s.extractor.
		On("Extract", mock.Anything, r, mock.Anything).
		Return(largeFileErr).
		Once()

	s.invalidIdx.
		On("Index", mock.Anything, r.Resource.ID, mock.MatchedBy(func(f *indexTypes.Invalid) bool {
			return s.Equal("resource invalid: blabla file too large", f.Error)
		})).
		Return(nil).
		Once()

	s.assertNotExists(r.Resource.ID)

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlStatTimeout() {
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

	s.protocol.
		On("Stat", mock.Anything, r).
		Return(context.DeadlineExceeded).
		Once()

	s.assertNotExists(r.Resource.ID)

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	s.Equal(err, context.DeadlineExceeded)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlReferencedFile() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
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
		},
	}

	s.extractor.
		On("Extract", mock.Anything, r, mock.Anything).
		Return(nil).
		Once()

	s.fileIdx.
		On("Index", mock.Anything, r.Resource.ID, mock.MatchedBy(func(f *indexTypes.File) bool {
			return s.Equal(f.References, indexTypes.References{
				indexTypes.Reference{
					ParentHash: r.Reference.Parent.ID,
					Name:       r.Reference.Name,
				},
			})
		})).
		Return(nil).
		Once()

	s.assertNotExists(r.Resource.ID)

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlReferencedDirectory() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
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
		},
	}

	// Empty dir
	s.protocol.
		On("Ls", mock.Anything, r, mock.AnythingOfType("chan<- *types.AnnotatedResource")).
		Return(nil).
		Once()

	s.dirIdx.
		On("Index", mock.Anything, r.Resource.ID, mock.MatchedBy(func(f *indexTypes.Directory) bool {
			return s.Equal(f.References, indexTypes.References{
				indexTypes.Reference{
					ParentHash: r.Reference.Parent.ID,
					Name:       r.Reference.Name,
				},
			})
		})).
		Return(nil).
		Once()

	s.assertNotExists(r.Resource.ID)

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

	unsupportedEntry := t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv",
		},
		Stat: t.Stat{
			Type: t.UnsupportedType,
		},
	}

	unknownEntry := t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv",
		},
		Stat: t.Stat{
			Type: t.UndefinedType,
		},
	}

	s.protocol.
		On("Ls", mock.Anything, r, mock.AnythingOfType("chan<- *types.AnnotatedResource")).
		Run(func(args mock.Arguments) {
			// Write bogus directory entry
			entryChan := args.Get(2).(chan<- *t.AnnotatedResource)
			entryChan <- &fileEntry
			entryChan <- &dirEntry
			entryChan <- &unsupportedEntry
			entryChan <- &unknownEntry
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
					indexTypes.Link{
						Hash: unsupportedEntry.ID,
						Name: unsupportedEntry.Reference.Name,
						Size: unsupportedEntry.Size,
						Type: indexTypes.UnsupportedLinkType,
					},
					indexTypes.Link{
						Hash: unknownEntry.ID,
						Name: unknownEntry.Reference.Name,
						Size: unknownEntry.Size,
						Type: indexTypes.UnknownLinkType,
					},
				})
		})).
		Return(nil).
		Once()

	s.invalidIdx.
		On("Index", mock.Anything, unsupportedEntry.ID, mock.MatchedBy(func(f *indexTypes.Invalid) bool {
			return s.Equal("unsupported type", f.Error)
		})).
		Return(nil).
		Once()

	s.fileQ.
		On("Publish", mock.Anything, mock.MatchedBy(func(f *t.AnnotatedResource) bool {
			return s.Equal(fileEntry, *f)
		}), mock.AnythingOfType("uint8")).
		Return(nil).
		Once()

	s.dirQ.
		On("Publish", mock.Anything, mock.MatchedBy(func(f *t.AnnotatedResource) bool {
			return s.Equal(dirEntry, *f)
		}), mock.AnythingOfType("uint8")).
		Return(nil).
		Once()

	s.hashQ.
		On("Publish", mock.Anything, mock.MatchedBy(func(f *t.AnnotatedResource) bool {
			return s.Equal(unknownEntry, *f)
		}), mock.AnythingOfType("uint8")).
		Return(nil).
		Once()

	s.assertNotExists(r.Resource.ID)

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlDirectoryUnexpectedType() {
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

	unexpectedEntry := t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv",
		},
		Stat: t.Stat{
			Type: 255,
		},
	}

	s.protocol.
		On("Ls", mock.Anything, r, mock.AnythingOfType("chan<- *types.AnnotatedResource")).
		Run(func(args mock.Arguments) {
			// Write bogus directory entry
			entryChan := args.Get(2).(chan<- *t.AnnotatedResource)
			entryChan <- &unexpectedEntry
		}).
		Return(nil).
		Once()

	s.assertNotExists(r.Resource.ID)

	//// THIS PANIC IS NOT PROPERLY CAUGHT FIXME!!!
	// The issue is the panic is in a different goroutine and does not propagate to the parent.

	// Unexpected types yield a panic; undefined behaviour
	s.Panics(func() { _ = s.c.Crawl(s.ctx, r) })
}

func (s *CrawlerTestSuite) TestCrawlLargeDirectory() {
	s.cfg = DefaultConfig()

	// Override MaxDirSize
	s.cfg.MaxDirSize = 3

	s.c = New(s.cfg, s.indexes, s.queues, s.protocol, s.extractor, s.instr)

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

	s.protocol.
		On("Ls", mock.Anything, r, mock.AnythingOfType("chan<- *types.AnnotatedResource")).
		Run(func(args mock.Arguments) {
			// Write bogus directory entry
			entryChan := args.Get(2).(chan<- *t.AnnotatedResource)
			entryChan <- &fileEntry
			entryChan <- &fileEntry
			entryChan <- &fileEntry
			entryChan <- &fileEntry
			entryChan <- &fileEntry
		}).
		Return(nil).
		Once()

	s.invalidIdx.
		On("Index", mock.Anything, r.Resource.ID, mock.MatchedBy(func(f *indexTypes.Invalid) bool {
			return s.Equal("directory too large", f.Error)
		})).
		Return(nil).
		Once()

	s.fileQ.
		On("Publish", mock.Anything, mock.MatchedBy(func(f *t.AnnotatedResource) bool {
			return s.Equal(fileEntry, *f)
		}), mock.AnythingOfType("uint8")).
		Return(nil).
		Times(5)

	s.assertNotExists(r.Resource.ID)

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlDirEntryTimeout() {
	s.cfg = DefaultConfig()

	// Override dir entry timeout
	s.cfg.DirEntryTimeout = 5 * time.Millisecond

	s.c = New(s.cfg, s.indexes, s.queues, s.protocol, s.extractor, s.instr)

	entryDelay := 2 * s.cfg.DirEntryTimeout

	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Stat: t.Stat{
			Type: t.DirectoryType,
		},
	}

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

	s.fileQ.
		On("Publish", mock.Anything, mock.MatchedBy(func(f *t.AnnotatedResource) bool {
			return s.Equal(*f, fileEntry)
		}), mock.AnythingOfType("uint8")).
		Return(nil).
		Once()

	s.protocol.
		On("Ls", mock.Anything, r, mock.AnythingOfType("chan<- *types.AnnotatedResource")).
		Run(func(args mock.Arguments) {
			entryChan := args.Get(2).(chan<- *t.AnnotatedResource)
			entryChan <- &fileEntry
			time.Sleep(entryDelay)
			entryChan <- &dirEntry
		}).
		Return(nil).
		Once()

	s.assertNotExists(r.Resource.ID)

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	s.Equal(err, context.DeadlineExceeded)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlUpdateLastSeen() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
	}

	// File is found, last seen 1 hour
	s.fileIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Run(func(args mock.Arguments) {
			u := args.Get(2).(*indexTypes.Update)
			lastSeen := time.Now().Add(-2 * time.Hour)
			u.LastSeen = &lastSeen
		}).
		Return(true, nil).
		Once()

	s.dirIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.invalidIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.partialIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.fileIdx.
		On("Update", mock.Anything, r.Resource.ID, mock.MatchedBy(func(u *indexTypes.Update) bool {
			return s.Empty(u.References) &&
				s.WithinDuration(*u.LastSeen, time.Now(), time.Second)
		})).
		Return(nil).
		Once()

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlNotUpdateInvalid() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
	}

	// File is found, last seen 1 hour
	s.fileIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Once()

	s.dirIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.partialIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.invalidIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(true, nil).
		Maybe()

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlAddReference() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Source: t.DirectorySource,
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8",
			},
			Name: "NewReference.pdf",
		},
	}

	// File is found, very recently, but a new reference is found.
	s.fileIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Run(func(args mock.Arguments) {
			u := args.Get(2).(*indexTypes.Update)
			lastSeen := time.Now()
			u.LastSeen = &lastSeen
			u.References = indexTypes.References{
				indexTypes.Reference{
					ParentHash: "Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
					Name:       "ExistingReference.pdf",
				},
			}
		}).
		Return(true, nil).
		Once()

	s.dirIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.invalidIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.partialIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.fileIdx.
		On("Update", mock.Anything, r.Resource.ID, mock.MatchedBy(func(u *indexTypes.Update) bool {
			return s.ElementsMatch(u.References, indexTypes.References{
				indexTypes.Reference{
					ParentHash: "Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
					Name:       "ExistingReference.pdf",
				},
				indexTypes.Reference{
					ParentHash: "QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8",
					Name:       "NewReference.pdf",
				},
			}) && u.LastSeen == nil
		})).
		Return(nil).
		Once()

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlUpdateGetError() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Source: t.DirectorySource,
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8",
			},
			Name: "NewReference.pdf",
		},
	}

	testErr := errors.New("test")

	s.fileIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, testErr).
		Maybe()

	s.dirIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.partialIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.invalidIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.Equal(testErr, err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlUpdateUpdateError() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Source: t.DirectorySource,
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8",
			},
			Name: "NewReference.pdf",
		},
	}

	// File is found, very recently, but a new reference is found.
	s.fileIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Run(func(args mock.Arguments) {
			u := args.Get(2).(*indexTypes.Update)
			lastSeen := time.Now()
			u.LastSeen = &lastSeen
			u.References = indexTypes.References{
				indexTypes.Reference{
					ParentHash: "Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
					Name:       "ExistingReference.pdf",
				},
			}
		}).
		Return(true, nil).
		Once()

	testErr := errors.New("test")

	s.dirIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.invalidIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.partialIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.fileIdx.
		On("Update", mock.Anything, r.Resource.ID, mock.MatchedBy(func(u *indexTypes.Update) bool {
			return s.ElementsMatch(u.References, indexTypes.References{
				indexTypes.Reference{
					ParentHash: "Qmc8mmzycvXnzgwBHokZQd97iWAmtdFMqX4FZUAQ5AQdQi",
					Name:       "ExistingReference.pdf",
				},
				indexTypes.Reference{
					ParentHash: "QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8",
					Name:       "NewReference.pdf",
				},
			}) && u.LastSeen == nil
		})).
		Return(testErr).
		Once()

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.Equal(testErr, err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlSameReference() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Source: t.DirectorySource,
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8",
			},
			Name: "NewReference.pdf",
		},
	}

	// File is found, very recently, but a new reference is found.
	s.fileIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Run(func(args mock.Arguments) {
			u := args.Get(2).(*indexTypes.Update)
			lastSeen := time.Now()
			u.LastSeen = &lastSeen
			u.References = indexTypes.References{
				indexTypes.Reference{
					ParentHash: "QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8",
					Name:       "NewReference.pdf",
				},
			}
		}).
		Return(true, nil).
		Once()

	s.dirIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.partialIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.invalidIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	// Crawl
	err := s.c.Crawl(s.ctx, r)

	// Test result, side effects
	s.NoError(err)
	s.assertExpectations()
}

func (s *CrawlerTestSuite) TestCrawlSameReferenceDifferentName() {
	// Prepare resource
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		},
		Source: t.DirectorySource,
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8",
			},
			Name: "NewReference.pdf",
		},
	}

	// File is found, very recently, but a new reference is found.
	s.fileIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Run(func(args mock.Arguments) {
			u := args.Get(2).(*indexTypes.Update)
			lastSeen := time.Now()
			u.LastSeen = &lastSeen
			u.References = indexTypes.References{
				indexTypes.Reference{
					ParentHash: "QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8",
					Name:       "NewName.pdf",
				},
			}
		}).
		Return(true, nil).
		Once()

	s.dirIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.partialIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.invalidIdx.
		On("Get", mock.Anything, r.Resource.ID, mock.Anything, []string{"references", "last-seen"}).
		Return(false, nil).
		Maybe()

	s.fileIdx.
		On("Update", mock.Anything, r.Resource.ID, mock.MatchedBy(func(u *indexTypes.Update) bool {
			return s.ElementsMatch(u.References, indexTypes.References{
				indexTypes.Reference{
					ParentHash: "QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8",
					Name:       "NewName.pdf",
				},
				indexTypes.Reference{
					ParentHash: "QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8",
					Name:       "NewReference.pdf",
				},
			}) &&
				u.LastSeen == nil
		})).
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
