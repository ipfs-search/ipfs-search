package crawler

import (
	"fmt"
)

// hashURL returns the IPFS URL for a particular hash
func hashURL(hash string) string {
	return fmt.Sprintf("/ipfs/%s", hash)
}
