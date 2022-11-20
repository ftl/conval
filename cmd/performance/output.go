package performance

import (
	"fmt"
	"io"

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
	// app.YamlOutput: OutputFunc(writeYAML),
	// app.JsonOutput: OutputFunc(writeJSON),
	// app.CsvOutput: OutputFunc(writeCSV),
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
	return nil
}
