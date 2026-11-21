package main

import (
	"os"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

var version = "dev"

func main() {
	os.Exit(cli.App{Version: version}.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
