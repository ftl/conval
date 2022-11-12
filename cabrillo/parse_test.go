package cabrillo

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/ftl/hamradio/callsign"
	"github.com/ftl/hamradio/locator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead_FromStartToEnd(t *testing.T) {
	tt := []struct {
		desc    string
		value   string
		invalid bool
	}{
		{
			desc:  "happy path",
			value: "START-OF-LOG: 3.0\nX-CUSTOM: some custom content\nEND-OF-LOG:\n",
		},
		{
			desc:  "no trailing newline",
			value: "START-OF-LOG: 3.0\nEND-OF-LOG:",
		},
		{
			desc:  "leading and trailing lines",
			value: "This line is ignored\n#this line either\nSTART-OF-LOG: 3.0\nEND-OF-LOG:\nas well as this line",
		},
		{
			desc:    "empty",
			value:   "",
			invalid: true,
		},
		{
			desc:    "no start",
			value:   "X-CUSTOM: some custom content\nEND-OF-LOG:\n",
			invalid: true,
		},
		{
			desc:    "no end",
			value:   "START-OF-LOG: 2.0\nX-CUSTOM: some custom content\n",
			invalid: true,
		},
		{
			desc:    "double start",
			value:   "START-OF-LOG: 3.0\nSTART-OF-LOG: 3.0\nX-CUSTOM: some custom content\nEND-OF-LOG:\n",
			invalid: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			buffer := bytes.NewBufferString(tc.value)
			actualLog, err := Read(buffer)
			if tc.invalid {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, actualLog)
			}
		})
	}
}

func TestParser_ParseAllTags(t *testing.T) {
	lines := []string{
		"START-OF-LOG: 3.0",
		"CALLSIGN: DL1ABC",
		"CONTEST: WAG",
		"CATEGORY-ASSISTED: assisted",
		"CATEGORY-BAND: ALL",
		"CATEGORY-MODE: CW",
		"CATEGORY-OPERATOR: SINGLE-OP",
		"CATEGORY-POWER: LOW",
		"CATEGORY-STATION: FIXED",
		"CATEGORY-TIME: 24-HOURS",
		"CATEGORY-TRANSMITTER: ONE",
		"CATEGORY-OVERLAY: CLASSIC",
		"CERTIFICATE: yes",
		"CLAIMED-SCORE: 123",
		"CLUB: BCC",
		"CREATED-BY: conval - the CONtest eVALuator",
		"EMAIL: conval@example.com",
		"GRID-LOCATOR: jn59",
		"LOCATION: DX",
		"NAME: Constantin Valberg",
		"ADDRESS: beside the big river",
		"ADDRESS-CITY: Musterstadt",
		"ADDRESS-STATE-PROVINCE: Bavaria",
		"ADDRESS-POSTALCODE: 80123",
		"ADDRESS-COUNTRY: Germany",
		"OPERATORS: @DL1ABC",
		"OFFTIME: 2002-03-22 0300 2002-03-22 0743",
		"SOAPBOX: This is an example that contains all officially",
		"SOAPBOX: defined Cabrillo tags.",
		"SOAPBOX:",
		"SOAPBOX: With an extra line and an empty line in between.",
		"QSO:  3559 CW 1999-03-06 0711 DL1ABC           599 B01    W1AW           599 001     0",
		"X-QSO:  3559 CW 1999-03-06 0712 DL1ABC           599 B01    N5KO           599 001     1",
		"X-CUSTOM: this is content in a custom tag",
		"END-OF-LOG:",
	}

	actualLog := NewLog()
	parser := newParser(actualLog)

	for _, line := range lines {
		err := parser.AddLine(line)
		assert.NoError(t, err)
	}
	assert.NoError(t, parser.CheckComplete())

	assert.Equal(t, "3.0", actualLog.CabrilloVersion, "cabrillo version")
	assert.Equal(t, callsign.MustParse("DL1ABC"), actualLog.Callsign, "callsign")
	assert.Equal(t, ContestIdentifier("WAG"), actualLog.Contest, "contest identifier")
	assert.Equal(t, Assisted, actualLog.Category.Assisted, "category assisted")
	assert.Equal(t, BandAll, actualLog.Category.Band, "category band")
	assert.Equal(t, ModeCW, actualLog.Category.Mode, "category mode")
	assert.Equal(t, SingleOperator, actualLog.Category.Operator, "category operator")
	assert.Equal(t, LowPower, actualLog.Category.Power, "category power")
	assert.Equal(t, FixedStation, actualLog.Category.Station, "category station")
	assert.Equal(t, Hours24, actualLog.Category.Time, "category time")
	assert.Equal(t, OneTransmitter, actualLog.Category.Transmitter, "category transmitter")
	assert.Equal(t, ClassicOverlay, actualLog.Category.Overlay, "category overlay")
	assert.True(t, actualLog.Certificate, "certificate")
	assert.Equal(t, 123, actualLog.ClaimedScore, "claimed score")
	assert.Equal(t, "BCC", actualLog.Club, "club")
	assert.Equal(t, "conval - the CONtest eVALuator", actualLog.CreatedBy, "created by")
	assert.Equal(t, "conval@example.com", actualLog.Email, "email")
	assert.Equal(t, locator.MustParse("JN59"), actualLog.GridLocator, "grid locator")
	assert.Equal(t, "DX", actualLog.Location, "location")
	assert.Equal(t, "Constantin Valberg", actualLog.Name, "name")
	assert.Equal(t, "beside the big river", actualLog.Address.Text, "address text")
	assert.Equal(t, "Musterstadt", actualLog.Address.City, "address city")
	assert.Equal(t, "Bavaria", actualLog.Address.StateProvince, "address state/province")
	assert.Equal(t, "80123", actualLog.Address.Postalcode, "address postalcode")
	assert.Equal(t, "Germany", actualLog.Address.Country, "address country")
	assert.Equal(t, []callsign.Callsign{callsign.MustParse("DL1ABC")}, actualLog.Operators, "operators")
	assert.Equal(t, callsign.MustParse("DL1ABC"), actualLog.Host, "host callsign")
	assert.Equal(t, time.Date(2002, time.March, 22, 3, 0, 0, 0, time.UTC), actualLog.Offtime.Begin, "offtime begin")
	assert.Equal(t, time.Date(2002, time.March, 22, 7, 43, 0, 0, time.UTC), actualLog.Offtime.End, "offtime end")
	assert.Equal(t, "This is an example that contains all officially\ndefined Cabrillo tags.\n\nWith an extra line and an empty line in between.", actualLog.Soapbox, "soapbox")
	assert.Equal(t, []QSO{
		{
			Frequency: QSOFrequency("3559"),
			Mode:      QSOModeCW,
			Timestamp: time.Date(1999, time.March, 6, 7, 11, 0, 0, time.UTC),
			Sent: QSOInfo{
				Call:     callsign.MustParse("DL1ABC"),
				RST:      "599",
				Exchange: []string{"B01"},
			},
			Received: QSOInfo{
				Call:     callsign.MustParse("W1AW"),
				RST:      "599",
				Exchange: []string{"001"},
			},
		},
	}, actualLog.QSOData, "qso data")
	assert.Equal(t, []QSO{
		{
			Frequency: QSOFrequency("3559"),
			Mode:      QSOModeCW,
			Timestamp: time.Date(1999, time.March, 6, 7, 12, 0, 0, time.UTC),
			Sent: QSOInfo{
				Call:     callsign.MustParse("DL1ABC"),
				RST:      "599",
				Exchange: []string{"B01"},
			},
			Received: QSOInfo{
				Call:     callsign.MustParse("N5KO"),
				RST:      "599",
				Exchange: []string{"001"},
			},
			Transmitter: 1,
		},
	}, actualLog.IgnoredQSOs, "ignored qsos")
}

func TestParseOperators(t *testing.T) {
	tt := []struct {
		desc              string
		value             string
		expectedOperators []string
		expectedHost      string
		invalid           bool
	}{
		{
			desc:              "",
			value:             "",
			expectedOperators: []string{},
		},
		{
			desc:              "single operator, no host",
			value:             "DL1ABC",
			expectedOperators: []string{"DL1ABC"},
		},
		{
			desc:              "single operator and host",
			value:             "@DL1ABC",
			expectedOperators: []string{"DL1ABC"},
			expectedHost:      "DL1ABC",
		},
		{
			desc:              "multiple operators, no host, spaces",
			value:             "DL1ABC DL2ABC DL3ABC",
			expectedOperators: []string{"DL1ABC", "DL2ABC", "DL3ABC"},
		},
		{
			desc:              "multiple operators, no host, kommas",
			value:             "DL1ABC, DL2ABC ,DL3ABC , DL4ABC",
			expectedOperators: []string{"DL1ABC", "DL2ABC", "DL3ABC", "DL4ABC"},
		},
		{
			desc:              "multiple operators, host, kommas",
			value:             "DL1ABC, DL2ABC ,DL3ABC , @DL4ABC",
			expectedOperators: []string{"DL1ABC", "DL2ABC", "DL3ABC", "DL4ABC"},
			expectedHost:      "DL4ABC",
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			expectedOperators := make([]callsign.Callsign, len(tc.expectedOperators))
			for i, op := range tc.expectedOperators {
				expectedOperators[i] = callsign.MustParse(op)
			}
			var expectedHost callsign.Callsign
			if tc.expectedHost != "" {
				expectedHost = callsign.MustParse(tc.expectedHost)
			}
			log := NewLog()
			parser := tagParsers[OperatorsTag]
			err := parser.Parse(log, tc.value)
			if tc.invalid {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedOperators, log.Operators)
				assert.Equal(t, expectedHost, log.Host)
			}
		})
	}
}

func TestParseOfftime(t *testing.T) {
	tt := []struct {
		desc          string
		value         string
		expectedBegin time.Time
		expectedEnd   time.Time
		invalid       bool
	}{
		{
			desc:          "empty",
			value:         "",
			expectedBegin: time.Time{},
			expectedEnd:   time.Time{},
		},
		{
			desc:          "happy path",
			value:         "2022-11-12 1000 2022-11-12 1200",
			expectedBegin: time.Date(2022, time.November, 12, 10, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2022, time.November, 12, 12, 0, 0, 0, time.UTC),
		},
		{
			desc:          "more space between timestamps",
			value:         "2022-11-12 1000   2022-11-12 1200",
			expectedBegin: time.Date(2022, time.November, 12, 10, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2022, time.November, 12, 12, 0, 0, 0, time.UTC),
		},
		{
			desc:    "more space before time",
			value:   "2022-11-12  1000 2022-11-12  1200",
			invalid: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			log := NewLog()
			parser := tagParsers[OfftimeTag]
			err := parser.Parse(log, tc.value)
			if tc.invalid {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBegin, log.Offtime.Begin)
				assert.Equal(t, tc.expectedEnd, log.Offtime.End)
			}
		})
	}
}

func TestParseQSO(t *testing.T) {
	tt := []struct {
		desc     string
		value    string
		expected QSO
		invalid  bool
	}{
		{
			desc:    "empty",
			value:   "",
			invalid: true,
		},
		{
			desc:  "happy path without transmitter column",
			value: "3799 PH 1999-03-06 0711 HC8N           59 700    W1AW           59 CT",
			expected: QSO{
				Frequency: QSOFrequency("3799"),
				Mode:      QSOModePhone,
				Timestamp: time.Date(1999, time.March, 6, 7, 11, 0, 0, time.UTC),
				Sent: QSOInfo{
					Call:     callsign.MustParse("HC8N"),
					RST:      "59",
					Exchange: []string{"700"},
				},
				Received: QSOInfo{
					Call:     callsign.MustParse("W1AW"),
					RST:      "59",
					Exchange: []string{"CT"},
				},
			},
		},
		{
			desc:  "happy path with transmitter column",
			value: "3799 PH 1999-03-06 0711 HC8N           59 700    W1AW           59 CT     1",
			expected: QSO{
				Frequency: QSOFrequency("3799"),
				Mode:      QSOModePhone,
				Timestamp: time.Date(1999, time.March, 6, 7, 11, 0, 0, time.UTC),
				Sent: QSOInfo{
					Call:     callsign.MustParse("HC8N"),
					RST:      "59",
					Exchange: []string{"700"},
				},
				Received: QSOInfo{
					Call:     callsign.MustParse("W1AW"),
					RST:      "59",
					Exchange: []string{"CT"},
				},
				Transmitter: 1,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			actual, err := ParseQSO(tc.value)
			if tc.invalid {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestReadTestdata(t *testing.T) {
	entries, err := os.ReadDir("testdata")
	require.NoError(t, err)
	for _, entry := range entries {
		t.Run(entry.Name(), func(t *testing.T) {
			file, err := os.Open("testdata/" + entry.Name())
			require.NoError(t, err)
			defer file.Close()

			_, err = Read(file)
			assert.NoError(t, err)
		})
	}
}
