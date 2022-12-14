package score

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	"github.com/ftl/conval/app"
)

type Output interface {
	Write(io.Writer, Result) error
}

type OutputFunc func(io.Writer, Result) error

func (f OutputFunc) Write(w io.Writer, result Result) error {
	return f(w, result)
}

var outputFormats = map[app.OutputFormat]Output{
	"total":        OutputFunc(writeTotal),
	app.TextOutput: OutputFunc(writeText),
	app.YamlOutput: OutputFunc(writeYAML),
	app.JsonOutput: OutputFunc(writeJSON),
}

func OutputFormats() []app.OutputFormat {
	result := make([]app.OutputFormat, 0, len(outputFormats))
	for format := range outputFormats {
		result = append(result, format)
	}
	return result
}

func WriteOutput(w io.Writer, format app.OutputFormat, result Result) error {
	output, ok := outputFormats[format]
	if !ok {
		return fmt.Errorf("unknown output format %s", format)
	}

	return output.Write(w, result)
}

func writeTotal(w io.Writer, result Result) error {
	_, err := fmt.Fprintf(w, "%d\n", result.Total)
	return err
}

func writeText(w io.Writer, result Result) error {
	_, err := fmt.Fprintf(w, "QSOs   : % 8d\nMultis : % 8d\nPoints : % 8d\nTotal  : % 8d\n", result.QSOs, result.Multis, result.Points, result.Total)
	return err
}

func writeYAML(w io.Writer, result Result) error {
	encoder := yaml.NewEncoder(w)
	return encoder.Encode(&result)
}

func writeJSON(w io.Writer, result Result) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(&result)
}
