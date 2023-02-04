package conval

import (
	"testing"

	"github.com/ftl/hamradio/callsign"
	"github.com/stretchr/testify/assert"
)

func TestVeronEntity(t *testing.T) {
	tt := []struct {
		call       string
		dxccEntity DXCCEntity
		expected   string
	}{
		{"JA1ABC", "JA", "JA1"},
		{"LU/G3XYZ", "LU", "LU0"},
		{"W8/G3XYZ", "K", "K8"},
		{"K5ZD", "K", "K5"},
		{"K5ZD/1", "K", "K1"},
		{"XK2ABC", "VE", "VE2"},
	}
	for _, tc := range tt {
		t.Run(tc.call, func(t *testing.T) {
			actual := VeronEntity(callsign.MustParse(tc.call), tc.dxccEntity)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestCanadianTerritory(t *testing.T) {
	tt := []struct {
		call     string
		expected string
	}{
		{"VE2ABC", "VE2"},
		{"VA2ABC", "VE2"},
		{"CG2ABC", "VE2"},
		{"XK2ABC", "VE2"},
		{"VO1ABC", "VO1"},
		{"VY1ABC", "VY1"},
	}
	for _, tc := range tt {
		t.Run(tc.call, func(t *testing.T) {
			actual := CanadianTerritory(callsign.MustParse(tc.call))

			assert.Equal(t, tc.expected, actual)
		})
	}
}
