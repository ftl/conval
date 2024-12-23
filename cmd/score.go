package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/ftl/conval"
	"github.com/ftl/conval/app"
	"github.com/ftl/conval/cmd/score"
)

var scoreFlags = struct {
	outputFormat string
}{}

var scoreCmd = &cobra.Command{
	Use:   "score <filename>",
	Short: "calculate the score of a contest log file",
	Run:   runScore,
}

func init() {
	scoreCmd.Flags().StringVar(&scoreFlags.outputFormat, "output", "total", "select the output format (total, text, yaml, json)")

	rootCmd.AddCommand(scoreCmd)
}

func runScore(cmd *cobra.Command, args []string) {
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

		result, err := score.Evaluate(logfile, definition, setup)
		if err != nil {
			log.Fatal(err)
		}

		err = score.WriteOutput(os.Stdout, app.ParseOutputFormat(scoreFlags.outputFormat), result)
		if err != nil {
			log.Fatal(err)
		}
	}
}
