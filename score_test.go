package conval

import (
	"testing"

	"github.com/ftl/hamradio/callsign"
	"github.com/stretchr/testify/assert"
)

func TestCounter_SimplestHappyPath(t *testing.T) {
	definition := Definition{
		Exchange: []ExchangeDefinition{{Fields: []ExchangeField{}}},
		Scoring: Scoring{
			QSORules:    []ScoringRule{{Value: 1}},
			QSOBandRule: OncePerBand,
			MultiRules:  []ScoringRule{{Value: 2}},
		},
	}
	setup := Setup{}
	qso := QSO{}

	counter := NewCounter(definition, setup, nil)

	qsoScore := counter.Add(qso)

	assert.Equal(t, 1, qsoScore.Points, "points")
	assert.Equal(t, 2, qsoScore.Multis, "multis")
}

func TestFilterScoringRules(t *testing.T) {
	tt := []struct {
		desc             string
		rules            []ScoringRule
		myContinent      Continent
		myCountry        DXCCEntity
		myPrefix         string
		theirContinent   Continent
		theirCountry     DXCCEntity
		theirPrefix      string
		band             ContestBand
		myExchange       QSOExchange
		theirExchange    QSOExchange
		onlyMostRelevant bool
		expected         []ScoringRule
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
			desc: "my continent not matching the excluded continent",
			rules: []ScoringRule{
				{MyContinent: []Continent{NotContinent, Europa}, Value: 2},
				{Value: 1},
			},
			myContinent: Africa,
			expected: []ScoringRule{
				{MyContinent: []Continent{NotContinent, Europa}, Value: 2},
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
			desc: "my country not matching the excluded country",
			rules: []ScoringRule{
				{MyCountry: []DXCCEntity{"not", "dl"}, Value: 2},
				{Value: 1},
			},
			myContinent: Europa,
			myCountry:   "f",
			expected: []ScoringRule{
				{MyCountry: []DXCCEntity{"not", "dl"}, Value: 2},
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
			desc: "their continent not matching the excluded continent",
			rules: []ScoringRule{
				{TheirContinent: []Continent{NotContinent, Africa, Antarctica, Asia, NorthAmerica, Oceania, SouthAmerica}, Value: 2},
				{Value: 1},
			},
			myContinent:    Europa,
			myCountry:      "dl",
			theirContinent: Europa,
			theirCountry:   "dl",
			expected: []ScoringRule{
				{TheirContinent: []Continent{NotContinent, Africa, Antarctica, Asia, NorthAmerica, Oceania, SouthAmerica}, Value: 2},
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
			desc: "their country not matching the excluded country",
			rules: []ScoringRule{
				{TheirCountry: []DXCCEntity{"not", "f"}, Value: 2},
				{Value: 1},
			},
			myContinent:    Europa,
			myCountry:      "dl",
			theirContinent: Europa,
			theirCountry:   "dl",
			expected: []ScoringRule{
				{TheirCountry: []DXCCEntity{"not", "f"}, Value: 2},
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
			theirExchange: QSOExchange{
				MemberNumberProperty: "1234",
			},
			onlyMostRelevant: true,
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
			theirExchange:  QSOExchange{},
			expected: []ScoringRule{
				{Value: 1},
			},
		},
		{
			desc: "use additional weight to override a more specific rule",
			rules: []ScoringRule{
				{MyContinent: []Continent{NorthAmerica}, TheirContinent: []Continent{NorthAmerica}, Value: 1},
				{TheirCountry: []DXCCEntity{SameCountry}, AdditionalWeight: 10, Value: 0},
			},
			myContinent:    NorthAmerica,
			myCountry:      "k",
			theirContinent: NorthAmerica,
			theirCountry:   "k",
			expected: []ScoringRule{
				{TheirCountry: []DXCCEntity{SameCountry}, AdditionalWeight: 10, Value: 0},
			},
		},
		{
			desc: "find multiple matching rules for multis",
			rules: []ScoringRule{
				{Property: CQZoneProperty, BandRule: OncePerBand, Value: 1},
				{Property: DXCCEntityProperty, BandRule: OncePerBand, Value: 1},
			},
			myContinent:    NorthAmerica,
			myCountry:      "k",
			theirContinent: NorthAmerica,
			theirCountry:   "k",
			theirExchange: QSOExchange{
				CQZoneProperty:     "5",
				DXCCEntityProperty: "k",
			},
			expected: []ScoringRule{
				{Property: CQZoneProperty, BandRule: OncePerBand, Value: 1},
				{Property: DXCCEntityProperty, BandRule: OncePerBand, Value: 1},
			},
		},
		{
			desc: "minimum age",
			rules: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Min: "23"},
					},
					Value: 1,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Min: "29"},
					},
					Value: 2,
				},
			},
			theirExchange: QSOExchange{
				AgeProperty: "28",
			},
			expected: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Min: "23"},
					},
					Value: 1,
				},
			},
		},
		{
			desc: "maximum age",
			rules: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Max: "23"},
					},
					Value: 1,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Max: "29"},
					},
					Value: 2,
				},
			},
			theirExchange: QSOExchange{
				AgeProperty: "28",
			},
			expected: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Max: "29"},
					},
					Value: 2,
				},
			},
		},
		{
			desc: "age in range",
			rules: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Min: "23", Max: "29"},
					},
					Value: 1,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Min: "13", Max: "19"},
					},
					Value: 2,
				},
			},
			theirExchange: QSOExchange{
				AgeProperty: "28",
			},
			expected: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Min: "23", Max: "29"},
					},
					Value: 1,
				},
			},
		},
		{
			desc: "age out of range",
			rules: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Min: "23", Max: "29"},
					},
					Value: 1,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Min: "13", Max: "19"},
					},
					Value: 2,
				},
			},
			theirExchange: QSOExchange{
				AgeProperty: "20",
			},
			expected: []ScoringRule{},
		},
		{
			desc: "age below gap",
			rules: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Min: "23"},
					},
					Value: 1,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Max: "19"},
					},
					Value: 2,
				},
			},
			theirExchange: QSOExchange{
				AgeProperty: "18",
			},
			expected: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: AgeProperty, Max: "19"},
					},
					Value: 2,
				},
			},
		},
		{
			desc: "only my class equals",
			rules: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, MyValue: "my"},
					},
					Value: 1,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, TheirValue: "their"},
					},
					Value: 2,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, MyValue: "my", TheirValue: "their"},
					},
					Value: 3,
				},
			},
			myExchange: QSOExchange{
				ClassProperty: "my",
			},
			theirExchange: QSOExchange{},
			expected: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, MyValue: "my"},
					},
					Value: 1,
				},
			},
		},
		{
			desc: "no values matches all",
			rules: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty},
					},
					Value: 4,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, MyValue: "my"},
					},
					Value: 1,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, TheirValue: "their"},
					},
					Value: 2,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, MyValue: "my", TheirValue: "their"},
					},
					Value: 3,
				},
			},
			myExchange: QSOExchange{
				ClassProperty: "my",
			},
			theirExchange: QSOExchange{},
			expected: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty},
					},
					Value: 4,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, MyValue: "my"},
					},
					Value: 1,
				},
			},
		},
		{
			desc: "only my class equals",
			rules: []ScoringRule{
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, MyValue: "v1", TheirValue: "v1"},
					},
					Value: 1,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, MyValue: "v2", TheirValue: "v2"},
					},
					Value: 2,
				},
				{
					Value: 3,
				},
			},
			myExchange: QSOExchange{
				ClassProperty: "v1",
			},
			theirExchange: QSOExchange{
				ClassProperty: "v2",
			},
			expected: []ScoringRule{
				{
					Value: 3,
				},
			},
		},
		{
			desc: "only the most relevant per property",
			rules: []ScoringRule{
				{
					Property: DXCCEntityProperty,
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, MyValue: "v1"},
					},
					Value: 1,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, TheirValue: "v2"},
					},
					Value: 2,
				},
				{
					Property: DXCCEntityProperty,
					Value:    3,
				},
			},
			myExchange: QSOExchange{
				ClassProperty: "v1",
			},
			theirExchange: QSOExchange{
				ClassProperty:      "v2",
				DXCCEntityProperty: "dl",
			},
			expected: []ScoringRule{
				{
					Property: DXCCEntityProperty,
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, MyValue: "v1"},
					},
					Value: 1,
				},
				{
					PropertyConstraints: []PropertyConstraint{
						{Name: ClassProperty, TheirValue: "v2"},
					},
					Value: 2,
				},
			},
		},
		{
			desc: "only my special prefix included (1)",
			rules: []ScoringRule{
				{
					MyCountry: []DXCCEntity{"k", "ve"},
					MyPrefix:  []string{"kh6", "kl7", "cy9", "cy0"},
					Value:     3,
				},
				{
					MyCountry: []DXCCEntity{"not", "k", "ve"},
					Value:     3,
				},
			},
			myCountry: "k",
			myPrefix:  "kh6",
			expected: []ScoringRule{
				{
					MyCountry: []DXCCEntity{"k", "ve"},
					MyPrefix:  []string{"kh6", "kl7", "cy9", "cy0"},
					Value:     3,
				},
			},
		},
		{
			desc: "only my special prefix included (2)",
			rules: []ScoringRule{
				{
					MyCountry: []DXCCEntity{"k", "ve"},
					MyPrefix:  []string{"kh6", "kl7", "cy9", "cy0"},
					Value:     3,
				},
				{
					MyCountry: []DXCCEntity{"not", "k", "ve"},
					Value:     3,
				},
			},
			myCountry: "k",
			myPrefix:  "w6",
			expected:  []ScoringRule{},
		},
		{
			desc: "only my special prefix excluded (1)",
			rules: []ScoringRule{
				{
					MyCountry: []DXCCEntity{"k", "ve"},
					MyPrefix:  []string{"not", "kh6", "kl7", "cy9", "cy0"},
					Value:     3,
				},
				{
					MyCountry: []DXCCEntity{"not", "k", "ve"},
					Value:     3,
				},
			},
			myCountry: "k",
			myPrefix:  "kh6",
			expected:  []ScoringRule{},
		},
		{
			desc: "only my special prefix excluded (2)",
			rules: []ScoringRule{
				{
					MyCountry: []DXCCEntity{"k", "ve"},
					MyPrefix:  []string{"not", "kh6", "kl7", "cy9", "cy0"},
					Value:     3,
				},
				{
					MyCountry: []DXCCEntity{"not", "k", "ve"},
					Value:     3,
				},
			},
			myCountry: "k",
			myPrefix:  "w",
			expected: []ScoringRule{
				{
					MyCountry: []DXCCEntity{"k", "ve"},
					MyPrefix:  []string{"not", "kh6", "kl7", "cy9", "cy0"},
					Value:     3,
				},
			},
		},
		{
			desc: "only their special prefix included (1)",
			rules: []ScoringRule{
				{
					TheirCountry: []DXCCEntity{"k", "ve"},
					TheirPrefix:  []string{"kh6", "kl7", "cy9", "cy0"},
					Value:        3,
				},
				{
					TheirCountry: []DXCCEntity{"not", "k", "ve"},
					Value:        3,
				},
			},
			theirCountry: "k",
			theirPrefix:  "kh6",
			expected: []ScoringRule{
				{
					TheirCountry: []DXCCEntity{"k", "ve"},
					TheirPrefix:  []string{"kh6", "kl7", "cy9", "cy0"},
					Value:        3,
				},
			},
		},
		{
			desc: "only their special prefix included (2)",
			rules: []ScoringRule{
				{
					TheirCountry: []DXCCEntity{"k", "ve"},
					TheirPrefix:  []string{"kh6", "kl7", "cy9", "cy0"},
					Value:        3,
				},
				{
					TheirCountry: []DXCCEntity{"not", "k", "ve"},
					Value:        3,
				},
			},
			theirCountry: "k",
			theirPrefix:  "w6",
			expected:     []ScoringRule{},
		},
		{
			desc: "only their special prefix excluded (1)",
			rules: []ScoringRule{
				{
					TheirCountry: []DXCCEntity{"k", "ve"},
					TheirPrefix:  []string{"not", "kh6", "kl7", "cy9", "cy0"},
					Value:        3,
				},
				{
					TheirCountry: []DXCCEntity{"not", "k", "ve"},
					Value:        3,
				},
			},
			theirCountry: "k",
			theirPrefix:  "kh6",
			expected:     []ScoringRule{},
		},
		{
			desc: "only their special prefix excluded (2)",
			rules: []ScoringRule{
				{
					TheirCountry: []DXCCEntity{"k", "ve"},
					TheirPrefix:  []string{"not", "kh6", "kl7", "cy9", "cy0"},
					Value:        3,
				},
				{
					TheirCountry: []DXCCEntity{"not", "k", "ve"},
					Value:        3,
				},
			},
			theirCountry: "k",
			theirPrefix:  "w6",
			expected: []ScoringRule{
				{
					TheirCountry: []DXCCEntity{"k", "ve"},
					TheirPrefix:  []string{"not", "kh6", "kl7", "cy9", "cy0"},
					Value:        3,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			getMyProperty := func(property Property) string {
				return tc.myExchange[property]
			}
			getTheirProperty := func(property Property) string {
				return tc.theirExchange[property]
			}
			counter := Counter{}

			actual := counter.filterScoringRules(tc.rules, tc.onlyMostRelevant, tc.myContinent, tc.myCountry, tc.myPrefix, tc.theirContinent, tc.theirCountry, tc.theirPrefix, tc.band, getMyProperty, getTheirProperty)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestFilterExchangeFields(t *testing.T) {
	tt := []struct {
		desc           string
		definitions    []ExchangeDefinition
		myContinent    Continent
		myCountry      DXCCEntity
		theirContinent Continent
		theirCountry   DXCCEntity
		expected       []ExchangeField
	}{
		{
			desc: "single definition, no constraints",
			definitions: []ExchangeDefinition{
				{Fields: []ExchangeField{{RSTProperty}, {SerialNumberProperty}}},
			},
			expected: []ExchangeField{{RSTProperty}, {SerialNumberProperty}},
		},
		{
			desc: "one country-specific, one general, get specific",
			definitions: []ExchangeDefinition{
				{TheirCountry: []DXCCEntity{"f"}, Fields: []ExchangeField{{RSTProperty}, {SerialNumberProperty}}},
				{Fields: []ExchangeField{{RSTProperty}, {CQZoneProperty}}},
			},
			theirCountry: "f",
			expected:     []ExchangeField{{RSTProperty}, {SerialNumberProperty}},
		},
		{
			desc: "one country-specific, one general, get general",
			definitions: []ExchangeDefinition{
				{TheirCountry: []DXCCEntity{"f"}, Fields: []ExchangeField{{RSTProperty}, {SerialNumberProperty}}},
				{Fields: []ExchangeField{{RSTProperty}, {CQZoneProperty}}},
			},
			theirCountry: "dl",
			expected:     []ExchangeField{{RSTProperty}, {CQZoneProperty}},
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			counter := Counter{}

			actual := counter.filterExchangeFields(tc.definitions, tc.myContinent, tc.myCountry, tc.theirContinent, tc.theirCountry, "")

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestCounter_Add_Points_Once(t *testing.T) {
	definition := Definition{
		Exchange: []ExchangeDefinition{{Fields: []ExchangeField{}}},
		Scoring: Scoring{
			QSORules:    []ScoringRule{{Value: 1}},
			QSOBandRule: Once,
			MultiRules:  []ScoringRule{{Property: GenericTextProperty, Value: 2}},
		},
	}
	setup := Setup{}
	qsos := []QSO{
		{TheirCall: callsign.MustParse("DL1ABC"), TheirExchange: QSOExchange{GenericTextProperty: "1"}},
		{TheirCall: callsign.MustParse("DL2ABC"), TheirExchange: QSOExchange{GenericTextProperty: "2"}},
		{TheirCall: callsign.MustParse("DL1ABC"), TheirExchange: QSOExchange{GenericTextProperty: "3"}},
	}
	expectedScores := []QSOScore{
		{Points: 1, Multis: 2, Duplicate: false, MultiValues: map[Property]string{GenericTextProperty: "1"}, MultiBandAndMode: map[Property]BandAndMode{GenericTextProperty: {}}},
		{Points: 1, Multis: 2, Duplicate: false, MultiValues: map[Property]string{GenericTextProperty: "2"}, MultiBandAndMode: map[Property]BandAndMode{GenericTextProperty: {}}},
		{Points: 1, Multis: 2, Duplicate: true, MultiValues: map[Property]string{GenericTextProperty: "3"}, MultiBandAndMode: map[Property]BandAndMode{GenericTextProperty: {}}},
	}
	expectedTotalScore := BandScore{QSOs: 3, Points: 2, Multis: 4}

	counter := NewCounter(definition, setup, nil)

	for i, qso := range qsos {
		actualScore := counter.Add(qso)
		assert.Equal(t, expectedScores[i], actualScore, i)
	}
	actualTotalScore := counter.TotalScore()
	assert.Equal(t, expectedTotalScore, actualTotalScore, "total score")
	assert.Equal(t, expectedTotalScore.Points*expectedTotalScore.Multis, counter.Total(actualTotalScore))
}

func TestCounter_Add_Points_OncePerBand(t *testing.T) {
	definition := Definition{
		Exchange: []ExchangeDefinition{{Fields: []ExchangeField{}}},
		Scoring: Scoring{
			QSORules:    []ScoringRule{{Value: 1}},
			QSOBandRule: OncePerBand,
		},
	}
	setup := Setup{}
	qsos := []QSO{
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band80m},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band80m},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band80m},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band40m},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band40m},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band40m},
	}
	expectedScores := []QSOScore{
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
	}
	expectedTotalScore := BandScore{QSOs: 6, Points: 4}

	counter := NewCounter(definition, setup, nil)

	for i, qso := range qsos {
		actualScore := counter.Add(qso)
		assert.Equal(t, expectedScores[i], actualScore)
	}
	actualTotalScore := counter.TotalScore()
	assert.Equal(t, expectedTotalScore, actualTotalScore)
}

func TestCounter_Add_Points_OncePerBandAndMode(t *testing.T) {
	definition := Definition{
		Exchange: []ExchangeDefinition{{Fields: []ExchangeField{}}},
		Scoring: Scoring{
			QSORules:    []ScoringRule{{Value: 1}},
			QSOBandRule: OncePerBandAndMode,
		},
	}
	setup := Setup{}
	qsos := []QSO{
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band80m, Mode: ModeCW},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band80m, Mode: ModeCW},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band80m, Mode: ModeCW},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band80m, Mode: ModeSSB},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band80m, Mode: ModeSSB},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band80m, Mode: ModeSSB},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band40m, Mode: ModeCW},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band40m, Mode: ModeCW},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band40m, Mode: ModeCW},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band40m, Mode: ModeSSB},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band40m, Mode: ModeSSB},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band40m, Mode: ModeSSB},
	}
	expectedScores := []QSOScore{
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
	}
	expectedTotalScore := BandScore{QSOs: 12, Points: 8}

	counter := NewCounter(definition, setup, nil)

	for i, qso := range qsos {
		actualScore := counter.Add(qso)
		assert.Equal(t, expectedScores[i], actualScore)
	}
	actualTotalScore := counter.TotalScore()
	assert.Equal(t, expectedTotalScore, actualTotalScore)
}

func TestCounter_Add_Multis_Once(t *testing.T) {
	definition := Definition{
		Exchange: []ExchangeDefinition{{Fields: []ExchangeField{}}},
		Scoring: Scoring{
			MultiRules: []ScoringRule{{Property: "cq_zone", BandRule: Once, Value: 1}},
		},
	}
	setup := Setup{}
	qsos := []QSO{
		{TheirCall: callsign.MustParse("DL1ABC"), TheirExchange: QSOExchange{"cq_zone": "14"}},
		{TheirCall: callsign.MustParse("AB1C"), TheirExchange: QSOExchange{"cq_zone": "5"}},
		{TheirCall: callsign.MustParse("DL2ABC"), TheirExchange: QSOExchange{"cq_zone": "14"}},
	}
	expectedScores := []QSOScore{
		{Multis: 1, MultiValues: map[Property]string{"cq_zone": "14"}, MultiBandAndMode: map[Property]BandAndMode{"cq_zone": {BandAll, ModeALL}}},
		{Multis: 1, MultiValues: map[Property]string{"cq_zone": "5"}, MultiBandAndMode: map[Property]BandAndMode{"cq_zone": {BandAll, ModeALL}}},
		{Multis: 0, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
	}
	expectedTotalScore := BandScore{QSOs: 3, Multis: 2}

	counter := NewCounter(definition, setup, nil)

	for i, qso := range qsos {
		actualScore := counter.Add(qso)
		assert.Equal(t, expectedScores[i], actualScore, "%d", i)
	}
	actualTotalScore := counter.TotalScore()
	assert.Equal(t, expectedTotalScore, actualTotalScore)
}

func TestCounter_Add_Multis_OncePerBand(t *testing.T) {
	definition := Definition{
		Exchange: []ExchangeDefinition{{Fields: []ExchangeField{}}},
		Scoring: Scoring{
			QSOBandRule: OncePerBand,
			MultiRules:  []ScoringRule{{Property: "cq_zone", BandRule: OncePerBand, Value: 1}},
		},
	}
	setup := Setup{}
	qsos := []QSO{
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band80m, TheirExchange: QSOExchange{"cq_zone": "14"}},
		{TheirCall: callsign.MustParse("AB1C"), Band: Band80m, TheirExchange: QSOExchange{"cq_zone": "5"}},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band80m, TheirExchange: QSOExchange{"cq_zone": "14"}},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band40m, TheirExchange: QSOExchange{"cq_zone": "14"}},
		{TheirCall: callsign.MustParse("AB1C"), Band: Band40m, TheirExchange: QSOExchange{"cq_zone": "5"}},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band40m, TheirExchange: QSOExchange{"cq_zone": "14"}},
	}
	expectedScores := []QSOScore{
		{Multis: 1, MultiValues: map[Property]string{"cq_zone": "14"}, MultiBandAndMode: map[Property]BandAndMode{"cq_zone": {Band80m, ModeALL}}},
		{Multis: 1, MultiValues: map[Property]string{"cq_zone": "5"}, MultiBandAndMode: map[Property]BandAndMode{"cq_zone": {Band80m, ModeALL}}},
		{Multis: 0, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Multis: 1, MultiValues: map[Property]string{"cq_zone": "14"}, MultiBandAndMode: map[Property]BandAndMode{"cq_zone": {Band40m, ModeALL}}},
		{Multis: 1, MultiValues: map[Property]string{"cq_zone": "5"}, MultiBandAndMode: map[Property]BandAndMode{"cq_zone": {Band40m, ModeALL}}},
		{Multis: 0, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
	}
	expectedTotalScore := BandScore{QSOs: 6, Multis: 4}

	counter := NewCounter(definition, setup, nil)

	for i, qso := range qsos {
		actualScore := counter.Add(qso)
		assert.Equal(t, expectedScores[i], actualScore, "%d", i)
	}
	actualTotalScore := counter.TotalScore()
	assert.Equal(t, expectedTotalScore, actualTotalScore)
}

func TestCounter_BandsPerMulti(t *testing.T) {
	definition := Definition{
		Exchange: []ExchangeDefinition{{Fields: []ExchangeField{}}},
		Scoring: Scoring{
			QSOBandRule: OncePerBand,
			MultiRules:  []ScoringRule{{Property: CQZoneProperty, BandRule: OncePerBand, Value: 1}},
		},
	}
	setup := Setup{}
	qsos := []QSO{
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band80m, TheirExchange: QSOExchange{"cq_zone": "14"}},
		{TheirCall: callsign.MustParse("AB1C"), Band: Band80m, TheirExchange: QSOExchange{"cq_zone": "5"}},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band80m, TheirExchange: QSOExchange{"cq_zone": "14"}},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band40m, TheirExchange: QSOExchange{"cq_zone": "14"}},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band40m, TheirExchange: QSOExchange{"cq_zone": "14"}},
	}

	counter := NewCounter(definition, setup, nil)
	for _, qso := range qsos {
		counter.Add(qso)
	}

	bandsPerZone8 := counter.BandsPerMulti(CQZoneProperty, "8")
	assert.ElementsMatch(t, []ContestBand{}, bandsPerZone8)
	bandsPerZone5 := counter.BandsPerMulti(CQZoneProperty, "5")
	assert.ElementsMatch(t, []ContestBand{Band80m}, bandsPerZone5)
	bandsPerZone14 := counter.BandsPerMulti(CQZoneProperty, "14")
	assert.ElementsMatch(t, []ContestBand{Band80m, Band40m}, bandsPerZone14)
}
