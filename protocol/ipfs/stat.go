package ipfs

import (
	"context"

	t "github.com/ipfs-search/ipfs-search/types"
)

const partialSize = 262144

type statResult struct {
	Hash string
	Type string
	Size int64 // unixfs size
}

func typeFromString(strType string) t.ResourceType {
	switch strType {
	case "file":
		return t.FileType
	case "directory":
		return t.DirectoryType
	default:
		return t.UnsupportedType
	}
}

func isPartial(r *t.AnnotatedResource) bool {
	return r.Size == partialSize
}

// Stat returns a AnnotatedResource with Type and Size populated.
func (i *IPFS) Stat(ctx context.Context, r *t.Resource) (*t.AnnotatedResource, error) {
	const cmd = "files/stat"

	ctx, cancel := context.WithTimeout(ctx, i.config.StatTimeout)
	defer cancel()

	path := absolutePath(r)
	req := i.shell.Request(cmd, path)

	result := new(statResult)

	if err := req.Exec(ctx, result); err != nil {
		return nil, err
	}

	annotatedResource := &t.AnnotatedResource{
		Resource: r,
		Stat: t.Stat{
			Type: typeFromString(result.Type),
			Size: uint64(result.Size),
		},
	}

	if isPartial(annotatedResource) {
		annotatedResource.Type = t.PartialType
	}

	return annotatedResource, nil
}
