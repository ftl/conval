package conval

import (
	"testing"

	"github.com/ftl/hamradio/dxcc"
	"github.com/stretchr/testify/assert"
)

func TestParseExchange(t *testing.T) {
	tt := []struct {
		desc     string
		fields   []ExchangeField
		values   []string
		expected QSOExchange
	}{
		{
			desc:     "only rst",
			fields:   []ExchangeField{[]Property{RSTProperty}},
			values:   []string{"123"},
			expected: QSOExchange{RSTProperty: "123"},
		},
		{
			desc:     "rst and member number",
			fields:   []ExchangeField{[]Property{RSTProperty}, []Property{MemberNumberProperty, NoMemberProperty}},
			values:   []string{"59", "123"},
			expected: QSOExchange{RSTProperty: "59", MemberNumberProperty: "123"},
		},
		{
			desc:     "rst and no member",
			fields:   []ExchangeField{[]Property{RSTProperty}, []Property{MemberNumberProperty, NoMemberProperty}},
			values:   []string{"59", "nm"},
			expected: QSOExchange{RSTProperty: "59", NoMemberProperty: "NM"},
		},
		{
			desc:     "rst and empty no member",
			fields:   []ExchangeField{[]Property{RSTProperty}, []Property{MemberNumberProperty, NoMemberProperty}},
			values:   []string{"59", ""},
			expected: QSOExchange{RSTProperty: "59"},
		},
		{
			desc:     "rst and serial",
			fields:   []ExchangeField{[]Property{RSTProperty}, []Property{SerialNumberProperty}},
			values:   []string{"59", "123"},
			expected: QSOExchange{RSTProperty: "59", SerialNumberProperty: "123"},
		},
		{
			desc:     "rst and empty serial",
			fields:   []ExchangeField{[]Property{RSTProperty}, []Property{SerialNumberProperty}},
			values:   []string{"59", ""},
			expected: QSOExchange{RSTProperty: "59"},
		},
		{
			desc:     "rst and no serial",
			fields:   []ExchangeField{[]Property{RSTProperty}, []Property{SerialNumberProperty}},
			values:   []string{"59"},
			expected: QSOExchange{RSTProperty: "59"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			actual := ParseExchange(tc.fields, tc.values, nil, PropertyValidatorsFunc(CommonPropertyValidator))
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPrefixDatabase_Find(t *testing.T) {
	prefixes, _, err := dxcc.DefaultPrefixes(true)
	assert.NoError(t, err)

	tt := []struct {
		call         string
		compliant    DXCCEntity
		notCompliant DXCCEntity
	}{
		{
			call:         "4U1VIC",
			compliant:    "oe",
			notCompliant: "4u1v",
		},
		{
			call:         "IT9ABC",
			compliant:    "i",
			notCompliant: "it9",
		},
		{
			call:         "IP9P",
			compliant:    "i",
			notCompliant: "ig9",
		},
		{
			call:         "JW/LB2PG",
			compliant:    "jw",
			notCompliant: "jw/b",
		},
	}
	for _, tc := range tt {
		t.Run(tc.call, func(t *testing.T) {
			db := &prefixDatabase{prefixes: prefixes, arrlCompliant: true}
			_, entity, _, _, found := db.Find(tc.call)
			assert.True(t, found)
			assert.Equal(t, tc.compliant, entity, "compliant")
			db.arrlCompliant = false
			_, entity, _, _, found = db.Find(tc.call)
			assert.True(t, found)
			assert.Equal(t, tc.notCompliant, entity, "not compliant")
		})
	}
}
