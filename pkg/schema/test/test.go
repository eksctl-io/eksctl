package test

import "github.com/weaveworks/eksctl/pkg/schema/test/subpkg"

// Config describes some settings for _some_ things
type Config struct {
	// This number describes the number of subthings
	Num           int                `json:"num"`
	Option        DirectType         `json:"option"`
	PointerOption *PointerType       `json:"pointeroption"`
	PackageOption subpkg.PackageType `json:"packageoption"`
	AliasedInt    Alias              `json:"aliasedint"`
	Unknown       interface{}        `json:"unknown"`
	Other         map[string]string  `json:"other"`
}

// Alias is just an int
type Alias int

// DirectType describes a sub configuration of the Config
type DirectType struct {
	Kind string `json:"kind"`
}

// PointerType describes a sub configuration of the Config
type PointerType struct {
	Kind string `json:"kind"`
}
