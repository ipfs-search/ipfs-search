package multi

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DefaultSelectorTestSuite struct {
	suite.Suite
	gi func(string, Properties) string
}

func (s *DefaultSelectorTestSuite) SetupTest() {
	s.gi = DefaultSelector.GetIndex
}

func (s *DefaultSelectorTestSuite) TestDocument() {
	exp := "ipfs_documents"
	s.Equal(exp, s.gi("", makeMetadataProp("image/vnd.djvu")))
	s.Equal(exp, s.gi("", makeMetadataProp("application/epub+zip")))     // Tricky character
	s.Equal(exp, s.gi("", makeMetadataProp("text/plain;charset=UTF-8"))) // Charset
}

func (s *DefaultSelectorTestSuite) TestImage() {
	exp := "ipfs_images"
	s.Equal(exp, s.gi("", makeMetadataProp("application/dicom")))
	s.Equal(exp, s.gi("", makeMetadataProp("image/crazybananas")))
}

func (s *DefaultSelectorTestSuite) TestData() {
	s.Equal("ipfs_data", s.gi("", makeMetadataProp("application/rdf+xml"))) // Tricky character
}

func (s *DefaultSelectorTestSuite) TestArchive() {
	s.Equal("ipfs_archives", s.gi("", makeMetadataProp("application/vnd.google-earth.kmz")))
}

func (s *DefaultSelectorTestSuite) TestUnknown() {
	s.Equal("ipfs_unknown", s.gi("", nil))
}

func (s *DefaultSelectorTestSuite) TestVideo() {
	exp := "ipfs_videos"
	s.Equal(exp, s.gi("", makeMetadataProp("application/x-matroska")))
	s.Equal(exp, s.gi("", makeMetadataProp("video/crazyvideo")))
}

func (s *DefaultSelectorTestSuite) TestAudio() {
	exp := "ipfs_audio"
	s.Equal(exp, s.gi("", makeMetadataProp("application/ogg")))
	s.Equal(exp, s.gi("", makeMetadataProp("audio/scream")))
}

func (s *DefaultSelectorTestSuite) TestOther() {
	exp := "ipfs_other"
	s.Equal(exp, s.gi("", makeMetadataProp("somethingelse")))
}

func TestDefaultSelectorTestSuite(t *testing.T) {
	suite.Run(t, new(DefaultSelectorTestSuite))
}
