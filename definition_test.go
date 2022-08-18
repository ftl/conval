package conval

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadYAML(t *testing.T) {
	tt := []struct {
		desc     string
		yaml     string
		expected Definition
	}{
		{
			desc: "valid name, identifier, official rules",
			yaml: `name: Test Contest
identifier: TEST-CONTEST-VALID
official_rules: https://github.com/ftl/conval/testdata`,
			expected: Definition{
				Name:          "Test Contest",
				Identifier:    "TEST-CONTEST-VALID",
				OfficialRules: "https://github.com/ftl/conval/testdata",
			},
		},
		{
			desc: "constrained and unconstrained duration",
			yaml: `name: Test Contest
durations:
- duration: 24h
  operator: single
- duration: 48h`,
			expected: Definition{
				Name: "Test Contest",
				Durations: []ConstrainedDuration{
					{Constraint: Constraint{Operator: SingleOperator}, Duration: 24 * time.Hour},
					{Duration: 48 * time.Hour},
				},
			},
		},
		{
			desc: "constrained break",
			yaml: `name: Test Contest
breaks:
- operator: single
  overlay: classic
  duration: 1h`,
			expected: Definition{
				Name: "Test Contest",
				Breaks: []ConstrainedDuration{
					{Constraint: Constraint{Operator: SingleOperator, Overlay: ClassicOverlay}, Duration: 1 * time.Hour},
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			buffer := bytes.NewBufferString(tc.yaml)
			actual, err := LoadYAML(buffer)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, *actual)
		})
	}
}
