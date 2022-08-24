package conval

import (
	"testing"

	"github.com/ftl/hamradio/callsign"
	"github.com/stretchr/testify/assert"
)

func TestCounter_SimplestHappyPath(t *testing.T) {
	setup := Setup{}
	rules := Scoring{
		QSORules:    []ScoringRule{{Value: 1}},
		QSOBandRule: OncePerBand,
		MultiRules:  []ScoringRule{{Value: 1}, {Value: 2}},
	}
	qso := QSO{}

	counter := NewCounter(setup, rules)

	qsoScore := counter.Add(qso)

	assert.Equal(t, 1, qsoScore.Points, "points")
	assert.Equal(t, 3, qsoScore.Multis, "multis")
}

func TestFilterScoringRules(t *testing.T) {
	tt := []struct {
		desc           string
		rules          []ScoringRule
		myContinent    Continent
		myCountry      DXCCEntity
		theirContinent Continent
		theirCountry   DXCCEntity
		band           ContestBand
		exchange       QSOExchange
		expected       []ScoringRule
	}{
		{
			desc:     "one simple unspecific rule",
			rules:    []ScoringRule{{Value: 1}},
			expected: []ScoringRule{{Value: 1}},
		},
		{
			desc: "my continent matching the specified continent",
			rules: []ScoringRule{
				{MyContinent: []Continent{Europa}, Value: 2},
				{Value: 1},
			},
			myContinent: Europa,
			expected: []ScoringRule{
				{MyContinent: []Continent{Europa}, Value: 2},
			},
		},
		{
			desc: "my country matching the specified country",
			rules: []ScoringRule{
				{MyCountry: []DXCCEntity{"dl"}, Value: 2},
				{Value: 1},
			},
			myContinent: Europa,
			myCountry:   "dl",
			expected: []ScoringRule{
				{MyCountry: []DXCCEntity{"dl"}, Value: 2},
			},
		},
		{
			desc: "their continent matching the specified continent",
			rules: []ScoringRule{
				{TheirContinent: []Continent{Africa, Antarctica, Asia, NorthAmerica, Oceania, SouthAmerica}, Value: 2},
				{Value: 1},
			},
			myContinent:    Europa,
			myCountry:      "dl",
			theirContinent: Africa,
			theirCountry:   "zs",
			expected: []ScoringRule{
				{TheirContinent: []Continent{Africa, Antarctica, Asia, NorthAmerica, Oceania, SouthAmerica}, Value: 2},
			},
		},
		{
			desc: "their country matching the specified country",
			rules: []ScoringRule{
				{TheirCountry: []DXCCEntity{"f"}, Value: 2},
				{Value: 1},
			},
			myContinent:    Europa,
			myCountry:      "dl",
			theirContinent: Europa,
			theirCountry:   "f",
			expected: []ScoringRule{
				{TheirCountry: []DXCCEntity{"f"}, Value: 2},
			},
		},
		{
			desc: "their continent is the same as my continent",
			rules: []ScoringRule{
				{TheirContinent: []Continent{SameContinent}, Value: 2},
				{Value: 1},
			},
			myContinent:    Europa,
			myCountry:      "dl",
			theirContinent: Europa,
			theirCountry:   "f",
			expected: []ScoringRule{
				{TheirContinent: []Continent{SameContinent}, Value: 2},
			},
		},
		{
			desc: "their continent is other than my continent",
			rules: []ScoringRule{
				{TheirContinent: []Continent{OtherContinent}, Value: 2},
				{Value: 1},
			},
			myContinent:    Europa,
			myCountry:      "dl",
			theirContinent: Africa,
			theirCountry:   "zs",
			expected: []ScoringRule{
				{TheirContinent: []Continent{OtherContinent}, Value: 2},
			},
		},
		{
			desc: "their country is the same as my country",
			rules: []ScoringRule{
				{TheirCountry: []DXCCEntity{SameCountry}, Value: 2},
				{Value: 1},
			},
			myContinent:    Europa,
			myCountry:      "dl",
			theirContinent: Europa,
			theirCountry:   "dl",
			expected: []ScoringRule{
				{TheirCountry: []DXCCEntity{SameCountry}, Value: 2},
			},
		},
		{
			desc: "their country is other than my country",
			rules: []ScoringRule{
				{TheirCountry: []DXCCEntity{OtherCountry}, Value: 2},
				{Value: 1},
			},
			myContinent:    Europa,
			myCountry:      "dl",
			theirContinent: Africa,
			theirCountry:   "zs",
			expected: []ScoringRule{
				{TheirCountry: []DXCCEntity{OtherCountry}, Value: 2},
			},
		},
		{
			desc: "band matching the specified band",
			rules: []ScoringRule{
				{TheirCountry: []DXCCEntity{"dl"}, Value: 2},
				{Bands: []ContestBand{Band10m, Band15m, Band20m}, Value: 3},
				{Bands: []ContestBand{Band40m, Band80m, Band160m}, Value: 4},
				{Value: 1},
			},
			myContinent:    Europa,
			myCountry:      "dl",
			theirContinent: Africa,
			theirCountry:   "zs",
			band:           Band10m,
			expected: []ScoringRule{
				{Bands: []ContestBand{Band10m, Band15m, Band20m}, Value: 3},
			},
		},
		{
			desc: "band not matching the specified band",
			rules: []ScoringRule{
				{TheirCountry: []DXCCEntity{"dl"}, Value: 2},
				{Bands: []ContestBand{Band10m, Band15m, Band20m}, Value: 3},
				{Value: 1},
			},
			myContinent:    Europa,
			myCountry:      "dl",
			theirContinent: Africa,
			theirCountry:   "zs",
			band:           Band40m,
			expected: []ScoringRule{
				{Value: 1},
			},
		},
		{
			desc: "exchange contains a specified property",
			rules: []ScoringRule{
				{TheirCountry: []DXCCEntity{"dl"}, Value: 2},
				{Property: MemberNumberProperty, Value: 3},
				{Value: 1},
			},
			myContinent:    Europa,
			myCountry:      "dl",
			theirContinent: Africa,
			theirCountry:   "zs",
			exchange: QSOExchange{
				MemberNumberProperty: "1234",
			},
			expected: []ScoringRule{
				{Property: MemberNumberProperty, Value: 3},
			},
		},
		{
			desc: "exchange does not contain a specified property",
			rules: []ScoringRule{
				{TheirCountry: []DXCCEntity{"dl"}, Value: 2},
				{Property: MemberNumberProperty, Value: 3},
				{Value: 1},
			},
			myContinent:    Europa,
			myCountry:      "dl",
			theirContinent: Africa,
			theirCountry:   "zs",
			exchange:       QSOExchange{},
			expected: []ScoringRule{
				{Value: 1},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			actual := filterScoringRules(tc.rules, true, tc.myContinent, tc.myCountry, tc.theirContinent, tc.theirCountry, tc.band, tc.exchange)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestCounter_Add_Points_Once(t *testing.T) {
	setup := Setup{}
	rules := Scoring{
		QSORules:    []ScoringRule{{Value: 1}},
		QSOBandRule: Once,
	}
	qsos := []QSO{
		{TheirCall: callsign.MustParse("DL1ABC")},
		{TheirCall: callsign.MustParse("DL2ABC")},
		{TheirCall: callsign.MustParse("DL1ABC")},
	}
	expectedScores := []QSOScore{
		{Points: 1, Duplicate: false},
		{Points: 1, Duplicate: false},
		{Points: 1, Duplicate: true},
	}
	expectedTotalScore := BandScore{Points: 2}

	counter := NewCounter(setup, rules)

	for i, qso := range qsos {
		actualScore := counter.Add(qso)
		assert.Equal(t, expectedScores[i], actualScore)
	}
	actualTotalScore := counter.TotalScore()
	assert.Equal(t, expectedTotalScore, actualTotalScore)
}

func TestCounter_Add_Points_OncePerBand(t *testing.T) {
	setup := Setup{}
	rules := Scoring{
		QSORules:    []ScoringRule{{Value: 1}},
		QSOBandRule: OncePerBand,
	}
	qsos := []QSO{
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band80m},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band80m},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band80m},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band40m},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band40m},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band40m},
	}
	expectedScores := []QSOScore{
		{Points: 1, Duplicate: false},
		{Points: 1, Duplicate: false},
		{Points: 1, Duplicate: true},
		{Points: 1, Duplicate: false},
		{Points: 1, Duplicate: false},
		{Points: 1, Duplicate: true},
	}
	expectedTotalScore := BandScore{Points: 4}

	counter := NewCounter(setup, rules)

	for i, qso := range qsos {
		actualScore := counter.Add(qso)
		assert.Equal(t, expectedScores[i], actualScore)
	}
	actualTotalScore := counter.TotalScore()
	assert.Equal(t, expectedTotalScore, actualTotalScore)
}
