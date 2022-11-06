package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ftl/hamradio/dxcc"
	"github.com/spf13/cobra"

	"github.com/ftl/conval"
)

var scoreFlags = struct {
	setupFilename      string
	definitionFilename string
	cabrilloName       string
}{}

var scoreCmd = &cobra.Command{
	Use:   "score <filename>",
	Short: "calculate the score of a contest log given as Cabrillo or ADIF file",
	Run:   runScore,
}

func init() {
	scoreCmd.Flags().StringVar(&scoreFlags.setupFilename, "setup", "setup.yaml", "the setup file")
	scoreCmd.Flags().StringVar(&scoreFlags.definitionFilename, "definition", "", "the contest definition file")
	scoreCmd.Flags().StringVar(&scoreFlags.cabrilloName, "cabrillo", "", "the cabrillo name (see https://www.contestcalendar.com/cabnames.php)")

	rootCmd.AddCommand(scoreCmd)
}

func runScore(cmd *cobra.Command, args []string) {
	var err error

	prefixes, err := newPrefixDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the CLI information

	var definition *conval.Definition
	if scoreFlags.definitionFilename != "" {
		definition, err = loadDefinitionFromFile(scoreFlags.definitionFilename)
	} else if scoreFlags.cabrilloName != "" {
		definition, err = conval.IncludedDefinition(scoreFlags.cabrilloName)
	}
	if err != nil {
		log.Fatal(err)
	}
	if definition != nil {
		log.Printf("%s: %s", definition.Identifier, definition.Name)
	}

	var setup *conval.Setup
	if scoreFlags.setupFilename != "" {
		setup, err = loadSetupFromFile(scoreFlags.setupFilename, prefixes)
		if err != nil {
			log.Fatal(err)
		}
	}
	// TODO override the setup information with values given with explicit CLI flags (TBD)
	if setup != nil {
		log.Printf("setup: %+v", setup)
	}

	if len(args) < 1 {
		log.Fatal("missing input filename")
	}
	for _, filename := range args {
		log.Printf("evaluating %s", filename)

		_ = definition
		_ = setup

		// # Load the logfile

		// detect the input file format
		// read the input file and get a Logfile object

		// get the contest definition using the following order
		// - the definition given by the  CLI flags (definition, cabrillo)
		// - the cabrillo name from the log file

		// get setup information from the input file
		// add the setup information from the setup file
		// add the setup information from the CLI flags

		// # Evaluate the log file

		// create a conval.Counter
		// add the QSOs to the counter
		// calculate the score
		// print the calculated score to stdout
	}
}

func loadDefinitionFromFile(filename string) (*conval.Definition, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return conval.LoadDefinitionYAML(file)
}

func loadSetupFromFile(filename string, prefixes prefixDatabase) (*conval.Setup, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	result, err := conval.LoadSetupYAML(file)
	if err != nil {
		return nil, err
	}

	matchingEntities, found := prefixes.Find(result.MyCall.String())
	if found {
		dxccEntity := matchingEntities[0]
		result.MyContinent = conval.Continent(dxccEntity.Continent)
		result.MyCountry = conval.DXCCEntity(dxccEntity.PrimaryPrefix)
	}

	return result, nil
}

type prefixDatabase interface {
	Find(s string) ([]dxcc.Prefix, bool)
}

func newPrefixDatabase() (prefixDatabase, error) {
	localFilename, err := dxcc.LocalFilename()
	if err != nil {
		return nil, err
	}
	updated, err := dxcc.Update(dxcc.DefaultURL, localFilename)
	if err != nil {
		fmt.Printf("update of local copy failed: %v\n", err)
	}
	if updated {
		fmt.Printf("updated local copy: %v\n", localFilename)
	}

	return dxcc.LoadLocal(localFilename)
}
