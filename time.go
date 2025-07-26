package conval

import "time"

type TimeRecord struct {
	MinBreakDuration time.Duration
	ActiveMinutes    int
	IdleMinutes      int
	BreakMinutes     int
	Breaks           int
}

func (r TimeRecord) OperationTime() time.Duration {
	return time.Duration(r.ActiveMinutes+r.IdleMinutes) * time.Minute
}

type TimeSheet struct {
	startTime     time.Time
	endTime       time.Time
	duration      time.Duration
	activeMinutes []bool
}

func NewTimeSheet(startTime time.Time, duration time.Duration) *TimeSheet {
	minutes := int(duration.Minutes())

	return &TimeSheet{
		startTime:     startTime,
		endTime:       startTime.Add(duration),
		duration:      duration,
		activeMinutes: make([]bool, minutes),
	}
}

func (ts *TimeSheet) MarkActive(now time.Time) {
	if now.Before(ts.startTime) || now.After(ts.endTime) {
		return
	}
	minuteIndex := ts.getMinuteIndex(now)
	if minuteIndex >= 0 && minuteIndex < len(ts.activeMinutes) {
		ts.activeMinutes[minuteIndex] = true
	}
}

func (ts *TimeSheet) getMinuteIndex(now time.Time) int {
	if now.Before(ts.startTime) {
		return -1
	}
	minutes := int(now.Sub(ts.startTime).Minutes())
	if minutes < 0 {
		return -1
	}
	return minutes
}

func (ts *TimeSheet) TimeRecord(minBreakDuration time.Duration) TimeRecord {
	result := TimeRecord{
		MinBreakDuration: minBreakDuration,
	}

	minBreakMinutes := int(minBreakDuration.Minutes())
	breakStart := -1
	for i, active := range ts.activeMinutes {
		if active {
			result.ActiveMinutes++
		}

		switch {
		case !active && breakStart == -1:
			breakStart = i
		case active && breakStart != -1:
			breakMinutes := i - breakStart
			if breakMinutes >= minBreakMinutes {
				result.BreakMinutes += breakMinutes
				result.Breaks++
			} else {
				result.IdleMinutes += breakMinutes
			}
			breakStart = -1
		}
	}

	if breakStart != -1 {
		breakMinutes := len(ts.activeMinutes) - breakStart
		if breakMinutes >= minBreakMinutes {
			result.BreakMinutes += breakMinutes
			result.Breaks++
		} else {
			result.IdleMinutes += breakMinutes
		}
	}

	return result
}
