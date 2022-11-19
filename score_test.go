package conval

import (
	"testing"

	"github.com/ftl/hamradio/callsign"
	"github.com/stretchr/testify/assert"
)

func TestCounter_SimplestHappyPath(t *testing.T) {
	setup := Setup{}
	exchange := []ExchangeDefinition{{Fields: []ExchangeField{}}}
	rules := Scoring{
		QSORules:    []ScoringRule{{Value: 1}},
		QSOBandRule: OncePerBand,
		MultiRules:  []ScoringRule{{Value: 1}, {Value: 2}},
	}
	qso := QSO{}

	counter := NewCounter(setup, exchange, rules)

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
			exchange: QSOExchange{
				CQZoneProperty:     "5",
				DXCCEntityProperty: "k",
			},
			expected: []ScoringRule{
				{Property: CQZoneProperty, BandRule: OncePerBand, Value: 1},
				{Property: DXCCEntityProperty, BandRule: OncePerBand, Value: 1},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			getProperty := func(property Property) string {
				return tc.exchange[property]
			}

			actual := filterScoringRules(tc.rules, true, tc.myContinent, tc.myCountry, tc.theirContinent, tc.theirCountry, tc.band, getProperty)

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
				{Fields: []ExchangeField{{TheirRSTProperty}, {SerialNumberProperty}}},
			},
			expected: []ExchangeField{{TheirRSTProperty}, {SerialNumberProperty}},
		},
		{
			desc: "one country-specific, one general, get specific",
			definitions: []ExchangeDefinition{
				{TheirCountry: []DXCCEntity{"f"}, Fields: []ExchangeField{{TheirRSTProperty}, {SerialNumberProperty}}},
				{Fields: []ExchangeField{{TheirRSTProperty}, {CQZoneProperty}}},
			},
			theirCountry: "f",
			expected:     []ExchangeField{{TheirRSTProperty}, {SerialNumberProperty}},
		},
		{
			desc: "one country-specific, one general, get general",
			definitions: []ExchangeDefinition{
				{TheirCountry: []DXCCEntity{"f"}, Fields: []ExchangeField{{TheirRSTProperty}, {SerialNumberProperty}}},
				{Fields: []ExchangeField{{TheirRSTProperty}, {CQZoneProperty}}},
			},
			theirCountry: "dl",
			expected:     []ExchangeField{{TheirRSTProperty}, {CQZoneProperty}},
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			actual := filterExchangeFields(tc.definitions, tc.myContinent, tc.myCountry, tc.theirContinent, tc.theirCountry)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestCounter_Add_Points_Once(t *testing.T) {
	setup := Setup{}
	exchange := []ExchangeDefinition{{Fields: []ExchangeField{}}}
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
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 0, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
	}
	expectedTotalScore := BandScore{QSOs: 3, Points: 2}

	counter := NewCounter(setup, exchange, rules)

	for i, qso := range qsos {
		actualScore := counter.Add(qso)
		assert.Equal(t, expectedScores[i], actualScore)
	}
	actualTotalScore := counter.TotalScore()
	assert.Equal(t, expectedTotalScore, actualTotalScore)
}

func TestCounter_Add_Points_OncePerBand(t *testing.T) {
	setup := Setup{}
	exchange := []ExchangeDefinition{{Fields: []ExchangeField{}}}
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
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 0, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 0, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
	}
	expectedTotalScore := BandScore{QSOs: 6, Points: 4}

	counter := NewCounter(setup, exchange, rules)

	for i, qso := range qsos {
		actualScore := counter.Add(qso)
		assert.Equal(t, expectedScores[i], actualScore)
	}
	actualTotalScore := counter.TotalScore()
	assert.Equal(t, expectedTotalScore, actualTotalScore)
}

func TestCounter_Add_Points_OncePerBandAndMode(t *testing.T) {
	setup := Setup{}
	exchange := []ExchangeDefinition{{Fields: []ExchangeField{}}}
	rules := Scoring{
		QSORules:    []ScoringRule{{Value: 1}},
		QSOBandRule: OncePerBandAndMode,
	}
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
		{Points: 0, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 0, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 0, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 1, Duplicate: false, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
		{Points: 0, Duplicate: true, MultiValues: map[Property]string{}, MultiBandAndMode: map[Property]BandAndMode{}},
	}
	expectedTotalScore := BandScore{QSOs: 12, Points: 8}

	counter := NewCounter(setup, exchange, rules)

	for i, qso := range qsos {
		actualScore := counter.Add(qso)
		assert.Equal(t, expectedScores[i], actualScore)
	}
	actualTotalScore := counter.TotalScore()
	assert.Equal(t, expectedTotalScore, actualTotalScore)
}

func TestCounter_Add_Multis_Once(t *testing.T) {
	setup := Setup{}
	exchange := []ExchangeDefinition{{Fields: []ExchangeField{}}}
	rules := Scoring{
		MultiRules: []ScoringRule{{Property: "cq_zone", BandRule: Once, Value: 1}},
	}
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

	counter := NewCounter(setup, exchange, rules)

	for i, qso := range qsos {
		actualScore := counter.Add(qso)
		assert.Equal(t, expectedScores[i], actualScore, "%d", i)
	}
	actualTotalScore := counter.TotalScore()
	assert.Equal(t, expectedTotalScore, actualTotalScore)
}

func TestCounter_Add_Multis_OncePerBand(t *testing.T) {
	setup := Setup{}
	exchange := []ExchangeDefinition{{Fields: []ExchangeField{}}}
	rules := Scoring{
		QSOBandRule: OncePerBand,
		MultiRules:  []ScoringRule{{Property: "cq_zone", BandRule: OncePerBand, Value: 1}},
	}
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

	counter := NewCounter(setup, exchange, rules)

	for i, qso := range qsos {
		actualScore := counter.Add(qso)
		assert.Equal(t, expectedScores[i], actualScore, "%d", i)
	}
	actualTotalScore := counter.TotalScore()
	assert.Equal(t, expectedTotalScore, actualTotalScore)
}

func TestCounter_BandsPerMulti(t *testing.T) {
	setup := Setup{}
	exchange := []ExchangeDefinition{{Fields: []ExchangeField{}}}
	rules := Scoring{
		QSOBandRule: OncePerBand,
		MultiRules:  []ScoringRule{{Property: CQZoneProperty, BandRule: OncePerBand, Value: 1}},
	}
	qsos := []QSO{
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band80m, TheirExchange: QSOExchange{"cq_zone": "14"}},
		{TheirCall: callsign.MustParse("AB1C"), Band: Band80m, TheirExchange: QSOExchange{"cq_zone": "5"}},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band80m, TheirExchange: QSOExchange{"cq_zone": "14"}},
		{TheirCall: callsign.MustParse("DL1ABC"), Band: Band40m, TheirExchange: QSOExchange{"cq_zone": "14"}},
		{TheirCall: callsign.MustParse("DL2ABC"), Band: Band40m, TheirExchange: QSOExchange{"cq_zone": "14"}},
	}

	counter := NewCounter(setup, exchange, rules)
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
