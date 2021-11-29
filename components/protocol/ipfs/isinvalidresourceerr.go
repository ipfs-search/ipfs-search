package ipfs

import (
	"context"
	"errors"
	ipfs "github.com/ipfs/go-ipfs-api"
	"log"
	"strings"
)

var invalidErrorPrefixes = [...]string{
	"proto: required field",                // Example: QmYAqhbqNDpU7X9VW6FV5imtngQ3oBRY35zuDXduuZnyA8
	"proto: unixfs_pb.Data: illegal tag 0", // Example: QmQkaTUmqcdGAXKaFXpe8t8yaEDGHe7xGQJHcfihrzAFTj
	"unexpected EOF",                       // Example: QmdtMPULYK2xBVt2stYdAdxmuQukbJNFEgsdB5KV3jvsBq
	"unrecognized object type",             // Example: z43AaGEvwdfzjrCZ3Sq7DKxdDHrwoaPQDtqF4jfdkNEVTiqGVFW
	"not unixfs node",                      // Example: z8mWaJHXieAVxxLagBpdaNWFEBKVWmMiE
	"proto: can't skip unknown wire type",
	"failed to decode Protocol Buffers",
	"protobuf: (PBNode) invalid wireType",
}

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

	for _, p := range invalidErrorPrefixes {
		if strings.HasPrefix(ipfsErr.Message, p) {
			// Known invalid.
			return true
		}
	}

	log.Printf("Unexpected *ipfs.Error: %v", ipfsErr.Message)

	return false
}
