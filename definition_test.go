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
		{
			desc: "single mode, dual mode, all mode categories",
			yaml: `name: Test Contest
categories:
- name: Single Mode
  modes: [cw]
- name: Dual Mode
  modes: [cw, rtty]
- name: All Mode
  modes: [all]`,
			expected: Definition{
				Name: "Test Contest",
				Categories: []Category{
					{Name: "Single Mode", Modes: []Mode{"cw"}},
					{Name: "Dual Mode", Modes: []Mode{"cw", "rtty"}},
					{Name: "All Mode", Modes: []Mode{"all"}},
				},
			},
		},
		{
			desc: "known and unknown overlays",
			yaml: `name: Test Contest
overlays:
- tb_wires
- something_special`,
			expected: Definition{
				Name:     "Test Contest",
				Overlays: []Overlay{"tb_wires", "something_special"},
			},
		},
		{
			desc: "modes and bands",
			yaml: `name: Test Contest
modes:
- cw
- ssb
bands:
- 80m
- 40m
- 20m`,
			expected: Definition{
				Name:  "Test Contest",
				Modes: []Mode{"cw", "ssb"},
				Bands: []ContestBand{"80m", "40m", "20m"},
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
