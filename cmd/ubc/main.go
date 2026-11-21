package main

import (
	"os"
	"time"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

var version = "dev"

func main() {
	app := cli.App{Version: version, RandomSeed: clockSeed}
	os.Exit(app.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func clockSeed() uint64 {
	return uint64(time.Now().UnixNano())
}
