package main

import (
	"fmt"
	"os"
	"tctg-automation/internal/ticketbotcli"
)

func main() {
	if err := ticketbotcli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
