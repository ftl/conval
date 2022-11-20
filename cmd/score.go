package main

import (
	"log"
	"os"
	"strings"

	"github.com/ftl/hamradio/dxcc"
	"github.com/spf13/cobra"

	"github.com/ftl/conval"
	"github.com/ftl/conval/cmd/score"
)

var scoreFlags = struct {
	setupFilename      string
	definitionFilename string
	outputFormat       string
}{}

var scoreCmd = &cobra.Command{
	Use:   "score <filename>",
	Short: "calculate the score of a contest log given as Cabrillo or ADIF file",
	Run:   runScore,
}

func init() {
	scoreCmd.Flags().StringVar(&scoreFlags.setupFilename, "setup", "", "the setup file")
	scoreCmd.Flags().StringVar(&scoreFlags.definitionFilename, "definition", "", "the contest definition as filename or cabrillo name")
	scoreCmd.Flags().StringVar(&scoreFlags.outputFormat, "output", "total", "select the output format (total, text, yaml, json, csv)")

	rootCmd.AddCommand(scoreCmd)
}

func runScore(cmd *cobra.Command, args []string) {
	var err error
	prefixes, err := newPrefixDatabase()
	if err != nil {
		log.Fatal(err)
	}
	definition, err := score.PrepareDefinition(scoreFlags.definitionFilename)
	if err != nil {
		log.Fatal(err)
	}
	setup, err := score.PrepareSetup(scoreFlags.setupFilename, prefixes)
	if err != nil {
		log.Fatal(err)
	}

	if len(args) < 1 {
		log.Fatal("missing input filename")
	}

	for _, filename := range args {
		logfile, err := score.ReadCabrilloLogFromFile(filename, prefixes)
		if err != nil {
			log.Fatal(err)
		}

		result, err := score.Evaluate(logfile, definition, setup)
		if err != nil {
			log.Fatal(err)
		}

		err = score.WriteOutput(os.Stdout, score.OutputFormat(strings.ToLower(scoreFlags.outputFormat)), result)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func newPrefixDatabase() (*prefixDatabase, error) {
	prefixes, _, err := dxcc.DefaultPrefixes(true)
	if err != nil {
		return nil, err
	}
	return &prefixDatabase{prefixes}, nil
}

type prefixDatabase struct {
	prefixes *dxcc.Prefixes
}

func (d prefixDatabase) Find(s string) (conval.Continent, conval.DXCCEntity, bool) {
	entities, found := d.prefixes.Find(s)
	if !found || len(entities) == 0 {
		return "", "", false
	}

	return conval.Continent(strings.ToLower(entities[0].Continent)), conval.DXCCEntity(strings.ToLower(entities[0].PrimaryPrefix)), true
}
