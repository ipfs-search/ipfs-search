package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	// "github.com/stretchr/testify/mock"
	radix "github.com/mediocregopher/radix/v4"
	"github.com/stretchr/testify/suite"

	"github.com/ipfs-search/ipfs-search/components/index/types"
	"github.com/ipfs-search/ipfs-search/instr"
)

type ExistsIndexTestSuite struct {
	suite.Suite
	ctx   context.Context
	instr *instr.Instrumentation
	now   *time.Time

	i *ExistsIndex
}

func (s *ExistsIndexTestSuite) stubClient(fn stubFunc) *Client {
	conn := radix.NewStubConn("", "", fn)
	rClient := radix.NewMultiClient(radix.ReplicaSet{
		Primary: conn,
	})
	return &Client{
		cfg:             &ClientConfig{},
		Instrumentation: instr.New(),
		radixClient:     rClient,
	}
}

func (s *ExistsIndexTestSuite) stubIndex(fn stubFunc) *Index {
	client := s.stubClient(fn)
	cfg := &Config{
		Name:   indexName,
		Prefix: indexPrefix,
	}

	return &Index{
		cfg, client,
	}
}

func (s *ExistsIndexTestSuite) SetupTest() {
	s.ctx = context.Background()
	now := time.Now().Truncate(time.Second).UTC()
	s.now = &now
}

func (s *ExistsIndexTestSuite) TestGetKey() {
	i := s.stubIndex(func(_ context.Context, _ []string) interface{} {
		return nil
	})

	k := i.getKey(testId)
	s.Equal(indexPrefix+":"+testId, k)
}

func (s *ExistsIndexTestSuite) TestSetLastSeenOnly() {
	i := s.stubIndex(func(_ context.Context, args []string) interface{} {
		s.Len(args, 4)
		s.Equal("HSET", args[0])
		// We test key generation elsewhere, ignore args[1]
		s.Equal("l", args[2]) // Use the `redis` tag.
		s.Equal(s.now.Format(time.RFC3339), args[3])

		return nil
	})
	u := &types.Update{
		LastSeen: s.now,
	}

	err := i.set(s.ctx, testId, u)
	s.NoError(err)
}

func (s *ExistsIndexTestSuite) TestSetReferencesOnly() {
	r1 := types.Reference{
		ParentHash: "p1",
		Name:       "f1",
	}
	r2 := types.Reference{
		ParentHash: "p2",
		Name:       "f2",
	}
	r := types.References{
		r1, r2,
	}
	u := &types.Update{
		References: r,
	}

	i := s.stubIndex(func(_ context.Context, args []string) interface{} {
		s.Len(args, 4)
		s.Equal("HSET", args[0])
		// We test key generation elsewhere, ignore args[1]
		s.Equal("r", args[2]) // Use the `redis` tag.

		refs := types.References{}
		err := refs.UnmarshalBinary([]byte(args[3]))
		s.NoError(err)

		s.Equal(refs, u.References)

		return nil
	})

	err := i.set(s.ctx, testId, u)
	s.NoError(err)
}

// func (s *ExistsIndexTestSuite) TestEmpty() {
// 	u := &types.Update{}

// 	i := s.stubIndex(func(_ context.Context, args []string) interface{} {
// 		s.Len(args, 4)
// 		s.Equal("HSET", args[0])
// 		// We test key generation elsewhere, ignore args[1]

// 		// s.Equal("l", args[2]) // Use the `redis` tag.
// 		return nil
// 	})

// 	err := i.set(s.ctx, testId, u)
// 	s.NoError(err)
// }

func (s *ExistsIndexTestSuite) TestSetAll() {
	r1 := types.Reference{
		ParentHash: "p1",
		Name:       "f1",
	}
	r := types.References{
		r1,
	}
	u := &types.Update{
		LastSeen:   s.now,
		References: r,
	}

	i := s.stubIndex(func(_ context.Context, args []string) interface{} {
		s.Len(args, 6)
		s.Equal("HSET", args[0])
		// We test key generation elsewhere, ignore args[1]

		s.Equal("l", args[2]) // Use the `redis` tag.
		s.Equal(s.now.Format(time.RFC3339), args[3])

		s.Equal("r", args[4]) // Use the `redis` tag.

		refs := types.References{}
		err := refs.UnmarshalBinary([]byte(args[5]))
		s.NoError(err)

		s.Equal(refs, u.References)

		return nil
	})

	err := i.set(s.ctx, testId, u)
	s.NoError(err)
}

func (s *ExistsIndexTestSuite) TestString() {
	i := s.stubIndex(func(_ context.Context, _ []string) interface{} {
		return nil
	})

	str := i.String()
	s.Equal(indexName, str)
}

// Identical to set
// func (s *ExistsIndexTestSuite) TestIndex() {}
// func (s *ExistsIndexTestSuite) TestUpdate() {}

func (s *ExistsIndexTestSuite) TestDelete() {
	i := s.stubIndex(func(_ context.Context, args []string) interface{} {
		s.Len(args, 2)
		s.Equal("UNLINK", args[0])
		// We test key generation elsewhere, ignore args[1]
		return nil
	})
	err := i.Delete(s.ctx, testId)
	s.NoError(err)
}

func (s *ExistsIndexTestSuite) TestGetFound() {
	nBytes, _ := s.now.MarshalText()

	r1 := types.Reference{
		ParentHash: "p1",
		Name:       "f1",
	}
	r := types.References{
		r1,
	}
	rBytes, _ := r.MarshalBinary()

	u := &types.Update{
		LastSeen:   s.now,
		References: r,
	}

	i := s.stubIndex(func(_ context.Context, args []string) interface{} {
		s.Len(args, 2)
		s.Equal("HGETALL", args[0])

		// return []string{"l", string(nBytes)}

		return [][]byte{
			{'l'}, nBytes,
			{'r'}, rBytes,
		}
	})

	dst := &types.Update{}

	found, err := i.Get(s.ctx, testId, dst)
	s.NoError(err)
	s.True(found)

	s.Equal(u, dst)
}

func (s *ExistsIndexTestSuite) TestGetNotFound() {
	i := s.stubIndex(func(_ context.Context, args []string) interface{} {
		s.Len(args, 2)
		s.Equal("HGETALL", args[0])

		// Return empty list when key does not exist.
		return []string{}
	})

	dst := &types.Update{}

	found, err := i.Get(s.ctx, testId, dst)
	s.NoError(err)
	s.False(found)
}

func (s *ExistsIndexTestSuite) TestError() {
	testErr := errors.New("error")
	i := s.stubIndex(func(_ context.Context, args []string) interface{} {
		return testErr
	})

	dst := &types.Update{}

	found, err := i.Get(s.ctx, testId, dst)
	s.Error(err)
	s.False(found)
	s.Equal(testErr, err)
}

func TestExistsIndexTestSuite(t *testing.T) {
	suite.Run(t, new(ExistsIndexTestSuite))
}
