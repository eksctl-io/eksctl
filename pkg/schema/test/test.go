package test

import "github.com/weaveworks/eksctl/pkg/schema/test/subpkg"

// Config describes some settings for _some_ things
type Config struct {
	// Num describes the number of subthings
	Num int `json:"num"`
	// An option
	// +required
	Option        DirectType         `json:"option"`
	PointerOption *PointerType       `json:"pointeroption"`
	PackageOption subpkg.PackageType `json:"packageoption"`
	AliasedInt    Alias              `json:"aliasedint"`
	Unknown       interface{}        `json:"unknown"`
	Other         map[string]string  `json:"other"`
	// Determines the version of the main thing. Valid variants are:
	// `DefaultVersion` (default): This is the right option,
	// `LegacyVersion`: Will be deprecated,
	// `TwoPointO`,
	// `"2"`
	Version string `json:"version"`
	// Tells us which kind of config.
	// Valid variants are `Kind` constants
	Kind string `json:"kind"`
	// Valid entries are `Kind` constants
	Kinds   []string `json:"kinds"`
	SumType SumType  `json:"sumType"`
}

type SumTypeA struct {
	A string `json:"a"`
}

type SumTypeB struct {
	B string `json:"b"`
}

// Schema type is one of `SumTypeA`, `SumTypeB`
type SumType struct {
	// Valid variants are:
	// `"a"`: type A
	// `"b"`: type B
	Type string `json:"type"`
}

// Values for `Kind`
const (
	// firstKind is legacy
	firstKind = "FirstKind" // nolint:deadcode,varcheck,unused
	// secondKind should be used
	secondKind = "SecondKind" // nolint:deadcode,varcheck,unused
	// thirdKind is from some other package
	thirdKind = subpkg.SpecialKind // nolint:deadcode,varcheck,unused
	// defaultKind (default) and this comment combines with secondKind
	defaultKind = secondKind // nolint:deadcode,varcheck,unused
)

// Alias is just an int
type Alias int

// DefaultVersion is some valid value
const DefaultVersion = "X"

// LegacyVersion is old
const LegacyVersion = "Y"

// TwoPointO for literal version
const TwoPointO = "2.0"

// DirectType describes a sub configuration of the Config
type DirectType struct {
	Kind string `json:"kind"`
}

// PointerType describes a sub configuration of the Config
type PointerType struct {
	Kind string `json:"kind"`
}
