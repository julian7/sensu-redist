package main

import (
	"fmt"
	"os"
)

var version = "SNAPSHOT"

func main() {
	cmd, err := rootCmd()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
