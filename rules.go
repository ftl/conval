package conval

import (
	"embed"
	"fmt"
	"path"
	"strings"
)

//go:embed rules/*.yaml
var definitions embed.FS

func IncludedDefinition(name string) (*Definition, error) {
	f, err := definitions.Open(fmt.Sprintf("rules/%s.yaml", name))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return LoadYAML(f)
}

func IncludedDefinitionNames() ([]string, error) {
	entries, err := definitions.ReadDir("rules")
	if err != nil {
		return nil, err
	}
	result := make([]string, len(entries))
	for i, entry := range entries {
		result[i] = strings.TrimSuffix(entry.Name(), path.Ext(entry.Name()))
	}
	return result, nil
}
