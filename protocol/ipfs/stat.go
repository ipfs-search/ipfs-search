package ipfs

import (
	"context"

	t "github.com/ipfs-search/ipfs-search/types"
)

// 256KB is the default chunker block size. Therefore, unreferenced files with exactly
// this size are very likely to be chunks of files (partials) rather than full files.
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

// Stat returns a AnnotatedResource with Type and Size populated.
// Ref: http://docs.ipfs.io.ipns.localhost:8080/reference/http/api/#api-v0-files-stat
// 500 errors might be indicative of invalids
func (i *IPFS) Stat(ctx context.Context, r *t.AnnotatedResource) error {
	const cmd = "files/stat"

	ctx, cancel := context.WithTimeout(ctx, i.config.StatTimeout)
	defer cancel()

	path := absolutePath(r)
	req := i.shell.Request(cmd, path)

	result := new(statResult)

	if err := req.Exec(ctx, result); err != nil {
		return err
	}

	r.Stat = t.Stat{
		Type: typeFromString(result.Type),
		Size: uint64(result.Size),
	}

	// Override type for *unreferenced* partials, based on size
	if r.Size == partialSize && r.Reference.Parent == nil {
		r.Stat.Type = t.PartialType
	}

	return nil
}
