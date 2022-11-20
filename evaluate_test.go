package conval

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCounter_EvaluateAll_ResultSize(t *testing.T) {
	definition := Definition{
		Duration: 12 * time.Hour,
	}
	setup := Setup{}
	counter := NewCounter(definition, setup)
	startTime := time.Now().Truncate(time.Hour)

	perHour := counter.EvaluateAll(startTime, time.Hour)
	assert.Equal(t, 12, len(perHour))
	per5Minute := counter.EvaluateAll(startTime, 5*time.Minute)
	assert.Equal(t, 144, len(per5Minute))
}

func TestToBinIndex(t *testing.T) {
	startTime := time.Now().Truncate(time.Hour)
	tt := []struct {
		desc       string
		offset     time.Duration
		resolution time.Duration
		expected   int
	}{
		{
			desc:       "zero",
			offset:     0,
			resolution: 0,
			expected:   0,
		},
		{
			desc:       "first bin ends with resolution",
			offset:     time.Hour - time.Second,
			resolution: time.Hour,
			expected:   0,
		},
		{
			desc:       "second bin contains resolution",
			offset:     time.Hour,
			resolution: time.Hour,
			expected:   1,
		},
		{
			desc:       "the past is -1",
			offset:     -time.Second,
			resolution: time.Hour,
			expected:   -1,
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			actual := toBinIndex(startTime.Add(tc.offset), startTime, tc.resolution)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
