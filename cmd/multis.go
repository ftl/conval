package main

import (
	"log"
	"os"

	"github.com/ftl/conval"
	"github.com/ftl/conval/app"
	"github.com/ftl/conval/cmd/multis"
	"github.com/spf13/cobra"
)

var multisFlags = struct {
	outputFormat string
}{}

var multisCmd = &cobra.Command{
	Use:   "multis <filename>",
	Short: "calculate the multis of a contest log file",
	Run:   runMultis,
}

func init() {
	multisCmd.Flags().StringVar(&multisFlags.outputFormat, "output", "text", "select the output format (text, yaml, json)")

	rootCmd.AddCommand(multisCmd)
}

func runMultis(cmd *cobra.Command, args []string) {
	var err error
	definition, err := app.PrepareDefinition(rootFlags.definitionName)
	if err != nil {
		log.Fatal(err)
	}
	prefixes, err := conval.NewPrefixDatabase(definition.ARRLCountryList)
	if err != nil {
		log.Fatal(err)
	}
	setup, err := app.PrepareSetup(rootFlags.setupFilename, prefixes)
	if err != nil {
		log.Fatal(err)
	}

	if len(args) < 1 {
		log.Fatal("missing input filename")
	}

	for _, filename := range args {
		logfile, err := app.ReadCabrilloLogFromFile(filename, prefixes)
		if err != nil {
			log.Fatal(err)
		}

		result, err := multis.Evaluate(logfile, definition, setup)
		if err != nil {
			log.Fatal(err)
		}

		err = multis.WriteOutput(os.Stdout, app.ParseOutputFormat(multisFlags.outputFormat), result)
		if err != nil {
			log.Fatal(err)
		}
	}
}
