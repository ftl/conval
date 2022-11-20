package performance

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
	app.CsvOutput:  OutputFunc(writeCSV),
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
	for _, row := range result.DataPoints {
		bands := strings.Join(row.Bands, ", ")
		multiValues := strings.Join(row.MultiValues, ", ")
		_, err := fmt.Fprintf(w, "%8s: % 4d|% 4d|% 4d|% 6d|%s|%s\n", row.Offset, row.QSOs, row.Points, row.Multis, row.Total, bands, multiValues)
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

func writeCSV(w io.Writer, result Result) error {
	for _, row := range result.DataPoints {
		bands := strings.Join(row.Bands, ", ")
		multiValues := strings.Join(row.MultiValues, ", ")
		_, err := fmt.Fprintf(w, "%q;%d;%d;%d;%d;%q;%q\n", row.Offset, row.QSOs, row.Points, row.Multis, row.Total, bands, multiValues)
		if err != nil {
			return err
		}
	}
	return nil
}
