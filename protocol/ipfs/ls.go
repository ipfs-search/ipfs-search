package ipfs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	unixfs "github.com/ipfs/go-unixfs"
	unixfs_pb "github.com/ipfs/go-unixfs/pb"

	t "github.com/ipfs-search/ipfs-search/types"
)

var (
	ErrUnexpectedObjectsLen = errors.New("unexpected Objects len")
	ErrUnexpectedLinksLen   = errors.New("unexpected Links len")
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

func typeFromPb(pbType unixfs_pb.Data_DataType) t.ResourceType {
	// Note: even though both resolve type and size are set to false, it seems that object
	// types are resolved. This might be a bug in the underlying implementation.
	// Hence we should not expect returned objects to have a type defined. When they are
	// not, they default to the unixfs_pb zero type of Raw.
	//
	// Performance-wise, not resolving here is strongly preferable (otherwise, referred
	// blocks need to be fetched).
	//
	// Current trace analysis (a real price-winning implementation!):
	//
	// 1. HTTP API returns numeric type based on unixfs_pb
	//    Zero-value of Type (0) is Raw.
	//
	//    https://github.com/ipfs/go-unixfs/blob/0faf57387de7e336a68a7ed5a9c35308cb98f576/pb/unixfs.proto
	//    http://docs.ipfs.io.ipns.localhost:8080/reference/http/api/#api-v0-ls
	//
	// 2. go-ipfs core API
	//    Maps DirEntry.Type from interface-go-ipfs-core interface to unixfs_pb
	//    iface.TFile -> unixfs.TFile
	//    iface.TDirectory -> unixfs.TDirectory
	//    iface.TUnknown -> not mapped (unixfs_pb 0 value of Raw)
	//
	//    https://github.com/ipfs/interface-go-ipfs-core/blob/master/unixfs.go#L50
	//    https://github.com/ipfs/go-ipfs/blob/5ec98e14016950510d8004c7acf306876c7ef4c0/core/commands/ls.go#L146
	//
	// 3. go-ipfs unixfs core API
	//    Implements interface-go-ipfs-core. Maps unixfs_pb to DirEntry.Type (!):
	//
	//    unixfs.TFile, unixfs.TRaw -> iface.TFile
	//    unixfs.THAMTShard, unixfs.TDirectory, unixfs.TMetadata -> iface.TDirectory
	//
	//    But only when ResolveChildren is true. If not, for DagProtobuf the lnk.Type is not set, causing
	//    it to default type defined in the core interface, being iface.TUnknown.
	//
	//    For Raw leave nodes, the type is set to iface.TFile.
	//
	//	  https://github.com/ipfs/go-ipfs/blob/5ec98e14016950510d8004c7acf306876c7ef4c0/core/commands/ls.go#L135
	//
	// Hence, if `resolve-type` and `size` are both `false`, `ResolveChildren` *should* be `false` as well and
	// the UnixFS implementation of go-ipfs should have `Type = TUnknown`. The core API should map this
	// to unixfs_pb type Raw, which causes the HTTP API to return 0.
	//
	// However, if `ResolveChildren` is *not* `false`, as seems to be the case, unixfs.TRaw is
	// mapped to iface.TFile and then back to unixfs.TFile.
	//
	// Hence, we probably want to map unixfs `Raw` to `UndefinedType`.
	//
	// Note that it *seems* that HAMT sharded directories *include* type information in the directory and
	// hence do not rely on protobuf types. Hence, type and size information will be included at no
	// additional costs, while normal directories will always have type and size set to their 0-value.

	switch pbType {
	case unixfs.TRaw:
		// This could both be a file as well as an unresolved type.
		return t.UndefinedType
	case unixfs.TFile:
		return t.FileType
	case unixfs.THAMTShard, unixfs.TDirectory, unixfs.TMetadata:
		return t.DirectoryType
	default:
		return t.UnsupportedType
	}
}

// decodeLink decodes an lsOutput and returns a link.
func decodeLink(dec *json.Decoder) (*lsLink, error) {
	var link lsOutput

	if err := dec.Decode(&link); err != nil {
		// Propagate other decoding errors
		return nil, fmt.Errorf("decoding json: %w", err)
	}

	if len(link.Objects) != 1 {
		return nil, ErrUnexpectedObjectsLen
	}

	if len(link.Objects[0].Links) != 1 {
		return nil, ErrUnexpectedLinksLen
	}

	return &link.Objects[0].Links[0], nil
}

// Ls returns a channel with AnnotatedResource's with Type and Size populated.
func (i *IPFS) Ls(ctx context.Context, r *t.AnnotatedResource, out chan<- *t.AnnotatedResource) error {
	const cmd = "ls"

	path := absolutePath(r)

	resp, err := i.shell.Request(cmd, path).
		Option("resolve-type", false).
		Option("size", false).
		Option("stream", true).
		Send(ctx)
	if err != nil {
		return err
	}
	defer resp.Close()
	if resp.Error != nil {
		return resp.Error
	}

	dec := json.NewDecoder(resp.Output)

	for {
		link, err := decodeLink(dec)
		if errors.Is(err, io.EOF) {
			// EOF; end of the list is a happy return

			// Question: should we close the channel on return?
			// Probably not: channel created and 'owner' by calling context.
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
		case out <- &refR:
		}
	}
}
