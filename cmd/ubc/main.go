package main

import (
	"os"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
