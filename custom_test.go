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

func TestSACArea(t *testing.T) {
	tt := []struct {
		call     string
		entity   DXCCEntity
		expected string
	}{
		{"SM3ABC", "sm", "sm3"},
		{"SI3ABC", "sm", "sm3"},
		{"SK3AB", "sm", "sm3"},
		{"7S3ABC", "sm", "sm3"},
		{"8S3ABC", "sm", "sm3"},
		{"LA/G3XYZ", "la", "la0"},
		{"OZ150A", "oz", "oz1"},
		{"OZ1ABC", "oz", "oz1"},
		{"SM5/SM3ABC", "sm", "sm5"},
		{"OH0XX", "oh0", "oh00"},
		{"OJ0AM", "oj0", "oj00"},
		{"JW5E", "jw", "jw5"},
		{"DL1ABC", "", ""},
	}
	for _, tc := range tt {
		t.Run(tc.call, func(t *testing.T) {
			assert.Equal(t, tc.expected, sacArea(callsign.MustParse(tc.call), tc.entity))
		})
	}
}

func TestSARLArea(t *testing.T) {
	tt := []struct {
		call     string
		entity   DXCCEntity
		expected string
	}{
		{"ZS1ABC", "zs", "1"},
		{"ZR1ABC", "zs", "1"},
		{"ZU1ABC", "zs", "1"},
		{"ZS2ABC", "zs", "2"},
		{"ZS6ABC", "zs", "6"},
		{"ZS0ABC", "zs", "9"},
		{"ZS7ABC", "ce9", "8"},
		{"ZS8XX", "zs8", "8"},
		{"V51ABC", "v5", "7"},
		{"3DA0AB", "3da", "8"},
		{"FR5ABC", "fr", "8"},
		{"D68ABC", "d6", "8"},
		{"DL1ABC", "dl", "9"},
		{"K1ABC", "k", "9"},
	}
	for _, tc := range tt {
		t.Run(tc.call, func(t *testing.T) {
			assert.Equal(t, tc.expected, sarlArea(callsign.MustParse(tc.call), tc.entity))
		})
	}
}

func TestBalkanPrefix(t *testing.T) {
	tt := []struct {
		call     string
		expected string
	}{
		{"LZ6Y", "LZ6"},
		{"LZ07KM", "LZ0"},
		{"YO2014A", "YO2"},
		{"ER650MD", "ER6"},
		{"SV0XCA/5", "SV5"},
		{"YO8WW/QRP", "YO8"},
		{"4O4A", "4O4"},
		{"9A5MP", "9A5"},
		{"ZC4A", "ZC4"},
		{"TA4RC", "TA4"},
		{"DL1ABC", ""},
		{"K1ABC", ""},
	}
	for _, tc := range tt {
		t.Run(tc.call, func(t *testing.T) {
			assert.Equal(t, tc.expected, balkanPrefix(callsign.MustParse(tc.call)))
		})
	}
}

func TestCommonwealthArea(t *testing.T) {
	tt := []struct {
		call     string
		entity   DXCCEntity
		expected string
	}{
		{"G3ABC", "g", "G"},
		{"GM4ABC", "gm", "GM"},
		{"VE3ABC", "ve", "VE3"},
		{"VO1ABC", "ve", "VO1"},
		{"VY1ABC", "ve", "VY1"},
		{"VK3ABC", "vk", "VK3"},
		{"VK5ABC", "vk", "VK5"},
		{"ZL1ABC", "zl", "ZL1"},
		{"ZS1ABC", "zs", "ZS1"},
		{"ZS8XX", "zs8", "ZS8"},
		{"9V1AA", "9v", "9V"},
		{"DL1ABC", "dl", ""},
		{"K1ABC", "k", ""},
		{"UA3ABC", "ua", ""},
	}
	for _, tc := range tt {
		t.Run(tc.call, func(t *testing.T) {
			assert.Equal(t, tc.expected, commonwealthArea(callsign.MustParse(tc.call), tc.entity))
		})
	}
}
