package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/ftl/conval"
	"github.com/spf13/cobra"
)

var openRulesCmd = &cobra.Command{
	Use:   "rules <cabrillo name>",
	Short: "open the official rules of the given contest in the browser",
	Run:   runOpenRules,
}

var openUploadCmd = &cobra.Command{
	Use:   "upload <cabrillo name>",
	Short: "open the upload page of the given contest in the browser",
	Run:   runOpenUpload,
}

func init() {
	rootCmd.AddCommand(openRulesCmd)
	rootCmd.AddCommand(openUploadCmd)
}

func runOpenRules(cmd *cobra.Command, args []string) {
	for _, name := range args {
		definition, err := conval.IncludedDefinition(name)
		if err != nil {
			log.Fatal(err)
		}
		url := definition.OfficialRules
		if url == "" {
			log.Printf("the definition of %s has no rules URL", name)
			continue
		}
		err = openInBrowser(url)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func runOpenUpload(cmd *cobra.Command, args []string) {
	for _, name := range args {
		definition, err := conval.IncludedDefinition(name)
		if err != nil {
			log.Fatal(err)
		}
		url := definition.UploadURL
		if url == "" {
			log.Printf("the definition of %s has no upload URL", name)
			continue
		}
		err = openInBrowser(url)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func openInBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unknown platform %s", runtime.GOOS)
	}
}
