package multis

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

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

func writeText(w io.Writer, result Result) error {
	for _, row := range result.Rows {
		var bands strings.Builder
		for i, band := range row.Bands {
			if i > 0 {
				bands.WriteString(", ")
			}
			bands.WriteString(strings.ToUpper(string(band)))
		}
		_, err := fmt.Fprintf(w, "%s %-3s (%2d): %s\n", row.Property, strings.ToUpper(row.Multi), len(row.Bands), bands.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func writeYAML(w io.Writer, result Result) error {
	encoder := yaml.NewEncoder(w)
	return encoder.Encode(&result)
}

func writeJSON(w io.Writer, result Result) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(&result)
}
