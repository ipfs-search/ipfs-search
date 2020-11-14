package tika

import (
	"log"
	"net"
	"net/url"
	"syscall"
)

// shouldRetry(err) returns true when requests are to be retried, false when nog.
func shouldRetry(err error) bool {
	if uerr, ok := err.(*url.Error); ok {
		if uerr.Timeout() {
			// Fail on timeouts; this situation likely indicates timeouts on the IPFS side.
			return false
		}

		if uerr.Temporary() {
			// Retry on other temp errors
			log.Printf("Temporary URL error: %v", uerr)
			return true
		}

		// Somehow, the errors below are not temp errors !?
		switch t := uerr.Err.(type) {
		case *net.OpError:
			if t.Op == "dial" {
				log.Printf("Unknown host %v", t)
				return true

			} else if t.Op == "read" {
				log.Printf("Connection refused %v", t)
				return true
			}

		case syscall.Errno:
			if t == syscall.ECONNREFUSED {
				log.Printf("Connection refused %v", t)
				return true
			}
		}
	}

	// Any other errors, usually imply proper failure scenario's.
	return false
}
