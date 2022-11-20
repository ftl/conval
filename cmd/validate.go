package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/ftl/conval"
)

// var validateFlags = struct {
// add the flags for the validate command here
// }{}

var validateCmd = &cobra.Command{
	Use:   "validate <filename>",
	Short: "validate the given contest definition file using the included examples",
	Run:   runValidate,
}

func init() {
	// add the flags to the validate command: validateCmd.Flags()....

	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) {
	var err error
	prefixes, err := conval.NewPrefixDatabase()
	if err != nil {
		log.Fatal(err)
	}

	if len(args) < 1 {
		log.Fatal("missing filename")
	}
	filename := args[0]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	definition, err := conval.LoadDefinitionYAML(file)
	if err != nil {
		log.Fatal(err)
	}

	if len(definition.Examples) == 0 {
		log.Fatalf("%s does not contain any examples to validate\n", filename)
	}

	err = conval.ValidateExamples(definition, prefixes)
	if err != nil {
		log.Fatal(err)
	}
}
