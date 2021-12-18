package types

// SourceType is an enum of possible sources for resources.
type SourceType uint8

const (
	// UnknownSource represents items for which no source is defined.
	UnknownSource SourceType = iota
	// SnifferSource represents items sourced from the sniffer.
	SnifferSource
	// DirectorySource represents items sourced from a directory.
	DirectorySource
)

func (t SourceType) String() string {
	switch t {
	case UnknownSource:
		return "unknown"
	case SnifferSource:
		return "sniffer"
	case DirectorySource:
		return "directory"
	default:
		panic("Invalid value for SourceType.")
	}
}
