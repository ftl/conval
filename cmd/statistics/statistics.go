package statistics

import (
	"fmt"
	"math"
	"time"

	"github.com/ftl/conval"
	"github.com/ftl/conval/app"
)

type Rate struct {
	TotalAverage  float64       `yaml:"average_total" json:"average_total"`
	ActiveAverage float64       `yaml:"average_active" json:"average_active"`
	Max           float64       `yaml:"max" json:"max"`
	Min           float64       `yaml:"min" json:"min"`
	MaxHour       time.Duration `yaml:"best_hour" json:"best_hour"`
	MinHour       time.Duration `yaml:"slowest_hour" json:"slowest_hour"`
}

func (r *Rate) UpdateMinMax(value int, offset time.Duration) {
	if float64(value) < r.Min {
		r.Min = float64(value)
		r.MinHour = offset
	}
	if float64(value) > r.Max {
		r.Max = float64(value)
		r.MaxHour = offset
	}
}

type Result struct {
	ContestName string        `yaml:"contest_name" json:"contest_name"`
	ContestID   string        `yaml:"contest_id" json:"contest_id"`
	StartTime   time.Time     `yaml:"start_time" json:"start_time"`
	Duration    time.Duration `yaml:"duration" json:"duration"`
	ActiveHours int           `yaml:"active_hours" json:"active_hours"`
	QSOs        int           `yaml:"qsos" json:"qsos"`
	Points      int           `yaml:"points" json:"points"`
	Multis      int           `yaml:"multis" json:"multis"`
	Total       int           `yaml:"total" json:"total"`
	QSORate     Rate          `yaml:"qso_rate" json:"qso_rate"`
	PointsRate  Rate          `yaml:"points_rate" json:"points_rate"`
	MultiRate   Rate          `yaml:"multi_rate" json:"multi_rate"`
}

func Evaluate(logfile app.Logfile, definition *conval.Definition, setup *conval.Setup, startTime time.Time) (Result, error) {
	var err error

	definitionForFile := definition
	if definitionForFile == nil {
		definitionForFile, err = conval.IncludedDefinition(string(logfile.Identifier()))
		if err != nil {
			return Result{}, err
		}
	}
	if definitionForFile == nil {
		return Result{}, fmt.Errorf("no contest definition found")
	}

	setupForFile := setup
	if setupForFile == nil {
		setupForFile = logfile.Setup()
	}
	if setupForFile == nil {
		return Result{}, fmt.Errorf("no setup defined")
	}

	prefixes, err := conval.NewPrefixDatabase()
	if err != nil {
		return Result{}, err
	}

	counter := conval.NewCounter(*definitionForFile, *setupForFile, prefixes)
	qsos := logfile.QSOs(definition, counter.EffectiveExchangeFields)
	for _, qso := range qsos {
		counter.Add(qso)
	}

	scoreBins := counter.EvaluateAll(startTime, time.Hour)
	totalHours := float64(len(scoreBins))
	result := Result{
		ContestName: definitionForFile.Name,
		ContestID:   string(definitionForFile.Identifier),
		StartTime:   startTime,
		Duration:    definitionForFile.Duration,
		QSORate:     Rate{Min: math.MaxInt},
		PointsRate:  Rate{Min: math.MaxInt},
		MultiRate:   Rate{Min: math.MaxInt},
	}

	result.ActiveHours = 0
	for i, bin := range scoreBins {
		offset := time.Duration(i+1) * time.Hour
		if bin.QSOs == 0 {
			continue
		}
		result.ActiveHours++

		result.QSOs += bin.QSOs
		result.Points += bin.Points
		result.Multis += bin.Multis
		result.QSORate.UpdateMinMax(bin.QSOs, offset)
		result.PointsRate.UpdateMinMax(bin.Points, offset)
		result.MultiRate.UpdateMinMax(bin.Multis, offset)
	}
	result.Total = result.Points * result.Multis
	result.QSORate.TotalAverage = float64(result.QSOs) / totalHours
	result.QSORate.ActiveAverage = float64(result.QSOs) / float64(result.ActiveHours)
	result.PointsRate.TotalAverage = float64(result.Points) / totalHours
	result.PointsRate.ActiveAverage = float64(result.Points) / float64(result.ActiveHours)
	result.MultiRate.TotalAverage = float64(result.Multis) / totalHours
	result.MultiRate.ActiveAverage = float64(result.Multis) / float64(result.ActiveHours)

	return result, nil
}
