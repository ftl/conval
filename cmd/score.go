package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ftl/hamradio/dxcc"
	"github.com/spf13/cobra"

	"github.com/ftl/conval"
	"github.com/ftl/conval/cmd/score"
)

var scoreFlags = struct {
	setupFilename      string
	definitionFilename string
	cabrilloName       string
	verbose            bool
}{}

var scoreCmd = &cobra.Command{
	Use:   "score <filename>",
	Short: "calculate the score of a contest log given as Cabrillo or ADIF file",
	Run:   runScore,
}

func init() {
	scoreCmd.Flags().StringVar(&scoreFlags.setupFilename, "setup", "", "the setup file")
	scoreCmd.Flags().StringVar(&scoreFlags.definitionFilename, "definition", "", "the contest definition file")
	scoreCmd.Flags().StringVar(&scoreFlags.cabrilloName, "cabrillo", "", "the cabrillo name (see https://www.contestcalendar.com/cabnames.php)")
	scoreCmd.Flags().BoolVar(&scoreFlags.verbose, "verbose", false, "enable verbose output")

	rootCmd.AddCommand(scoreCmd)
}

func runScore(cmd *cobra.Command, args []string) {
	var err error
	prefixes, err := newPrefixDatabase()
	if err != nil {
		log.Fatal(err)
	}
	definition, err := score.PrepareDefinition(scoreFlags.definitionFilename, scoreFlags.cabrilloName)
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
		// TODO: detect the input file format
		logfile, err := score.ReadCabrilloLogFromFile(filename, prefixes)
		if err != nil {
			log.Fatal(err)
		}

		result, err := score.Evaluate(logfile, definition, setup)
		if err != nil {
			log.Fatal(err)
		}

		// print the multis board
		if scoreFlags.verbose {
			for _, row := range result.MultisBoard {
				var bands strings.Builder
				for i, band := range row.Bands {
					if i > 0 {
						bands.WriteString(", ")
					}
					bands.WriteString(strings.ToUpper(string(band)))
				}
				fmt.Printf("%s %-3s (%2d): %s\n", row.Property, strings.ToUpper(row.Multi), len(row.Bands), bands.String())
			}
		}

		// print the total score
		if scoreFlags.verbose {
			fmt.Printf("QSOs   : % 8d\nMultis : % 8d\nPoints : % 8d\nTotal  : % 8d\n", result.QSOs, result.Multis, result.Points, result.Total)
		} else {
			fmt.Printf("%d\n", result.Multis*result.Points)
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
