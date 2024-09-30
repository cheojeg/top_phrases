package main

import (
	"fmt"
	"github.com/cheojeg/top_phrases/cmd/cli"
)

func main() {
	rootCmd := cli.NewCmdRoot()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("error starting sign_guard")
	}
}
