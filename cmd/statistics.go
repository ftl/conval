package main

import (
	"log"
	"os"
	"time"

	"github.com/ftl/conval/app"
	"github.com/ftl/conval/cmd/statistics"
	"github.com/spf13/cobra"
)

var statisticsFlags = struct {
	startTime    string
	outputFormat string
}{}

var statisticsCmd = &cobra.Command{
	Use:   "statistics <filename>",
	Short: "calculate statistics of a contest log file",
	Run:   runStatistics,
}

func init() {
	statisticsCmd.Flags().StringVar(&statisticsFlags.startTime, "start", "", "the offical start time of the contest (2006-01-02T15:04:05Z)")
	statisticsCmd.Flags().StringVar(&statisticsFlags.outputFormat, "output", "text", "select the output format (text, yaml, json)")

	rootCmd.AddCommand(statisticsCmd)
}

func runStatistics(cmd *cobra.Command, args []string) {
	var err error
	prefixes, err := app.NewPrefixDatabase()
	if err != nil {
		log.Fatal(err)
	}
	definition, err := app.PrepareDefinition(rootFlags.definitionName)
	if err != nil {
		log.Fatal(err)
	}
	setup, err := app.PrepareSetup(rootFlags.setupFilename, prefixes)
	if err != nil {
		log.Fatal(err)
	}
	startTime, err := time.Parse(time.RFC3339, statisticsFlags.startTime)
	if err != nil {
		log.Fatalf("the start time is invalid: %v", err)
	}

	if len(args) < 1 {
		log.Fatal("missing input filename")
	}

	for _, filename := range args {
		logfile, err := app.ReadCabrilloLogFromFile(filename, prefixes)
		if err != nil {
			log.Fatal(err)
		}

		result, err := statistics.Evaluate(logfile, definition, setup, startTime)
		if err != nil {
			log.Fatal(err)
		}

		err = statistics.WriteOutput(os.Stdout, app.ParseOutputFormat(statisticsFlags.outputFormat), result)
		if err != nil {
			log.Fatal(err)
		}
	}
}
