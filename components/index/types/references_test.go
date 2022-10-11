package types

import (
	"log"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/suite"
)

var testRefs References = References{
	Reference{
		ParentHash: "QmSKboVigcD3AY4kLsob117KJcMHvMUu6vNFqk1PQzYUpp",
		Name:       "reference1",
	},
	Reference{
		ParentHash: "QmafrLBfzRLV4XSH1XcaMMeaXEUhDJjmtDfsYU95TrWG87",
		Name:       "reference2",
	},
	Reference{
		ParentHash: "sdfsdfsdfsd",
		Name:       "thrd",
	},
	Reference{
		ParentHash: "sdfsdfsdfsd",
		Name:       "thrd",
	},
	Reference{
		ParentHash: "sdfsdfsdfsd",
		Name:       "thrd",
	},
	Reference{
		ParentHash: "sdfsdfsdfsd",
		Name:       "thrd",
	},
}

type ReferencesTestSuite struct {
	suite.Suite
}

func (s *ReferencesTestSuite) TestMarshallUnmarshallBinary() {
	data, err := testRefs.MarshalBinary()
	s.NoError(err)
	s.NotEmpty(data)

	log.Printf("References serialised to %d bytes", len(data))

	// Check if we can back our original data
	newRefs := References{}
	err = newRefs.UnmarshalBinary(data)
	s.NoError(err)

	s.Equal(testRefs, newRefs)

	// Look at the raw stuff
	raw := new(interface{})
	err = cbor.Unmarshal(data, raw)
	s.NoError(err)

	log.Printf("%v", *raw)

	// json, err := json.MarshalIndent(raw, "", "  ")
	// s.NoError(err)
	// log.Printf("%s", json)
}

func TestReferencesTestSuite(t *testing.T) {
	suite.Run(t, new(ReferencesTestSuite))
}
