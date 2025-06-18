package main

import "notashelf.dev/tct/cmd"

var version = "dev" // will be set by build process

func main() {
	cmd.Version = version
	cmd.Execute()
}
