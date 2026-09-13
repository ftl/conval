package conval

import (
	"testing"

	"github.com/ftl/hamradio/callsign"
	"github.com/stretchr/testify/assert"
)

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
