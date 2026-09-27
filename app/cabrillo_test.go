package app

import (
	"testing"

	"github.com/ftl/cabrillo"
	"github.com/stretchr/testify/assert"

	"github.com/ftl/conval"
)

func TestCabrilloToOverlay(t *testing.T) {
	tt := []struct {
		value    cabrillo.CategoryOverlay
		expected conval.Overlay
	}{
		{"", conval.NoOverlay},
		{cabrillo.ClassicOverlay, conval.ClassicOverlay},
		{cabrillo.TBWiresOverlay, conval.ThreeBandAndWiresOverlay},
		{cabrillo.WireOnlyOverlay, conval.WireOnlyOverlay},
		{cabrillo.RookieOverlay, conval.RookieOverlay},
		{cabrillo.YouthOverlay, conval.YouthOverlay},
		{cabrillo.YLOverlay, conval.YLOverlay},
		{cabrillo.YNOverlay, conval.YNOverlay},
		{cabrillo.TeenOverlay, conval.TeenOverlay},
		{cabrillo.NewcomerOverlay, conval.NewcomerOverlay},
		{cabrillo.DXpeditionOverlay, conval.DXpeditionOverlay},
		{cabrillo.SingleElementOverlay, conval.SingleElementOverlay},
		{cabrillo.TwelveHourOverlay, conval.TwelveHourOverlay},
		{cabrillo.NoviceTechOverlay, conval.NoviceTechOverlay},
		{cabrillo.Over50Overlay, conval.Over50Overlay},
	}
	for _, tc := range tt {
		t.Run(string(tc.value), func(t *testing.T) {
			assert.Equal(t, tc.expected, cabrilloToOverlay(tc.value))
		})
	}
}
