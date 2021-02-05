package ipfs

import (
	"net/http"

	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
)

const gatewayURL = "http://ipfs:8080"

type GatewayURLTestSuite struct {
	suite.Suite

	ipfs *IPFS
}

func (s *GatewayURLTestSuite) SetupTest() {
	cfg := DefaultConfig()
	cfg.GatewayURL = gatewayURL

	s.ipfs = New(cfg, http.DefaultClient, instr.New())
}

func (s *GatewayURLTestSuite) TestGatewayURLReferenced() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p",
		},
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
			},
			Name: "fileName.pdf",
		},
	}

	url := s.ipfs.GatewayURL(r)

	s.Equal(url, gatewayURL+"/ipfs/QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp/fileName.pdf")
}

func (s *GatewayURLTestSuite) TestGatewayURLUnreferenced() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p",
		},
	}

	url := s.ipfs.GatewayURL(r)

	s.Equal(url, gatewayURL+"/ipfs/QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p")
}

func (s *GatewayURLTestSuite) TestGatewayURLUnnamedReference() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p",
		},
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
			},
			Name: "",
		},
	}

	url := s.ipfs.GatewayURL(r)

	s.Equal(url, gatewayURL+"/ipfs/QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p")
}

func (s *GatewayURLTestSuite) TestEscapeURL() {
	// Regression test:
	// http://ipfs-tika:8081/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp/Killing_Yourself_to_Live:_85%_of_a_True_Story.html
	// panic: creating request: parse http://ipfs-tika:8081/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp/Killing_Yourself_to_Live:_85%_of_a_True_Story.html: invalid URL escape "%_o"

	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p",
		},
		Reference: t.Reference{
			Parent: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       "QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp",
			},
			Name: "Killing_Yourself_to_Live:_85%_of_a_True_Story.html",
		},
	}

	url := s.ipfs.GatewayURL(r)

	s.Equal(url, gatewayURL+"/ipfs/QmehSxmTPRCr85Xjgzjut6uWQihoTfqg9VVihJ892bmZCp/Killing_Yourself_to_Live:_85%25_of_a_True_Story.html")
}

func TestGatewayURLTestSuite(t *testing.T) {
	suite.Run(t, new(GatewayURLTestSuite))
}
