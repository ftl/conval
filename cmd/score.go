package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ftl/hamradio/dxcc"
	"github.com/spf13/cobra"

	"github.com/ftl/conval"
	"github.com/ftl/conval/cabrillo"
)

var scoreFlags = struct {
	setupFilename      string
	definitionFilename string
	cabrilloName       string
}{}

var scoreCmd = &cobra.Command{
	Use:   "score <filename>",
	Short: "calculate the score of a contest log given as Cabrillo or ADIF file",
	Run:   runScore,
}

func init() {
	scoreCmd.Flags().StringVar(&scoreFlags.setupFilename, "setup", "", "the setup file")
	scoreCmd.Flags().StringVar(&scoreFlags.definitionFilename, "definition", "", "the contest definition file")
	scoreCmd.Flags().StringVar(&scoreFlags.cabrilloName, "cabrillo", "", "the cabrillo name (see https://www.contestcalendar.com/cabnames.php)")

	rootCmd.AddCommand(scoreCmd)
}

func runScore(cmd *cobra.Command, args []string) {
	var err error

	prefixes, err := newPrefixDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the CLI information

	var definition *conval.Definition
	if scoreFlags.definitionFilename != "" {
		definition, err = loadDefinitionFromFile(scoreFlags.definitionFilename)
	} else if scoreFlags.cabrilloName != "" {
		definition, err = conval.IncludedDefinition(scoreFlags.cabrilloName)
	}
	if err != nil {
		log.Fatal(err)
	}

	var setup *conval.Setup
	if scoreFlags.setupFilename != "" {
		setup, err = loadSetupFromFile(scoreFlags.setupFilename, prefixes)
		if err != nil {
			log.Fatal(err)
		}
	}
	// TODO override the setup information with values given with explicit CLI flags

	if len(args) < 1 {
		log.Fatal("missing input filename")
	}
	for _, filename := range args {
		log.Printf("evaluating %s", filename)

		// TODO: detect the input file format
		logfile, err := readCabrilloLogFromFile(filename, prefixes)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("claimed score: %d", logfile.log.ClaimedScore)

		definitionForFile := definition
		if definitionForFile == nil {
			definitionForFile, err = conval.IncludedDefinition(string(logfile.Identifier()))
			if err != nil {
				log.Fatal(err)
			}
		}
		if definitionForFile == nil {
			log.Fatal("no contest definition found")
		} else {
			log.Printf("%s: %s", definitionForFile.Identifier, definitionForFile.Name)
		}

		setupForFile := setup
		if setupForFile == nil {
			setupForFile = logfile.Setup()
		}
		if setupForFile == nil {
			log.Fatal("no setup defined")
		} else {
			log.Printf("setup: %+v", setupForFile)
		}

		counter := conval.NewCounter(*setupForFile, definitionForFile.Exchange, definitionForFile.Scoring)
		qsos := logfile.QSOs(counter.EffectiveExchangeFields)
		for _, qso := range qsos {
			counter.Add(qso)
		}
		totalScore := counter.TotalScore()
		log.Printf("QSOs: %d Score: %d * %d = %d", len(qsos), totalScore.Multis, totalScore.Points, totalScore.Multis*totalScore.Points)
	}
}

func loadDefinitionFromFile(filename string) (*conval.Definition, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return conval.LoadDefinitionYAML(file)
}

func loadSetupFromFile(filename string, prefixes prefixDatabase) (*conval.Setup, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	result, err := conval.LoadSetupYAML(file)
	if err != nil {
		return nil, err
	}

	result.MyContinent, result.MyCountry = toConvalContinentAndCountry(prefixes.Find(result.MyCall.String()))

	return result, nil
}

type prefixDatabase interface {
	Find(s string) ([]dxcc.Prefix, bool)
}

func toConvalContinentAndCountry(entities []dxcc.Prefix, found bool) (conval.Continent, conval.DXCCEntity) {
	if !found || len(entities) != 1 {
		return "", ""
	}
	return conval.Continent(strings.ToLower(entities[0].Continent)), conval.DXCCEntity(strings.ToLower(entities[0].PrimaryPrefix))
}

func newPrefixDatabase() (prefixDatabase, error) {
	localFilename, err := dxcc.LocalFilename()
	if err != nil {
		return nil, err
	}
	updated, err := dxcc.Update(dxcc.DefaultURL, localFilename)
	if err != nil {
		fmt.Printf("update of local copy failed: %v\n", err)
	}
	if updated {
		fmt.Printf("updated local copy: %v\n", localFilename)
	}

	return dxcc.LoadLocal(localFilename)
}

func readCabrilloLogFromFile(filename string, prefixes prefixDatabase) (*cabrilloLogfile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	log, err := cabrillo.Read(file)
	if err != nil {
		return nil, err
	}

	return &cabrilloLogfile{log, prefixes}, nil
}

type cabrilloLogfile struct {
	log      *cabrillo.Log
	prefixes prefixDatabase
}

func (l cabrilloLogfile) Identifier() conval.ContestIdentifier {
	return conval.ContestIdentifier(l.log.Contest)
}

func (l cabrilloLogfile) Setup() *conval.Setup {
	result := new(conval.Setup)
	result.MyCall = l.log.Callsign
	result.MyContinent, result.MyCountry = toConvalContinentAndCountry(l.prefixes.Find(result.MyCall.String()))

	result.GridLocator = l.log.GridLocator
	result.Operators = l.log.Operators

	result.OperatorMode = toOperatorMode(l.log.Category.Operator)
	result.Overlay = toOverlay(l.log.Category.Overlay)
	result.Power = conval.PowerMode(strings.ToLower(string(l.log.Category.Power)))
	result.Bands = []conval.ContestBand{conval.ContestBand(strings.ToLower(string(l.log.Category.Band)))}
	result.Modes = toModes(l.log.Category.Mode)

	return result
}

func (l cabrilloLogfile) QSOs(exchangeFields func(conval.Continent, conval.DXCCEntity) []conval.ExchangeField) []conval.QSO {
	result := make([]conval.QSO, len(l.log.QSOData))
	for i, qso := range l.log.QSOData {
		resultQSO := conval.QSO{
			TheirCall: qso.Received.Call,
			Timestamp: qso.Timestamp,
			Band:      toBand(qso.Frequency),
			Mode:      toQSOMode(qso.Mode),
		}
		resultQSO.TheirContinent, resultQSO.TheirCountry = toConvalContinentAndCountry(l.prefixes.Find(resultQSO.TheirCall.String()))
		fields := exchangeFields(resultQSO.TheirContinent, resultQSO.TheirCountry)
		resultQSO.TheirExchange = toQSOExchange(fields, qso.Received)

		result[i] = resultQSO
	}
	return result
}

func toOperatorMode(operator cabrillo.CategoryOperator) conval.OperatorMode {
	if operator == cabrillo.MultiOperator {
		return conval.MultiOperator
	}
	return conval.SingleOperator
}

func toOverlay(overlay cabrillo.CategoryOverlay) conval.Overlay {
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

func toModes(mode cabrillo.CategoryMode) []conval.Mode {
	switch mode {
	case cabrillo.ModeMIXED:
		return []conval.Mode{conval.ModeALL}
	case cabrillo.ModeDIGI:
		return []conval.Mode{conval.ModeDigital}
	default:
		return []conval.Mode{conval.Mode(strings.ToLower(string(mode)))}
	}
}

func toBand(frequency cabrillo.QSOFrequency) conval.ContestBand {
	band := frequency.ToBand()
	return conval.ContestBand(strings.ToLower(string(band)))
}

func toQSOMode(mode cabrillo.QSOMode) conval.Mode {
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

func toQSOExchange(fields []conval.ExchangeField, info cabrillo.QSOInfo) conval.QSOExchange {
	values := make([]string, 0, len(info.Exchange)+1)
	values = append(values, info.RST)
	values = append(values, info.Exchange...)
	return conval.ParseExchange(fields, values)
}
