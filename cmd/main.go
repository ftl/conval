package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootFlags = struct {
}{}

var rootCmd = &cobra.Command{
	Use:   "conval",
	Short: "A simple tool to evaluate contest logs.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
}

func main() {
	Execute()
}
