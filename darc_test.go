package conval

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/ftl/hamradio/callsign"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWAECallAreaNumber(t *testing.T) {
	tt := []struct {
		call       string
		dxccEntity DXCCEntity
		expected   string
	}{
		{"W1ABC", "K", "1"},
		{"KA1ABC", "K", "1"},
		{"K3ABC/1", "K", "1"},
		{"K/DL3ABC/1", "K", "1"},
		{"VO1ABC", "VE", "1"},
		{"VJ2ABC", "VK", "2"},
		{"VK3ABC", "ZL", "3"},
		{"V93ABC", "ZS", "3"},
		{"7M4ABC", "JA", "4"},
		{"3M5ABC", "BY", "5"},
		{"PP6ABC", "PY", "6"},
		{"RA8ABC", "UA9", "8"},
		{"RA9ABC", "UA9", "9"},
		{"RA0ABC", "UA9", "0"},
	}
	for _, tc := range tt {
		t.Run(tc.call, func(t *testing.T) {
			actual := waeCallAreaNumber(callsign.MustParse(tc.call), tc.dxccEntity)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestValidateSDOKs(t *testing.T) {
	validateWAGDOK := PropertyValidators[WAGDOKProperty]

	sdokFile, err := os.Open("testdata/sdok.txt")
	require.NoError(t, err)
	defer sdokFile.Close()

	scanner := bufio.NewScanner(sdokFile)
	for scanner.Scan() {
		sdok := strings.ToLower(scanner.Text())
		t.Run(sdok, func(t *testing.T) {
			assert.NoError(t, validateWAGDOK.ValidateProperty(sdok))
		})
	}
}
