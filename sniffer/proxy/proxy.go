package proxy

import (
	hook "github.com/alanshaw/ipfs-hookds"
	"github.com/alanshaw/ipfs-hookds/batch"
	"github.com/ipfs/go-datastore"
)

// New wraps a datastore in a proxy calling afterPut after every Put() operation.
func New(ds datastore.Batching, afterPut hook.AfterPutFunc) datastore.Batching {
	afterBatch := func(b datastore.Batch, err error) (datastore.Batch, error) {
		return batch.NewBatch(b, batch.WithAfterPut(batch.AfterPutFunc(afterPut))), err
	}

	return hook.NewBatching(ds, hook.WithAfterPut(afterPut), hook.WithAfterBatch(afterBatch))
}
