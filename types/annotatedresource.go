package types

type AnnotatedResource struct {
	*Resource
	Reference
	Stat
}
