package ipfs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	unixfs_pb "github.com/ipfs/go-unixfs/pb"

	t "github.com/ipfs-search/ipfs-search/types"
)

// Note: copied from https://github.com/ipfs/go-ipfs-http-client/blob/6062f4dc5c9edafa6f1b8301e420b8439588f2fa/unixfs.go#L133
type lsLink struct {
	Name, Hash string
	Size       uint64
	Type       unixfs_pb.Data_DataType
	Target     string
}

type lsObject struct {
	Hash  string
	Links []lsLink
}

type lsOutput struct {
	Objects []lsObject
}

// decodeLink decodes an lsOutput and returns a link.
func decodeLink(dec *json.Decoder) (*lsLink, error) {
	var link lsOutput

	if err := dec.Decode(&link); err != nil {
		// Propagate other decoding errors
		return nil, fmt.Errorf("decoding json: %w", err)
	}

	if len(link.Objects) != 1 {
		return nil, errors.New("unexpected Objects len")
	}

	if len(link.Objects[0].Links) != 1 {
		return nil, errors.New("unexpected Links len")
	}

	return &link.Objects[0].Links[0], nil
}

// Ls returns a channel with AnnotatedResource's with Type and Size populated.
func (i *IPFS) Ls(ctx context.Context, r *t.AnnotatedResource, out chan<- t.AnnotatedResource) error {
	const cmd = "ls"

	path := absolutePath(r.Resource)

	req := i.shell.Request(cmd, path).
		Option("resolve-type", false).
		Option("size", false).
		Option("stream", true)

	// IMPORTANT - WE WANT TO AVOID STAT'ING OBJECTS WHILE WE'RE LISTING!
	// Rather, we might make this quick and dirty (avoids a lot of duplicate work, e.g.
	// not stat'ing existing items!!1)

	resp, err := req.Send(ctx)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return resp.Error
	}
	defer resp.Close()

	dec := json.NewDecoder(resp.Output)

	for {
		link, err := decodeLink(dec)
		if errors.Is(err, io.EOF) {
			// EOF; end of the list is a happy return
			return nil
		}

		// TODO: Consider using an error channel here; don't abort on individual decoding errors?
		// Alternativel: propagate an InvalidType object instead and log the error without propagating.
		// Needs real world testing. How many directories with invalid entries are there,
		// and should we care about them?
		if err != nil {
			return err
		}

		refR := t.AnnotatedResource{
			Resource: &t.Resource{
				Protocol: t.IPFSProtocol,
				ID:       link.Hash,
			},
			Reference: t.Reference{
				Parent: r.Resource,
				Name:   link.Name,
			},
			Stat: t.Stat{
				Type: typeFromPb(link.Type),
				Size: link.Size,
			},
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- refR:
		}
	}
}
