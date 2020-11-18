package types

// ResourceType is an enum of possible resource types.
type ResourceType uint8

const (
	// UndefinedType is for resources for which the type is as of yet undefined.
	UndefinedType = iota
	// UnsupportedType represents an as-of-yet unsupported type.
	UnsupportedType
	// FileType is a regular file.
	FileType
	// DirectoryType is a directory.
	DirectoryType
)

func (t ResourceType) String() string {
	switch t {
	case UndefinedType:
		return "undefined"
	case FileType:
		return "file"
	case DirectoryType:
		return "directory"
	default:
		panic("Invalid value for ResourceType.")
	}
}
