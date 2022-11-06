package conval

import (
	"embed"
	"fmt"
	"path"
	"strings"
)

// use the cabrillo name as filename and identifier, see https://www.contestcalendar.com/cabnames.php

//go:embed rules/*.yaml
var definitions embed.FS

func IncludedDefinition(name string) (*Definition, error) {
	filename := fmt.Sprintf("rules/%s.yaml", strings.ToLower(strings.TrimSpace(name)))
	f, err := definitions.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return LoadDefinitionYAML(f)
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
