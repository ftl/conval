package conval

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefinition_ExchangeFields(t *testing.T) {
	tt := []struct {
		desc     string
		exchange []ExchangeDefinition
		expected []ExchangeField
	}{
		{
			desc: "one definition",
			exchange: []ExchangeDefinition{
				{
					Fields: []ExchangeField{
						{RSTProperty},
						{SerialNumberProperty},
					},
				},
			},
			expected: []ExchangeField{
				{RSTProperty},
				{SerialNumberProperty},
			},
		},
		{
			desc: "two definitions, common RST field",
			exchange: []ExchangeDefinition{
				{
					Fields: []ExchangeField{
						{RSTProperty},
						{SerialNumberProperty},
					},
				},
				{
					Fields: []ExchangeField{
						{RSTProperty},
						{NoMemberProperty, "wag_dok"},
					},
				},
			},
			expected: []ExchangeField{
				{RSTProperty},
				{SerialNumberProperty, NoMemberProperty, "wag_dok"},
			},
		},
		{
			desc: "two definitions, different field count",
			exchange: []ExchangeDefinition{
				{
					Fields: []ExchangeField{
						{RSTProperty},
						{SerialNumberProperty},
					},
				},
				{
					Fields: []ExchangeField{
						{RSTProperty},
						{NameProperty},
						{StateProvinceProperty},
					},
				},
			},
			expected: []ExchangeField{
				{RSTProperty},
				{SerialNumberProperty, NameProperty},
				{EmptyProperty, StateProvinceProperty},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			definition := Definition{
				Exchange: tc.exchange,
			}
			actual := definition.ExchangeFields()
			assert.Equal(t, tc.expected, actual)
		})
	}
}

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
duration: 48h
duration-constraints:
- duration: 24h
  operator_mode: single`,
			expected: Definition{
				Name:     "Test Contest",
				Duration: 48 * time.Hour,
				DurationConstraints: []ConstrainedDuration{
					{Constraint: Constraint{OperatorMode: SingleOperator}, Duration: 24 * time.Hour},
				},
			},
		},
		{
			desc: "constrained break",
			yaml: `name: Test Contest
breaks:
- operator_mode: single
  overlay: classic
  duration: 1h`,
			expected: Definition{
				Name: "Test Contest",
				Breaks: []ConstrainedDuration{
					{Constraint: Constraint{OperatorMode: SingleOperator, Overlay: ClassicOverlay}, Duration: 1 * time.Hour},
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
		{
			desc: "multiple band change rules",
			yaml: `name: Test Contest
band_change_rules:
- operator_mode: single
  overlay: classic
  grace_period: 10m
  multiplier_exception: false
- grace_period: 10m
  multiplier_exception: true`,
			expected: Definition{
				Name: "Test Contest",
				BandChangeRules: []BandChangeRule{
					{Constraint: Constraint{OperatorMode: SingleOperator, Overlay: ClassicOverlay}, GracePeriod: 10 * time.Minute},
					{GracePeriod: 10 * time.Minute, MultiplierException: true},
				},
			},
		},
		{
			desc: "three exchange field, one with two alternatives",
			yaml: `name: Test Contest
exchange:
- fields:
  - [rst]
  - [serial]
  - [member_number, nm]`,
			expected: Definition{
				Name: "Test Contest",
				Exchange: []ExchangeDefinition{
					{
						Fields: []ExchangeField{
							{RSTProperty},
							{SerialNumberProperty},
							{MemberNumberProperty, NoMemberProperty},
						},
					},
				},
			},
		},
		{
			desc: "scoring rules",
			yaml: `name: Test Contest
scoring:
  qsos:
  - their_continent: [other]
    bands:
    - 10m
    - 15m
    - 20m
    value: 3
  - their_continent: [other]
    bands:
    - 40m
    - 80m
    - 160m
    value: 6
  - value: 1
  qso_band_rule: once_per_band
  multis:
  - property: dxcc_entity
    my_continent: [eu]
    band_rule: once_per_band
    value: 1
  - property: dxcc_entity
    my_continent: [af, an, as, na, oc, sa]
    band_rule: once_per_band
    value: 2`,
			expected: Definition{
				Name: "Test Contest",
				Scoring: Scoring{
					QSORules: []ScoringRule{
						{TheirContinent: []Continent{OtherContinent}, Bands: []ContestBand{Band10m, Band15m, Band20m}, Value: 3},
						{TheirContinent: []Continent{OtherContinent}, Bands: []ContestBand{Band40m, Band80m, Band160m}, Value: 6},
						{Value: 1},
					},
					QSOBandRule: OncePerBand,
					MultiRules: []ScoringRule{
						{Property: DXCCEntityProperty, MyContinent: []Continent{Europa}, BandRule: OncePerBand, Value: 1},
						{Property: DXCCEntityProperty, MyContinent: []Continent{Africa, Antarctica, Asia, NorthAmerica, Oceania, SouthAmerica}, BandRule: OncePerBand, Value: 2},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			buffer := bytes.NewBufferString(tc.yaml)
			actual, err := LoadDefinitionYAML(buffer)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, *actual)
		})
	}
}

func TestSaveLoadYAMLRoundtrip(t *testing.T) {
	names, err := IncludedDefinitionNames()
	require.NoError(t, err)

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			definition, err := IncludedDefinition(name)
			require.NoError(t, err)

			buffer := bytes.NewBuffer([]byte{})
			err = SaveDefinitionYAML(buffer, definition, true)
			assert.NoError(t, err, "save")

			loadedDefinition, err := LoadDefinitionYAML(buffer)
			assert.NoError(t, err, "load")

			assert.Equal(t, *definition, *loadedDefinition)
		})
	}
}

func TestSaveYAMLWithoutExamples(t *testing.T) {
	names, err := IncludedDefinitionNames()
	require.NoError(t, err)

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			definition, err := IncludedDefinition(name)
			require.NoError(t, err)

			definitionWithoutExamples := *definition
			definitionWithoutExamples.Examples = nil

			buffer := bytes.NewBuffer([]byte{})
			err = SaveDefinitionYAML(buffer, definition, false)
			assert.NoError(t, err, "save")

			loadedDefinition, err := LoadDefinitionYAML(buffer)
			assert.NoError(t, err, "load")

			assert.Equal(t, definitionWithoutExamples, *loadedDefinition)
		})
	}
}

func TestPropertyConstraint_Matches(t *testing.T) {
	tt := []struct {
		desc       string
		myValue    string
		theirValue string
		constraint PropertyConstraint
		expected   bool
	}{
		{
			desc:       "not empty, true",
			theirValue: "123",
			constraint: PropertyConstraint{
				TheirValueNotEmpty: true,
			},
			expected: true,
		},
		{
			desc:       "not empty, false",
			theirValue: "",
			constraint: PropertyConstraint{
				TheirValueNotEmpty: true,
			},
			expected: false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			actual := tc.constraint.Matches(tc.myValue, tc.theirValue)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPropertyDefinition_Validate(t *testing.T) {
	tt := []struct {
		desc   string
		yaml   string
		values []string
		valid  []bool
	}{
		{
			desc: "simple value list",
			yaml: `name: Test Contest
properties:
- name: pty
  values: [A, B, C]`,
			values: []string{"A", "a", "D"},
			valid:  []bool{true, true, false},
		},
		{
			desc: "regular expression",
			yaml: `name: Test Contest
properties:
- name: pty
  expression: "\\d*[A-Z][A-Z0-9]+"`,
			values: []string{"B36", "b36", "70DARC", "1A"},
			valid:  []bool{true, true, true, false},
		},
		{
			desc: "regular expression over value list",
			yaml: `name: Test Contest
properties:
- name: pty
  values: [1A, B, C]
  expression: "\\d*[A-Z][A-Z0-9]+"`,
			values: []string{"B36", "b36", "70DARC", "1A", "B", "C"},
			valid:  []bool{true, true, true, false, false, false},
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			buffer := bytes.NewBufferString(tc.yaml)
			definition, err := LoadDefinitionYAML(buffer)
			assert.NoError(t, err)

			validator, ok := definition.PropertyValidator("pty")
			assert.True(t, ok)
			assert.NotNil(t, validator)

			for i, value := range tc.values {
				err := validator.ValidateProperty(value, nil)
				if tc.valid[i] {
					assert.NoError(t, err, i)
				} else {
					assert.Error(t, err, i)
				}
			}
		})
	}
}

func TestPropertyDefinition_Get(t *testing.T) {
	tt := []struct {
		desc     string
		yaml     string
		exchange []string
		expected []string
	}{
		{
			desc: "simple value list",
			yaml: `name: Test Contest
properties:
- name: pty
  values: [A, B, C]
exchange:
- fields:
  - [pty]`,
			exchange: []string{"A", "a", "D"},
			expected: []string{"A", "A", ""},
		},
		{
			desc: "regular expression",
			yaml: `name: Test Contest
properties:
- name: pty
  expression: "\\d*[A-Z][A-Z0-9]+"
exchange:
- fields:
  - [pty]`,
			exchange: []string{"B36", "b36", "70DARC", "1A"},
			expected: []string{"B36", "B36", "70DARC", ""},
		},
		{
			desc: "regular expression over value list",
			yaml: `name: Test Contest
properties:
- name: pty
  values: [1A, B, C]
  expression: "\\d*[A-Z][A-Z0-9]+"
exchange:
- fields:
  - [pty]`,
			exchange: []string{"B36", "b36", "70DARC", "1A", "B", "C"},
			expected: []string{"B36", "B36", "70DARC", "", "", ""},
		},
		{
			desc: "derived property",
			yaml: `name: Test Contest
properties:
- name: src
  expression: "\\d*[A-Z][A-Z0-9]+"
- name: pty
  source: src
  expression: "\\d*([A-Z]).*"
exchange:
- fields: 
  - [src]`,
			exchange: []string{"B36", "b36", "70DARC", "1A", "B", "C"},
			expected: []string{"B", "B", "D", "", "", ""},
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			buffer := bytes.NewBufferString(tc.yaml)
			definition, err := LoadDefinitionYAML(buffer)
			assert.NoError(t, err)

			getter, ok := definition.PropertyGetter("pty")
			assert.True(t, ok)
			assert.NotNil(t, getter)

			for i, exchange := range tc.exchange {
				parsedExchange := ParseExchange(definition.ExchangeFields(), []string{exchange}, nil, definition)
				qso := QSO{
					TheirExchange: parsedExchange,
				}

				actual := getter.GetProperty(qso)
				assert.Equal(t, tc.expected[i], actual, i)
			}
		})
	}
}
