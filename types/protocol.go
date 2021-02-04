package types

// Protocol is an enum specifying the protocol.
type Protocol uint8

const (
	// InvalidProtocol (default) value signifies an invalid protocol.
	InvalidProtocol Protocol = iota
	// IPFSProtocol (currently only supported protocol)
	IPFSProtocol
)

func (p Protocol) String() string {
	switch p {
	case IPFSProtocol:
		return "ipfs"
	default:
		panic("Invalid value for Protocol.")
	}
}
