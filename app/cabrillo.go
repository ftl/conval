package app

import (
	"os"
	"strings"

	"github.com/ftl/conval"
	"github.com/ftl/conval/cabrillo"
)

func ReadCabrilloLogFromFile(filename string, prefixes conval.PrefixDatabase) (*CabrilloLogfile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	log, err := cabrillo.Read(file)
	if err != nil {
		return nil, err
	}

	return &CabrilloLogfile{log, prefixes}, nil
}

type CabrilloLogfile struct {
	log      *cabrillo.Log
	prefixes conval.PrefixDatabase
}

func (l CabrilloLogfile) Identifier() conval.ContestIdentifier {
	return conval.ContestIdentifier(l.log.Contest)
}

func (l CabrilloLogfile) Setup() *conval.Setup {
	result := new(conval.Setup)
	result.MyCall = l.log.Callsign
	myContinent, myCountry, found := l.prefixes.Find(result.MyCall.String())
	if found {
		result.MyContinent = myContinent
		result.MyCountry = myCountry
	}

	result.GridLocator = l.log.GridLocator
	result.Operators = l.log.Operators

	result.OperatorMode = cabrilloToOperatorMode(l.log.Category.Operator)
	result.Overlay = cabrilloToOverlay(l.log.Category.Overlay)
	result.Power = conval.PowerMode(strings.ToLower(string(l.log.Category.Power)))
	result.Bands = []conval.ContestBand{conval.ContestBand(strings.ToLower(string(l.log.Category.Band)))}
	result.Modes = cabrilloToModes(l.log.Category.Mode)

	return result
}

func (l CabrilloLogfile) QSOs(exchangeFields func(conval.Continent, conval.DXCCEntity) []conval.ExchangeField) []conval.QSO {
	result := make([]conval.QSO, len(l.log.QSOData))
	for i, qso := range l.log.QSOData {
		resultQSO := conval.QSO{
			TheirCall: qso.Received.Call,
			Timestamp: qso.Timestamp,
			Band:      cabrilloToBand(qso.Frequency),
			Mode:      cabrilloToQSOMode(qso.Mode),
		}
		theirContinent, theirCountry, found := l.prefixes.Find(resultQSO.TheirCall.String())
		if found {
			resultQSO.TheirContinent = theirContinent
			resultQSO.TheirCountry = theirCountry
		}
		fields := exchangeFields(resultQSO.TheirContinent, resultQSO.TheirCountry)
		resultQSO.TheirExchange = cabrilloToQSOExchange(fields, qso.Received)

		result[i] = resultQSO
	}
	return result
}

func cabrilloToOperatorMode(operator cabrillo.CategoryOperator) conval.OperatorMode {
	if operator == cabrillo.MultiOperator {
		return conval.MultiOperator
	}
	return conval.SingleOperator
}

func cabrilloToOverlay(overlay cabrillo.CategoryOverlay) conval.Overlay {
	switch overlay {
	case cabrillo.ClassicOverlay:
		return conval.ClassicOverlay
	case cabrillo.TBWiresOverlay:
		return conval.ThreeBandAndWiresOverlay
	case cabrillo.RookieOverlay:
		return conval.RookieOverlay
	case cabrillo.YouthOverlay:
		return conval.YouthOverlay
	default:
		return conval.NoOverlay
	}
}

func cabrilloToModes(mode cabrillo.CategoryMode) []conval.Mode {
	switch mode {
	case cabrillo.ModeMIXED:
		return []conval.Mode{conval.ModeALL}
	case cabrillo.ModeDIGI:
		return []conval.Mode{conval.ModeDigital}
	default:
		return []conval.Mode{conval.Mode(strings.ToLower(string(mode)))}
	}
}

func cabrilloToBand(frequency cabrillo.QSOFrequency) conval.ContestBand {
	band := frequency.ToBand()
	return conval.ContestBand(strings.ToLower(string(band)))
}

func cabrilloToQSOMode(mode cabrillo.QSOMode) conval.Mode {
	switch mode {
	case cabrillo.QSOModePhone:
		return conval.ModeSSB
	case cabrillo.QSOModeDigi:
		return conval.ModeDigital
	case cabrillo.QSOModeRTTY:
		return conval.ModeRTTY
	default:
		return conval.Mode(strings.ToLower(string(mode)))
	}
}

func cabrilloToQSOExchange(fields []conval.ExchangeField, info cabrillo.QSOInfo) conval.QSOExchange {
	values := make([]string, 0, len(info.Exchange)+1)
	values = append(values, info.RST)
	values = append(values, info.Exchange...)
	return conval.ParseExchange(fields, values)
}
