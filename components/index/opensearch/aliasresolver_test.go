package opensearch

import (
	"context"
	"testing"

	"github.com/dankinder/httpmock"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/stretchr/testify/suite"

	"github.com/ipfs-search/ipfs-search/components/index/opensearch/testsuite"
)

type AliasResolverTestSuite struct {
	testsuite.Suite

	ctx context.Context

	r *defaultAliasResolver
}

func (s *AliasResolverTestSuite) SetupTest() {
	s.SetupSearchMock()

	s.ctx = context.Background()

	client, _ := opensearch.NewClient(opensearch.Config{
		Addresses: []string{s.MockAPIServer.URL()},
	})

	s.r = &defaultAliasResolver{
		client: client,
	}
}

func (s *AliasResolverTestSuite) TeardownTest() {
	s.TeardownSearchMock()
}

func (s *AliasResolverTestSuite) TestRefreshAliases() {
	s.MockAPIHandler.
		On("Handle", "GET", "/_alias", []byte{}).
		Return(httpmock.Response{
			Body: []byte(`{
			"index1": {
				"aliases": {
					"index1_alias": {}
				}
			},
			"index2": {
				"aliases": {
					"index2_alias": {}
				}
			}
		}`),
		}).
		Once()

	err := s.r.refreshAliases(s.ctx)
	s.NoError(err)

	s.Equal(map[string]string{
		"index1": "index1_alias",
		"index2": "index2_alias",
	}, s.r.indexToAlias)

	s.Equal(map[string]string{
		"index1_alias": "index1",
		"index2_alias": "index2",
	}, s.r.aliasToIndex)
}

func (s *AliasResolverTestSuite) TestGetIndex() {
	s.MockAPIHandler.
		On("Handle", "GET", "/_alias", []byte{}).
		Return(httpmock.Response{
			Body: []byte(`{
			"index1": {
				"aliases": {
					"index1_alias": {}
				}
			},
			"index2": {
				"aliases": {
					"index2_alias": {}
				}
			}
		}`),
		}).
		Once()

	alias := "index1_alias"
	index, err := s.r.GetIndex(s.ctx, alias)
	s.NoError(err)
	s.Equal("index1", index)
}

func (s *AliasResolverTestSuite) TestGetIndexNotFound() {
	s.MockAPIHandler.
		On("Handle", "GET", "/_alias", []byte{}).
		Return(httpmock.Response{
			Body: []byte(`{
			"index1": {
				"aliases": {
					"index1_alias": {}
				}
			},
			"index2": {
				"aliases": {
					"index2_alias": {}
				}
			}
		}`),
		}).
		Once()

	alias := "index3_alias"

	_, err := s.r.GetIndex(s.ctx, alias)
	s.Error(err)
	s.Equal(err, ErrNotFound)
}

func (s *AliasResolverTestSuite) TestGetIndexUnknownRefresh() {
	// Execute TestGetIndex to populate cache.
	s.TestGetIndex()

	// Expect refresh, returning a previously unknown index
	s.MockAPIHandler.
		On("Handle", "GET", "/_alias", []byte{}).
		Return(httpmock.Response{
			Body: []byte(`{
			"index1": {
				"aliases": {
					"index1_alias": {}
				}
			},
			"index2": {
				"aliases": {
					"index2_alias": {}
				}
			},
			"index3": {
				"aliases": {
					"index3_alias": {}
				}
			}
		}`),
		}).
		Once()

	alias := "index3_alias"

	index, err := s.r.GetIndex(s.ctx, alias)
	s.NoError(err)
	s.Equal("index3", index)
}

func (s *AliasResolverTestSuite) TestGetAlias() {
	s.MockAPIHandler.
		On("Handle", "GET", "/_alias", []byte{}).
		Return(httpmock.Response{
			Body: []byte(`{
			"index1": {
				"aliases": {
					"index1_alias": {}
				}
			},
			"index2": {
				"aliases": {
					"index2_alias": {}
				}
			}
		}`),
		}).
		Once()

	index := "index1"
	alias, err := s.r.GetAlias(s.ctx, index)
	s.NoError(err)
	s.Equal("index1_alias", alias)
}

func (s *AliasResolverTestSuite) TestGetAliasNotFound() {
	s.MockAPIHandler.
		On("Handle", "GET", "/_alias", []byte{}).
		Return(httpmock.Response{
			Body: []byte(`{
			"index1": {
				"aliases": {
					"index1_alias": {}
				}
			},
			"index2": {
				"aliases": {
					"index2_alias": {}
				}
			}
		}`),
		}).
		Once()

	index := "index3"

	_, err := s.r.GetAlias(s.ctx, index)
	s.Error(err)
	s.Equal(err, ErrNotFound)
}

func (s *AliasResolverTestSuite) TestGetAliasUnknownRefresh() {
	// Execute TestGetAlias to populate cache.
	s.TestGetAlias()

	// Expect refresh, returning a previously unknown alias
	s.MockAPIHandler.
		On("Handle", "GET", "/_alias", []byte{}).
		Return(httpmock.Response{
			Body: []byte(`{
			"index1": {
				"aliases": {
					"index1_alias": {}
				}
			},
			"index2": {
				"aliases": {
					"index2_alias": {}
				}
			},
			"index3": {
				"aliases": {
					"index3_alias": {}
				}
			}
		}`),
		}).
		Once()

	index := "index3"

	alias, err := s.r.GetAlias(s.ctx, index)
	s.NoError(err)
	s.Equal("index3_alias", alias)
}

func TestAliasResolverTestSuite(t *testing.T) {
	suite.Run(t, new(AliasResolverTestSuite))
}
