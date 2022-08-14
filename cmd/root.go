package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootFlags = struct {
	setupFilename string
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
	rootCmd.PersistentFlags().StringVar(&rootFlags.setupFilename, "setup", "setup.yaml", "the setup file")
}
