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
	"github.com/ftl/conval/cmd/score"
)

var scoreFlags = struct {
	setupFilename      string
	definitionFilename string
	cabrilloName       string
	verbose            bool
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
	scoreCmd.Flags().BoolVar(&scoreFlags.verbose, "verbose", false, "enable verbose output")

	rootCmd.AddCommand(scoreCmd)
}

func runScore(cmd *cobra.Command, args []string) {
	var err error
	prefixes, err := newPrefixDatabase()
	if err != nil {
		log.Fatal(err)
	}
	definition, err := score.PrepareDefinition(scoreFlags.definitionFilename, scoreFlags.cabrilloName)
	if err != nil {
		log.Fatal(err)
	}
	setup, err := score.PrepareSetup(scoreFlags.setupFilename, prefixes)
	if err != nil {
		log.Fatal(err)
	}

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

		// print the total score
		totalScore := counter.TotalScore()
		if scoreFlags.verbose {
			fmt.Printf("QSOs   : % 8d\nMultis : % 8d\nPoints : % 8d\nTotal  : % 8d\n", totalScore.QSOs, totalScore.Multis, totalScore.Points, totalScore.Total())
		} else {
			fmt.Printf("%d\n", totalScore.Multis*totalScore.Points)
		}

		// print the multis per band
		if scoreFlags.verbose {
			properties := counter.MultiProperties()
			bands := append(counter.UsedBands(), conval.BandAll)
			for _, property := range properties {
				fmt.Printf("%s:\n", property)
				for _, band := range bands {
					multis := counter.MultisPerBand(band, property)
					if len(multis) == 0 {
						continue
					}
					fmt.Printf("% 5s (% 3d): %s\n", band, len(multis), strings.ToUpper(strings.Join(multis, ", ")))
				}
				fmt.Printf("\n")
			}
		}
	}
}

func newPrefixDatabase() (*prefixDatabase, error) {
	prefixes, _, err := dxcc.DefaultPrefixes(true)
	if err != nil {
		return nil, err
	}
	return &prefixDatabase{prefixes}, nil
}

type prefixDatabase struct {
	prefixes *dxcc.Prefixes
}

func (d prefixDatabase) Find(s string) (conval.Continent, conval.DXCCEntity, bool) {
	entities, found := d.prefixes.Find(s)
	if !found || len(entities) == 0 {
		return "", "", false
	}

	return conval.Continent(strings.ToLower(entities[0].Continent)), conval.DXCCEntity(strings.ToLower(entities[0].PrimaryPrefix)), true
}

func readCabrilloLogFromFile(filename string, prefixes conval.PrefixDatabase) (*cabrilloLogfile, error) {
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
	prefixes conval.PrefixDatabase
}

func (l cabrilloLogfile) Identifier() conval.ContestIdentifier {
	return conval.ContestIdentifier(l.log.Contest)
}

func (l cabrilloLogfile) Setup() *conval.Setup {
	result := new(conval.Setup)
	result.MyCall = l.log.Callsign
	myContinent, myCountry, found := l.prefixes.Find(result.MyCall.String())
	if found {
		result.MyContinent = myContinent
		result.MyCountry = myCountry
	}

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
		theirContinent, theirCountry, found := l.prefixes.Find(resultQSO.TheirCall.String())
		if found {
			resultQSO.TheirContinent = theirContinent
			resultQSO.TheirCountry = theirCountry
		}
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
