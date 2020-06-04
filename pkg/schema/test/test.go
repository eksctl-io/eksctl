package test

import "github.com/weaveworks/eksctl/pkg/schema/test/subpkg"

// Config describes some settings for _some_ things
type Config struct {
	// Num describes the number of subthings
	Num           int                `json:"num"`
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
	Version       string  `json:"version"`
}

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
