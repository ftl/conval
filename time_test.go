package conval

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeSheet_TimeRecord(t *testing.T) {
	startTime := time.Now().Add(-2 * time.Hour)
	duration := 1 * time.Hour
	minBreakDuration := 10 * time.Minute

	timeSheet := NewTimeSheet(startTime, duration)
	timeSheet.MarkActive(startTime.Add(12 * time.Minute))
	timeSheet.MarkActive(startTime.Add(13 * time.Minute))
	timeSheet.MarkActive(startTime.Add(15 * time.Minute))
	timeSheet.MarkActive(startTime.Add(20 * time.Minute))
	timeSheet.MarkActive(startTime.Add(33 * time.Minute))
	timeSheet.MarkActive(startTime.Add(35 * time.Minute))
	timeRecord := timeSheet.TimeRecord(minBreakDuration)

	assert.Equal(t, 6, timeRecord.ActiveMinutes)
	assert.Equal(t, 6, timeRecord.IdleMinutes)
	assert.Equal(t, 48, timeRecord.BreakMinutes)
	assert.Equal(t, 3, timeRecord.Breaks)
	assert.Equal(t, 24, timeRecord.TotalMinutes)
	assert.Equal(t, 12*time.Minute, timeRecord.OperationTime())
}

func TestCounter_computeMinBreakDuration(t *testing.T) {
	tests := []struct {
		cabrilloName string
		expected     time.Duration
	}{
		{"agcw-hny", DefaultBreakDuration},
		{"cq-160-cw", 30 * time.Minute},
		{"cq-wpx-cw", 1 * time.Hour},
		{"cq-ww-cw", 1 * time.Hour},
		{"darc-waedc-cw", 1 * time.Hour},
	}
	for _, test := range tests {
		t.Run(test.cabrilloName, func(t *testing.T) {
			defintion, err := IncludedDefinition(test.cabrilloName)
			require.NoError(t, err)
			counter := NewCounter(*defintion, Setup{}, nil)

			actual := counter.ComputeMinBreakDuration()
			assert.Equal(t, test.expected, actual)
		})
	}
}
