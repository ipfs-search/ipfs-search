package proxy

import (
	"testing"

	"github.com/ipfs/go-datastore"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AfterPutMock struct {
	mock.Mock
}

func (m *AfterPutMock) AfterPut(k datastore.Key, v []byte, err error) error {
	return m.Called(k, v, err).Error(0)
}

type DataStoreTestSuite struct {
	suite.Suite
	mock *AfterPutMock
	ds   datastore.Batching
}

func (s *DataStoreTestSuite) SetupTest() {
	s.mock = &AfterPutMock{}
	s.mock.Test(s.T())

	s.ds = datastore.NewMapDatastore()
}

func (s *DataStoreTestSuite) TearDownTest() {
	s.ds.Close()
}

func (s *DataStoreTestSuite) TestPut() {
	ds := New(s.ds, s.mock.AfterPut)

	k := datastore.NewKey("test")
	v := []byte("test")

	s.mock.On("AfterPut", k, v, nil).Return(nil)

	err := ds.Put(k, v)

	s.NoError(err)
	s.mock.AssertExpectations(s.T())
}

func (s *DataStoreTestSuite) TestBatch() {
	ds := New(s.ds, s.mock.AfterPut)

	k := datastore.NewKey("test")
	v := []byte("test")

	s.mock.On("AfterPut", k, v, nil).Return(nil)

	b, err := ds.Batch()
	s.NoError(err)

	err = b.Put(k, v)
	s.NoError(err)

	err = b.Commit()
	s.NoError(err)

	s.mock.AssertExpectations(s.T())
}

func TestDataStoreTestSuite(t *testing.T) {
	suite.Run(t, new(DataStoreTestSuite))
}
