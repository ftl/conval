package score

import (
	"fmt"

	"github.com/ftl/conval"
)

type Logfile interface {
	Identifier() conval.ContestIdentifier
	Setup() *conval.Setup
	QSOs() []conval.QSO
}

type Score struct {
	Points int
	Multis int
}

func (s Score) Total() int {
	return s.Points * s.Multis
}

func Evaluate(logfile Logfile, setup conval.Setup, definition conval.Definition) (Score, error) {
	return Score{}, fmt.Errorf("not yet implemented")
}
