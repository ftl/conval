package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootFlags = struct {
	setupFilename  string
	definitionName string
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
	rootCmd.PersistentFlags().StringVar(&rootFlags.setupFilename, "setup", "", "the setup file")
	rootCmd.PersistentFlags().StringVar(&rootFlags.definitionName, "definition", "", "the contest definition as filename or cabrillo name")
}

func main() {
	Execute()
}
