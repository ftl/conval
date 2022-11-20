package statistics

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
	_, err := fmt.Fprintf(w, `Contest: %s (%s)

Start at:     %s
Duration:     %s
Active Hours: %d

Rate   | Total  | Active | Min    | Max    | Best Hour | Slowest Hour |
-----------------------------------------------------------------------
`,
		result.ContestName,
		result.ContestID,
		result.StartTime,
		result.Duration,
		result.ActiveHours)
	if err != nil {
		return err
	}

	err = writeRate(w, "QSOs", result.QSORate)
	if err != nil {
		return err
	}
	err = writeRate(w, "Points", result.PointsRate)
	if err != nil {
		return err
	}
	err = writeRate(w, "Multis", result.MultiRate)
	if err != nil {
		return err
	}

	return err
}

func writeRate(w io.Writer, name string, rate Rate) error {
	_, err := fmt.Fprintf(w, "%-6s | % 6.1f | % 6.1f | % 6.1f | % 6.1f | %9s | %12s |\n",
		name, rate.TotalAverage, rate.ActiveAverage, rate.Min, rate.Max, rate.MaxHour, rate.MinHour)
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
