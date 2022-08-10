package nsfw

import (
	"context"
	"net/http"
	"testing"

	"github.com/dankinder/httpmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/ipfs-search/ipfs-search/components/extractor"
	indexTypes "github.com/ipfs-search/ipfs-search/components/index/types"

	"github.com/ipfs-search/ipfs-search/instr"
	"github.com/ipfs-search/ipfs-search/utils"
	t "github.com/ipfs-search/ipfs-search/types"
)

const testCID = "QmehHHRh1a7u66r7fugebp6f6wGNMGCa7eho9cgjwhAcm2"

type NSFWTestSuite struct {
	suite.Suite

	ctx context.Context
	e   extractor.Extractor
	getter utils.HTTPBodyGetter

	cfg *Config

	mockAPIHandler *httpmock.MockHandler
	mockAPIServer  *httpmock.Server
	responseHeader http.Header
}

func (s *NSFWTestSuite) SetupTest() {
	s.ctx = context.Background()

	s.mockAPIHandler = &httpmock.MockHandler{}
	s.mockAPIServer = httpmock.NewServer(s.mockAPIHandler)
	s.responseHeader = http.Header{
		"Content-Type": []string{"application/json"},
	}

	s.cfg = DefaultConfig()
	s.cfg.NSFWServerURL = s.mockAPIServer.URL()

	i := instr.New()
	s.getter = utils.NewHTTPBodyGetter(http.DefaultClient, i)

	s.e = New(s.cfg, s.getter, i)
}

func (s *NSFWTestSuite) TearDownTest() {
	s.mockAPIServer.Close()
}

func (s *NSFWTestSuite) TestExtract() {
	testJSON := []byte(`
		{
		  "classification": {
		    "neutral": 0.9980410933494568,
		    "drawing": 0.001135041005909443,
		    "porn": 0.00050011818530038,
		    "hentai": 0.00016194644558709115,
		    "sexy": 0.00016178081568796188
		  },
		  "nsfwjsVersion": "2.4.1"
		}
    `)

	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       testCID,
		},
		Stat: t.Stat{
			Size: 400,
		},
	}

	extractorURL := "/classify/" + testCID

	s.mockAPIHandler.
		On("Handle", "GET", extractorURL, mock.Anything).
		Return(httpmock.Response{
			Body: testJSON,
		}).
		Once()

	f := indexTypes.File{
		Metadata: indexTypes.Metadata{
			"Content-Type": "image/bmp",
		},
	}

	err := s.e.Extract(s.ctx, r, &f)

	s.NoError(err)
	s.mockAPIHandler.AssertExpectations(s.T())

	s.NotNil(f.NSFW)
	s.Equal(0.9980410933494568, f.NSFW.Classification.Neutral)
	s.Equal(0.001135041005909443, f.NSFW.Classification.Drawing)
	s.Equal(0.00050011818530038, f.NSFW.Classification.Porn)
	s.Equal(0.00016194644558709115, f.NSFW.Classification.Hentai)
	s.Equal(0.00016178081568796188, f.NSFW.Classification.Sexy)
	s.Equal("2.4.1", f.NSFW.NSFWVersion)
}

func (s *NSFWTestSuite) TestExtractMaxFileSize() {
	s.cfg.MaxFileSize = 100
	s.e = New(s.cfg, s.getter, instr.New())

	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       testCID,
		},
		Stat: t.Stat{
			Size: uint64(s.cfg.MaxFileSize + 1),
		},
	}

	f := indexTypes.File{
		Metadata: indexTypes.Metadata{
			"Content-Type": "image/jpeg",
		},
	}
	err := s.e.Extract(s.ctx, r, &f)

	s.Error(err, extractor.ErrFileTooLarge)
	s.mockAPIHandler.AssertExpectations(s.T())

	s.Nil(f.NSFW)
}

func (s *NSFWTestSuite) TestExtractUpstreamError() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       testCID,
		},
	}

	// Closing server early, generates a request error.
	s.mockAPIServer.Close()

	f := indexTypes.File{
		Metadata: indexTypes.Metadata{
			"Content-Type": "image/jpeg",
		},
	}

	err := s.e.Extract(s.ctx, r, &f)
	s.Error(err, t.ErrRequest)

	s.Nil(f.NSFW)
}

func (s *NSFWTestSuite) TestServer500() {
	// 500 will just propagate whatever error we're getting from a lower level
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       testCID,
		},
	}

	extractorURL := "/classify/" + testCID

	s.mockAPIHandler.
		On("Handle", "GET", extractorURL, mock.Anything).
		Return(httpmock.Response{
			Status: 500,
			Body:   []byte("{}"),
		}).
		Once()

	f := indexTypes.File{
		Metadata: indexTypes.Metadata{
			"Content-Type": "image/jpeg",
		},
	}

	err := s.e.Extract(s.ctx, r, &f)

	s.Error(err, t.ErrUnexpectedResponse)
	s.mockAPIHandler.AssertExpectations(s.T())

	s.Nil(f.NSFW)
}

func (s *NSFWTestSuite) TestExtractInvalidJSON() {
	testJSON := []byte(`invalid JSON`)

	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       testCID,
		},
		Stat: t.Stat{
			Size: 400,
		},
	}

	extractorURL := "/classify/" + testCID

	s.mockAPIHandler.
		On("Handle", "GET", extractorURL, mock.Anything).
		Return(httpmock.Response{
			Body: testJSON,
		}).
		Once()

	f := indexTypes.File{
		Metadata: indexTypes.Metadata{
			"Content-Type": "image/jpeg",
		},
	}

	err := s.e.Extract(s.ctx, r, &f)

	s.Error(err, t.ErrUnexpectedResponse)
	s.mockAPIHandler.AssertExpectations(s.T())

	s.Nil(f.NSFW)
}

func (s *NSFWTestSuite) TestIncompatibleType() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       testCID,
		},
		Stat: t.Stat{
			Size: 400,
		},
	}

	f := indexTypes.File{
		Metadata: indexTypes.Metadata{
			"Content-Type": "image/unsupported",
		},
	}

	err := s.e.Extract(s.ctx, r, &f)

	s.NoError(err)
	s.mockAPIHandler.AssertExpectations(s.T())

	s.Nil(f.NSFW)
}

func TestNSFWTestSuite(t *testing.T) {
	suite.Run(t, new(NSFWTestSuite))
}
