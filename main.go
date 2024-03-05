package main

import (
	"os"

	"github.com/xxxVitoxxx/imgo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
