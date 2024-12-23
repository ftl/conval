package main

import (
	"log"
	"os"
	"time"

	"github.com/ftl/conval"
	"github.com/ftl/conval/app"
	"github.com/ftl/conval/cmd/performance"
	"github.com/spf13/cobra"
)

var performanceFlags = struct {
	startTime    string
	resolution   time.Duration
	outputFormat string
}{}

var performanceCmd = &cobra.Command{
	Use:   "performance <filename>",
	Short: "evaluate the performance over time of a contest log file",
	Run:   runPerformance,
}

func init() {
	performanceCmd.Flags().StringVar(&performanceFlags.startTime, "start", "", "the offical start time of the contest (2006-01-02T15:04:05Z)")
	performanceCmd.Flags().DurationVar(&performanceFlags.resolution, "resolution", 1*time.Hour, "the time resolution")
	performanceCmd.Flags().StringVar(&performanceFlags.outputFormat, "output", "text", "select the output format (text, yaml, json, csv)")

	rootCmd.AddCommand(performanceCmd)
}

func runPerformance(cmd *cobra.Command, args []string) {
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
	startTime, err := time.Parse(time.RFC3339, performanceFlags.startTime)
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

		result, err := performance.Evaluate(logfile, definition, setup, startTime, performanceFlags.resolution)
		if err != nil {
			log.Fatal(err)
		}

		err = performance.WriteOutput(os.Stdout, app.ParseOutputFormat(performanceFlags.outputFormat), result)
		if err != nil {
			log.Fatal(err)
		}
	}
}
