package ipfs

import (
	"context"

	t "github.com/ipfs-search/ipfs-search/types"
)

// 256KB is the default chunker block size. Therefore, unreferenced files with exactly
// this size are very likely to be chunks of files (partials) rather than full files.
const partialSize = 262144

type statResult struct {
	Type           string
	Size           uint64
	CumulativeSize uint64
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

// Stat returns a AnnotatedResource with Type and Size populated.
// Ref: http://docs.ipfs.io.ipns.localhost:8080/reference/http/api/#api-v0-files-stat
func (i *IPFS) Stat(ctx context.Context, r *t.AnnotatedResource) error {
	const cmd = "files/stat"

	path := absolutePath(r)
	req := i.shell.Request(cmd, path)

	result := new(statResult)

	if err := req.Exec(ctx, result); err != nil {
		return err
	}

	rType := typeFromString(result.Type)

	var size uint64
	switch rType {
	case t.FileType:
		size = result.Size
	case t.DirectoryType:
		// For directories, the size is always 0, hence CumulativeSize is appropriate.
		size = result.CumulativeSize
	}

	r.Stat = t.Stat{
		Type: rType,
		Size: size,
	}

	// Override type for *unreferenced* partials, based on size
	if r.Size == partialSize && r.Reference.Parent == nil {
		r.Stat.Type = t.PartialType
	}

	return nil
}
