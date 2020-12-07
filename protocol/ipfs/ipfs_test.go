package ipfs

import (
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"

	"github.com/ipfs-search/ipfs-search/instr"
	t "github.com/ipfs-search/ipfs-search/types"
)

const gatewayURL = "http://ipfs:8080"

type IPFSTestSuite struct {
	suite.Suite

	ipfs *IPFS
}

func (s *IPFSTestSuite) SetupTest() {
	cfg := DefaultConfig()
	cfg.IPFSGatewayURL = gatewayURL

	s.ipfs = New(cfg, http.DefaultClient, instr.New())
}

func (s *IPFSTestSuite) TestGatewayURLReferenced() {
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

func (s *IPFSTestSuite) TestGatewayURLUnreferenced() {
	r := &t.AnnotatedResource{
		Resource: &t.Resource{
			Protocol: t.IPFSProtocol,
			ID:       "QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p",
		},
	}

	url := s.ipfs.GatewayURL(r)

	s.Equal(url, gatewayURL+"/ipfs/QmcBLKyRHjbGeLnjnmj74FFJpGJDz4YxFqUDYqMU7Mny1p")
}

func (s *IPFSTestSuite) TestGatewayURLUnnamedReference() {
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

func TestIPFSTestSuite(t *testing.T) {
	suite.Run(t, new(IPFSTestSuite))
}
