package conval

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValdiateIncludedContestExamples(t *testing.T) {
	names, err := IncludedDefinitionNames()
	require.NoError(t, err)

	prefixes, err := NewPrefixDatabase()
	require.NoError(t, err)

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			definition, err := IncludedDefinition(name)
			require.NoError(t, err)

			if len(definition.Examples) > 0 {
				err = ValidateExamples(definition, prefixes)
				assert.NoError(t, err)
			}
		})
	}
}
