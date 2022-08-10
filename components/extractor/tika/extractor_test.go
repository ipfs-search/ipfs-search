package tika

import (
    "context"
    "fmt"
    "net/http"
    "net/url"
    "testing"

    "github.com/dankinder/httpmock"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"

    "github.com/ipfs-search/ipfs-search/components/extractor"
    indexTypes "github.com/ipfs-search/ipfs-search/components/index/types"
    "github.com/ipfs-search/ipfs-search/components/protocol"

    "github.com/ipfs-search/ipfs-search/instr"
    t "github.com/ipfs-search/ipfs-search/types"
)

const testCID = "QmehHHRh1a7u66r7fugebp6f6wGNMGCa7eho9cgjwhAcm2"

type TikaTestSuite struct {
    suite.Suite

    ctx context.Context
    e   extractor.Extractor

    cfg      *Config
    protocol *protocol.Mock

    mockAPIHandler *httpmock.MockHandler
    mockAPIServer  *httpmock.Server
    responseHeader http.Header
}

func (s *TikaTestSuite) SetupTest() {
    s.ctx = context.Background()

    s.mockAPIHandler = &httpmock.MockHandler{}
    s.mockAPIServer = httpmock.NewServer(s.mockAPIHandler)
    s.responseHeader = http.Header{
        "Content-Type": []string{"application/json"},
    }

    s.cfg = DefaultConfig()
    s.cfg.TikaExtractorURL = s.mockAPIServer.URL()

    s.protocol = &protocol.Mock{}

    s.e = New(s.cfg, http.DefaultClient, s.protocol, instr.New())
}

func (s *TikaTestSuite) TearDownTest() {
    s.mockAPIServer.Close()
}

func (s *TikaTestSuite) TestExtract() {
    testJSON := []byte(`
        {
          "metadata": {
            "title": [
              "How Filecoin Supports Video Storage"
            ],
            "Content-Type": [
              "text/html; charset=UTF-8"
            ]
          },
          "content": "The Filecoin Space Race is now live! Learn More\n\t\t Thank you!",
          "language": {
            "language": "en",
            "confidence": "HIGH",
            "rawScore": 0.99999505
          },
          "urls": [
            "https://filecoin.io/uploads/video-storage-social.png",
            "https://proto.school/#/tutorials?course=filecoin"
          ],
          "ipfs_tika_version": "dev-build"
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

    gwURL := "http://localhost:8080/ipfs/" + testCID
    extractorURL := "/extract?url=http%3A%2F%2Flocalhost%3A8080%2Fipfs%2F" + testCID

    s.protocol.
        On("GatewayURL", r).
        Return(gwURL).
        Once()

    s.mockAPIHandler.
        On("Handle", "GET", extractorURL, mock.Anything).
        Return(httpmock.Response{
            Body: testJSON,
        }).
        Once()

    f := &indexTypes.File{
        Document: indexTypes.Document{
            Size: r.Size,
        },
    }

    err := s.e.Extract(s.ctx, r, &f)

    s.NoError(err)
    s.mockAPIHandler.AssertExpectations(s.T())

    s.Equal(uint64(400), f.Size)
    s.Equal([]interface{}{"How Filecoin Supports Video Storage"}, f.Metadata["title"])
    s.Equal("en", f.Language.Language)
    s.Equal(0.99999505, f.Language.RawScore)
    s.Contains(f.URLs, "https://proto.school/#/tutorials?course=filecoin")
}

func (s *TikaTestSuite) TestExtractMaxFileSize() {
    s.cfg.MaxFileSize = 100
    s.e = New(s.cfg, http.DefaultClient, s.protocol, instr.New())

    r := &t.AnnotatedResource{
        Resource: &t.Resource{
            Protocol: t.IPFSProtocol,
            ID:       testCID,
        },
        Stat: t.Stat{
            Size: uint64(s.cfg.MaxFileSize + 1),
        },
    }

    f := &indexTypes.File{}
    err := s.e.Extract(s.ctx, r, &f)

    s.Error(err, extractor.ErrFileTooLarge)
    s.mockAPIHandler.AssertExpectations(s.T())
}

func (s *TikaTestSuite) TestExtractUpstreamError() {
    r := &t.AnnotatedResource{
        Resource: &t.Resource{
            Protocol: t.IPFSProtocol,
            ID:       testCID,
        },
        Stat: t.Stat{
            Size: 400,
        },
    }

    gwURL := "http://localhost:8080/ipfs/%s" + testCID
    // extractorURL := fmt.Sprintf("/extract?url=%s", url.QueryEscape(gwURL))

    s.protocol.
        On("GatewayURL", r).
        Return(gwURL).
        Once()

    // Closing server early, generates a request error.
    s.mockAPIServer.Close()

    f := &indexTypes.File{
        Document: indexTypes.Document{
            Size: r.Size,
        },
    }

    err := s.e.Extract(s.ctx, r, &f)
    s.Error(err, extractor.ErrRequest)
}

func (s *TikaTestSuite) TestURLEscape() {
    // Regression test:
    // http://ipfs-tika:8081/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp/Killing_Yourself_to_Live:_85%_of_a_True_Story.html
    // panic: creating request: parse http://ipfs-tika:8081/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp/Killing_Yourself_to_Live:_85%_of_a_True_Story.html: invalid URL escape "%_o"

    tikaURL := "/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp/" + url.PathEscape("Killing_Yourself_to_Live:_85%_of_a_True_Story.html")
    gwURL := "http://localhost:8080" + tikaURL
    extractorURL := fmt.Sprintf("/extract?url=%s", url.QueryEscape(gwURL))

    r := &t.AnnotatedResource{
        Resource: &t.Resource{
            Protocol: t.IPFSProtocol,
            ID:       testCID,
        },
    }

    s.protocol.
        On("GatewayURL", r).
        Return(gwURL).
        Once()

    s.mockAPIHandler.
        On("Handle", "GET", extractorURL, mock.Anything).
        Return(httpmock.Response{
            Body: []byte("{}"),
        }).
        Once()

    f := &indexTypes.File{}

    err := s.e.Extract(s.ctx, r, &f)

    s.NoError(err)
    s.mockAPIHandler.AssertExpectations(s.T())
}

func (s *TikaTestSuite) TestTika500() {
    // 500 will just propagate whatever error we're getting from a lower level
    r := &t.AnnotatedResource{
        Resource: &t.Resource{
            Protocol: t.IPFSProtocol,
            ID:       testCID,
        },
    }

    gwURL := "http://localhost:8080/ipfs/%s" + testCID
    extractorURL := fmt.Sprintf("/extract?url=%s", url.QueryEscape(gwURL))

    s.protocol.
        On("GatewayURL", r).
        Return(gwURL).
        Once()

    s.mockAPIHandler.
        On("Handle", "GET", extractorURL, mock.Anything).
        Return(httpmock.Response{
            Status: 500,
            Body:   []byte("{}"),
        }).
        Once()

    f := &indexTypes.File{}

    err := s.e.Extract(s.ctx, r, &f)

    s.Error(err, extractor.ErrUnexpectedResponse)
    s.mockAPIHandler.AssertExpectations(s.T())
}

func (s *TikaTestSuite) TestExtractInvalidJSON() {
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

    gwURL := "http://localhost:8080/ipfs/%s" + testCID
    extractorURL := fmt.Sprintf("/extract?url=%s", url.QueryEscape(gwURL))

    s.protocol.
        On("GatewayURL", r).
        Return(gwURL).
        Once()

    s.mockAPIHandler.
        On("Handle", "GET", extractorURL, mock.Anything).
        Return(httpmock.Response{
            Body: testJSON,
        }).
        Once()

    f := &indexTypes.File{
        Document: indexTypes.Document{
            Size: r.Size,
        },
    }

    err := s.e.Extract(s.ctx, r, &f)

    s.Error(err, extractor.ErrUnexpectedResponse)
    s.mockAPIHandler.AssertExpectations(s.T())
}

func TestTikaTestSuite(t *testing.T) {
    suite.Run(t, new(TikaTestSuite))
}
