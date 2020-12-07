package ipfs

import (
	ipfs "github.com/ipfs/go-ipfs-api"
)

// IsInvalidResourceErr determines whether an error returned by protocol methods
// represents invalid content.
func (i *IPFS) IsInvalidResourceErr(err error) bool {
	ipfsErr, ok := err.(*ipfs.Error)

	if !ok {
		return false
	}

	switch ipfsErr.Message {
	case "proto: required field \"Type\" not set", // Example: QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8
		"proto: unixfs_pb.Data: illegal tag 0 (wire type 0)", // Example: QmQkaTUmqcdGAXKaFXpe8t8yaEDGHe7xGQJHcfihrzAFTj
		"unexpected EOF",                 // Example: QmdtMPULYK2xBVt2stYdAdxmuQukbJNFEgsdB5KV3jvsBq
		"unrecognized object type: 144",  // Example: z43AaGEvwdfzjrCZ3Sq7DKxdDHrwoaPQDtqF4jfdkNEVTiqGVFW
		"not unixfs node (proto or raw)": // Example: z8mWaJHXieAVxxLagBpdaNWFEBKVWmMiE
		return true
	}

	return false
}
