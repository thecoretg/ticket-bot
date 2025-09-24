package main

import (
	"fmt"
	"os"

	"github.com/thecoretg/ticketbot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("An error occured: %v\n", err)
		os.Exit(1)
	}
}
