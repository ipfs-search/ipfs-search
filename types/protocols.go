package types

const (
	// IPFSProtocol (currently only supported protocol)
	IPFSProtocol = iota
)

// Protocol is an enum specifying the protocol.
type Protocol uint8

func (p Protocol) String() string {
	switch p {
	case IPFSProtocol:
		return "ipfs"
	default:
		panic("Invalid value for Protocol.")
	}
}
