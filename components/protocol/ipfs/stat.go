package ipfs

import (
	"context"
	"fmt"

	t "github.com/ipfs-search/ipfs-search/types"
)

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

func getSize(rType t.ResourceType, result *statResult) uint64 {
	if rType == t.FileType {
		return result.Size
	}

	return result.CumulativeSize
}

// Stat returns a AnnotatedResource with Type and Size populated.
// Ref: http://docs.ipfs.io.ipns.localhost:8080/reference/http/api/#api-v0-files-stat
func (i *IPFS) Stat(ctx context.Context, r *t.AnnotatedResource) error {
	ctx, span := i.Tracer.Start(ctx, "protocol.ipfs.Stat")
	defer span.End()

	const cmd = "files/stat"

	path := absolutePath(r)
	req := i.shell.Request(cmd, path)

	result := new(statResult)

	if err := req.Exec(ctx, result); err != nil {
		if isInvalidResourceErr(err) {
			err = fmt.Errorf("%w: %v", t.ErrInvalidResource, err)
		}

		span.RecordError(err)
		return err
	}

	rType := typeFromString(result.Type)
	size := getSize(rType, result)

	r.Stat = t.Stat{
		Type: rType,
		Size: size,
	}

	// Override type for *unreferenced* partials, based on size
	if r.Size == uint64(i.config.PartialSize) && r.Reference.Parent == nil {
		r.Stat.Type = t.PartialType
	}

	return nil
}
