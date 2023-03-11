package conval

import (
	"testing"

	"github.com/ftl/hamradio/callsign"
	"github.com/stretchr/testify/assert"
)

func TestValidateNoMember(t *testing.T) {
	validator := commonPropertyValidators[NoMemberProperty]
	assert.NoError(t, validator.ValidateProperty("nm", nil), "nm")
	assert.NoError(t, validator.ValidateProperty("NM", nil), "NM")
	assert.Error(t, validator.ValidateProperty("", nil), "empty")
	assert.Error(t, validator.ValidateProperty("   ", nil), "whitespace")
	assert.Error(t, validator.ValidateProperty("12345", nil), "number")
}

func TestWPXPrefix(t *testing.T) {
	tt := []struct {
		call     string
		expected string
	}{
		{"DL1ABC", "DL1"},
		{"9A1A", "9A1"},
		{"LY1000A", "LY1000"},
		{"DL/9A1A", "DL0"},
		{"N8BJQ/KH9", "KH9"},
		{"N8BJQ/9", "N8"},
		{"DL1ABC/P", "DL1"},
		{"ON2ABC", "ON2"},
		{"OT0ABC", "OT0"},
	}
	for _, tc := range tt {
		t.Run(tc.call, func(t *testing.T) {
			actual := WPXPrefix(callsign.MustParse(tc.call))
			assert.Equal(t, tc.expected, actual)
		})
	}
}
