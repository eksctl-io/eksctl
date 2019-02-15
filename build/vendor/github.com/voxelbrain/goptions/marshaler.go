package goptions

type Marshaler interface {
	MarshalGoption(s string) error
}
