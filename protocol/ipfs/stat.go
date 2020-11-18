package ipfs

import (
	"context"

	t "github.com/ipfs-search/ipfs-search/types"
)

type statResult struct {
	Hash string
	Type string
	Size int64 // unixfs size
}

// Stat returns a ReferencedResource with Type and Size populated.
func (i *IPFS) Stat(ctx context.Context, r *t.Resource) (*t.ReferencedResource, error) {
	const cmd = "files/stat"

	ctx, cancel := context.WithTimeout(ctx, i.config.StatTimeout)
	defer cancel()

	path := absolutePath(r)
	req := i.shell.Request(cmd, path)

	result := new(statResult)

	if err := req.Exec(ctx, result); err != nil {
		return nil, err
	}

	return &t.ReferencedResource{
		Resource: r,
		Reference: &t.Reference{
			Type: typeFromString(result.Type),
			Size: uint64(result.Size),
		},
	}, nil
}
