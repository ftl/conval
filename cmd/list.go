package main

import (
	"fmt"
	"log"

	"github.com/ftl/conval"
	"github.com/spf13/cobra"
)

// var listFlags = struct {
// 	// add the flags for the list command here
// }{}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all included contest definitions",
	Run:   runList,
}

func init() {
	// add the flags to the list command: listCmd.Flags()....

	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	names, err := conval.IncludedDefinitionNames()
	if err != nil {
		log.Fatal(err)
	}
	if len(names) == 0 {
		return
	}

	for _, name := range names {
		definition, err := conval.IncludedDefinition(name)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\t%s\t%s\n", name, definition.Name, definition.OfficialRules)
	}
}
