package conval

import (
	"log"
	"time"
)

type ScoreBin struct {
	QSOs        int                          `yaml:"qsos" json:"qsos"`
	Points      int                          `yaml:"points" json:"points"`
	Multis      int                          `yaml:"multis" json:"multis"`
	MultiValues map[Property]map[string]bool `yaml:"-" json:"-"`
	Bands       map[ContestBand]bool         `yaml:"-" json:"-"`
}

func (b ScoreBin) Total() int {
	return b.Points * b.Multis
}

func (b *ScoreBin) Add(qso ScoredQSO) {
	b.QSOs++
	b.Points += qso.Points
	b.Multis += qso.Multis
	if b.MultiValues == nil {
		b.MultiValues = make(map[Property]map[string]bool)
	}
	for property, multi := range qso.MultiValues {
		multisPerProperty := b.MultiValues[property]
		if multisPerProperty == nil {
			multisPerProperty = make(map[string]bool)
		}
		multisPerProperty[multi] = true
		b.MultiValues[property] = multisPerProperty
	}
	if b.Bands == nil {
		b.Bands = make(map[ContestBand]bool)
	}
	b.Bands[qso.Band] = true
}

func (c Counter) EvaluateAll(startTime time.Time, resolution time.Duration) []ScoreBin {
	if resolution == 0 {
		resolution = time.Hour
	}
	binCount := c.definition.Duration / resolution

	result := make([]ScoreBin, binCount)

	for _, qso := range c.qsos {
		binIndex := toBinIndex(qso.Timestamp, startTime, resolution)
		if binIndex < 0 || binIndex >= len(result) {
			log.Printf("invalid qso timestamp %s", qso.Timestamp)
			continue
		}
		result[binIndex].Add(qso)
	}

	return result
}

func toBinIndex(t time.Time, startTime time.Time, resolution time.Duration) int {
	if resolution == 0 {
		resolution = time.Hour
	}
	if t.Before(startTime) {
		return -1
	}
	binTime := t.Truncate(resolution)
	binOffset := binTime.Sub(startTime)

	return int(binOffset / resolution)
}
