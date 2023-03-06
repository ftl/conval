package conval

import (
	"testing"

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
