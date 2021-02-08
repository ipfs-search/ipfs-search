package ipfs

import (
	"context"
	"errors"
	ipfs "github.com/ipfs/go-ipfs-api"
	"log"
)

// isInvalidResourceErr determines whether an error returned by protocol methods represents invalid content.
func isInvalidResourceErr(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		// Timeouts are explicitly not protocol errors
		return false
	}

	ipfsErr, ok := err.(*ipfs.Error)

	if !ok {
		log.Printf("Unexpected protocol error: %T:%v", err, err)
		return false
	}

	log.Printf("*ipfs.Error: %v", ipfsErr.Message)

	switch ipfsErr.Message {
	case "proto: required field \"Type\" not set", // Example: QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8
		"proto: unixfs_pb.Data: illegal tag 0 (wire type 0)", // Example: QmQkaTUmqcdGAXKaFXpe8t8yaEDGHe7xGQJHcfihrzAFTj
		"proto: unixfs_pb.Data: illegal tag 0 (wire type 2)",
		"unexpected EOF",                 // Example: QmdtMPULYK2xBVt2stYdAdxmuQukbJNFEgsdB5KV3jvsBq
		"unrecognized object type: 144",  // Example: z43AaGEvwdfzjrCZ3Sq7DKxdDHrwoaPQDtqF4jfdkNEVTiqGVFW
		"not unixfs node (proto or raw)": // Example: z8mWaJHXieAVxxLagBpdaNWFEBKVWmMiE
		return true
	}

	// *shell.Error
	// files/stat: proto: can't skip unknown wire type 6
	// files/stat: proto: can't skip unknown wire type 6
	// files/stat: proto: can't skip unknown wire type 6

	return false
}
