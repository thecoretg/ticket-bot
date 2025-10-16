package main

import (
	"os"

	"github.com/thecoretg/ticketbot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
