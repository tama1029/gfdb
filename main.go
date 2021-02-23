package main

import (
	"fmt"
	"github.com/tama1029/gfdb/cmd/gfdb"
	"os"
)

func main() {
	if err := gfdb.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(-1)
	}
}
