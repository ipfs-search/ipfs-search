package types

const (
	// Invalid (default) value signifies an invalid protocol.
	Invalid = iota
	// IPFSProtocol (currently only supported protocol)
	IPFSProtocol
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
