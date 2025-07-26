package conval

import "time"

const DefaultBreakDuration = 30 * time.Minute

type TimeReport struct {
	MinBreakDuration time.Duration
	// ActiveMinutes counts minutes with logged QSOs
	ActiveMinutes int
	// IdleMinutes counts minutes between active minutes without QSO that do not qualify as break
	IdleMinutes int
	// BreakMinutes counts consecutive minutes that qualify as break
	BreakMinutes int
	// Breaks counts the number of breaks
	Breaks int
	// TotalMinutes counts from the first to the last active minute
	TotalMinutes int
}

func (r TimeReport) OperationTime() time.Duration {
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

func (s *TimeSheet) MarkActive(now time.Time) {
	if now.Before(s.startTime) || now.After(s.endTime) {
		return
	}
	minuteIndex := s.getMinuteIndex(now)
	if minuteIndex >= 0 && minuteIndex < len(s.activeMinutes) {
		s.activeMinutes[minuteIndex] = true
	}
}

func (s *TimeSheet) getMinuteIndex(now time.Time) int {
	if now.Before(s.startTime) {
		return -1
	}
	minutes := int(now.Sub(s.startTime).Minutes())
	if minutes < 0 {
		return -1
	}
	return minutes
}

func (s *TimeSheet) TimeReport(minBreakDuration time.Duration) TimeReport {
	result := TimeReport{
		MinBreakDuration: minBreakDuration,
	}

	minBreakMinutes := int(minBreakDuration.Minutes())
	firstActive := -1
	lastActive := -1
	breakStart := -1
	for i, active := range s.activeMinutes {
		if active {
			result.ActiveMinutes++
			if firstActive == -1 {
				firstActive = i
			}
			lastActive = i
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
		breakMinutes := len(s.activeMinutes) - breakStart
		if breakMinutes >= minBreakMinutes {
			result.BreakMinutes += breakMinutes
			result.Breaks++
		} else {
			result.IdleMinutes += breakMinutes
		}
	}

	result.TotalMinutes = lastActive - firstActive + 1

	return result
}

func (c *Counter) ComputeMinBreakDuration() time.Duration {
	if len(c.definition.Breaks) == 1 {
		return c.definition.Breaks[0].Duration
	}

	for _, b := range c.definition.Breaks {
		if (b.Constraint.OperatorMode == c.setup.OperatorMode) &&
			(b.Constraint.Overlay == c.setup.Overlay) {
			return b.Duration
		}
	}

	return DefaultBreakDuration
}
