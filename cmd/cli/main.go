package main

import (
	"os"
	"path/filepath"

	"github.com/pivotal-cf-experimental/new_version_resource/cmd/check"
	"github.com/pivotal-cf-experimental/new_version_resource/cmd/in"
	"github.com/pivotal-cf-experimental/new_version_resource/cmd/out"
)

func main() {
	switch filepath.Base(os.Args[0]) {
	case "check":
		check.Main()
	case "in":
		in.Main()
	case "out":
		out.Main()
	case "cli":
		for _, name := range []string{"check", "in", "out"} {
			if err := os.Symlink("cli", name); err != nil {
				panic(err)
			}
		}
	default:
		panic("Unkown command: acceptable check/in/out/cli")
	}
}
