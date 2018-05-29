package main

// this is to keep dep happy, as we want authenticator cmd vendored

import (
	_ "github.com/heptio/authenticator/cmd/heptio-authenticator-aws"
)

func main() {}
