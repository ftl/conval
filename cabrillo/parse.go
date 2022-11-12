package cabrillo

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ftl/hamradio/callsign"
	"github.com/ftl/hamradio/locator"
)

func Read(r io.Reader) (*Log, error) {
	result := NewLog()
	parser := newParser(result)

	lineScanner := bufio.NewScanner(r)
	for lineScanner.Scan() {
		line := lineScanner.Text()
		err := parser.AddLine(line)
		if err != nil {
			return nil, err
		}
	}
	scanErr := lineScanner.Err()
	if scanErr != nil {
		return nil, scanErr
	}
	parserErr := parser.CheckComplete()
	if parserErr != nil {
		return nil, parserErr
	}

	return result, nil
}

func newParser(log *Log) *parser {
	return &parser{log: log}
}

type parser struct {
	log        *Log
	lineNumber int
	started    bool
	ended      bool
}

func (p *parser) AddLine(line string) error {
	p.lineNumber++
	if line == "" {
		return nil
	}
	withinLog := p.started && !p.ended
	tagStr, value, found := strings.Cut(line, ":")
	if !found {
		// ignore any lines outside the start and end tags
		if !withinLog {
			return nil
		}
		return p.lineErrorf("%s is not a valid Cabrillo log line", line)
	}
	tag := Tag(strings.ToUpper(strings.TrimSpace(tagStr)))
	value = strings.TrimSpace(value)

	switch tag {
	case StartOfLogTag:
		if withinLog {
			return p.lineErrorf("the log already started in a former line")
		}
		p.started = true
		p.log.CabrilloVersion = value
	case EndOfLogTag:
		p.ended = true
	default:
		return p.parseTag(tag, value)
	}

	return nil
}

func (p *parser) lineErrorf(format string, args ...any) error {
	prefix := fmt.Sprintf("line %d", p.lineNumber)
	message := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %s", prefix, message)
}

func (p *parser) parseTag(tag Tag, value string) error {
	tagParser, found := tagParsers[tag]
	if !found {
		p.appendCustomValue(tag, value)
		return nil
	}
	return tagParser.Parse(p.log, value)
}

func (p *parser) appendCustomValue(tag Tag, value string) {
	currentValue, found := p.log.Custom[tag]
	var newValue string
	if found {
		newValue = currentValue + "\n" + value
	} else {
		newValue = value
	}
	p.log.Custom[tag] = newValue
}

func (p *parser) CheckComplete() error {
	if !p.started {
		return fmt.Errorf("no START-OF-LOG tag found")
	}
	if !p.ended {
		return fmt.Errorf("no END-OF-LOG tag found")
	}

	return nil
}

type tagParser interface {
	Parse(*Log, string) error
}

type tagParserFunc func(*Log, string) error

func (f tagParserFunc) Parse(log *Log, value string) error {
	return f(log, value)
}

var (
	operatorsSeparator = regexp.MustCompile(`\s*(,\s*|\s+)`)
	validOfftime       = regexp.MustCompile(`([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{4})\s+([0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{4})`)
	qsoColumnSeparator = regexp.MustCompile(`\s+`)
)

var tagParsers = map[Tag]tagParser{
	CallsignTag: tagParserFunc(func(log *Log, value string) error {
		cs, err := callsign.Parse(value)
		if err != nil {
			return err
		}
		log.Callsign = cs
		return nil
	}),
	ContestTag: tagParserFunc(func(log *Log, value string) error {
		log.Contest = ContestIdentifier(strings.ToUpper(value))
		return nil
	}),
	CategoryAssistedTag: tagParserFunc(func(log *Log, value string) error {
		log.Category.Assisted = CategoryAssisted(strings.ToUpper(value))
		return nil
	}),
	CategoryBandTag: tagParserFunc(func(log *Log, value string) error {
		log.Category.Band = CategoryBand(strings.ToUpper(value))
		return nil
	}),
	CategoryModeTag: tagParserFunc(func(log *Log, value string) error {
		log.Category.Mode = CategoryMode(strings.ToUpper(value))
		return nil
	}),
	CategoryOperatorTag: tagParserFunc(func(log *Log, value string) error {
		log.Category.Operator = CategoryOperator(strings.ToUpper(value))
		return nil
	}),
	CategoryPowerTag: tagParserFunc(func(log *Log, value string) error {
		log.Category.Power = CategoryPower(strings.ToUpper(value))
		return nil
	}),
	CategoryStationTag: tagParserFunc(func(log *Log, value string) error {
		log.Category.Station = CategoryStation(strings.ToUpper(value))
		return nil
	}),
	CategoryTimeTag: tagParserFunc(func(log *Log, value string) error {
		log.Category.Time = CategoryTime(strings.ToUpper(value))
		return nil
	}),
	CategoryTransmitterTag: tagParserFunc(func(log *Log, value string) error {
		log.Category.Transmitter = CategoryTransmitter(strings.ToUpper(value))
		return nil
	}),
	CategoryOverlayTag: tagParserFunc(func(log *Log, value string) error {
		log.Category.Overlay = CategoryOverlay(strings.ToUpper(value))
		return nil
	}),
	CertificateTag: tagParserFunc(func(log *Log, value string) error {
		value = strings.ToUpper(value)
		switch value {
		case "1", "TRUE", "YES":
			log.Certificate = true
		default:
			log.Certificate = false
		}
		return nil
	}),
	ClaimedScoreTag: tagParserFunc(func(log *Log, value string) error {
		if value == "" {
			return nil
		}
		score, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		log.ClaimedScore = score
		return nil
	}),
	ClubTag: tagParserFunc(func(log *Log, value string) error {
		log.Club = value
		return nil
	}),
	CratedByTag: tagParserFunc(func(log *Log, value string) error {
		log.CreatedBy = value
		return nil
	}),
	EmailTag: tagParserFunc(func(log *Log, value string) error {
		log.Email = value
		return nil
	}),
	GridLocatorTag: tagParserFunc(func(log *Log, value string) error {
		if value == "" {
			return nil
		}
		result, err := locator.Parse(value)
		if err != nil {
			return err
		}
		log.GridLocator = result
		return nil
	}),
	LocationTag: tagParserFunc(func(log *Log, value string) error {
		log.Location = value
		return nil
	}),
	NameTag: tagParserFunc(func(log *Log, value string) error {
		log.Name = value
		return nil
	}),
	AddressTag: tagParserFunc(func(log *Log, value string) error {
		log.Address.Text = value
		return nil
	}),
	AddressCityTag: tagParserFunc(func(log *Log, value string) error {
		log.Address.City = value
		return nil
	}),
	AddressStateProviceTag: tagParserFunc(func(log *Log, value string) error {
		log.Address.StateProvince = value
		return nil
	}),
	AddressPostalcodeTag: tagParserFunc(func(log *Log, value string) error {
		log.Address.Postalcode = value
		return nil
	}),
	AddressCountryTag: tagParserFunc(func(log *Log, value string) error {
		log.Address.Country = value
		return nil
	}),
	OperatorsTag: tagParserFunc(func(log *Log, value string) error {
		parts := operatorsSeparator.Split(value, -1)
		operators := make([]callsign.Callsign, 0, len(parts))
		for _, part := range parts {
			if part == "" {
				continue
			}
			isHost := false
			if strings.HasPrefix(part, "@") {
				part = part[1:]
				isHost = true
			}
			call, err := callsign.Parse(part)
			if err != nil {
				return err
			}
			operators = append(operators, call)
			if isHost {
				log.Host = call
			}
		}
		log.Operators = operators
		return nil
	}),
	OfftimeTag: tagParserFunc(func(log *Log, value string) error {
		if value == "" {
			return nil
		}

		matches := validOfftime.FindStringSubmatch(value)
		if len(matches) != 3 {
			return fmt.Errorf("%s is not a valid offtime value", value)
		}

		begin, err := ParseTimestamp(matches[1])
		if err != nil {
			return fmt.Errorf("the offtime begin is not a valid timestamp: %w", err)
		}
		end, err := ParseTimestamp(matches[2])
		if err != nil {
			return fmt.Errorf("the offtime end is not a valid timestamp: %w", err)
		}

		log.Offtime.Begin = begin
		log.Offtime.End = end

		return nil
	}),
	SoapboxTag: tagParserFunc(func(log *Log, value string) error {
		if log.Soapbox != "" {
			log.Soapbox += "\n"
		}
		log.Soapbox += value
		return nil
	}),
	QSOTag: tagParserFunc(func(log *Log, value string) error {
		qso, err := ParseQSO(value)
		if err != nil {
			return err
		}
		log.QSOData = append(log.QSOData, qso)
		return nil
	}),
	XQSOTag: tagParserFunc(func(log *Log, value string) error {
		qso, err := ParseQSO(value)
		if err != nil {
			return err
		}
		log.IgnoredQSOs = append(log.IgnoredQSOs, qso)
		return nil
	}),
}

func ParseTimestamp(s string) (time.Time, error) {
	const timestampLayout = "2006-01-02 1504"
	return time.Parse(timestampLayout, s)
}

func ParseQSO(s string) (QSO, error) {
	columns := qsoColumnSeparator.Split(s, -1)
	if len(columns) < 8 { // need at least a callsign and rst per side
		return QSO{}, fmt.Errorf("not enough QSO columns: %d", len(columns))
	}

	timestamp, err := ParseTimestamp(columns[2] + " " + columns[3])
	if err != nil {
		return QSO{}, err
	}

	hasTransmitterColumn := (len(columns)%2 == 1)
	var qsoInfoLength int
	if hasTransmitterColumn {
		qsoInfoLength = (len(columns) - 5) / 2
	} else {
		qsoInfoLength = (len(columns) - 4) / 2
	}

	sentInfo, err := parseQSOInfo(columns[4 : 4+qsoInfoLength])
	if err != nil {
		return QSO{}, err
	}
	receivedInfo, err := parseQSOInfo(columns[4+qsoInfoLength : 4+2*qsoInfoLength])
	if err != nil {
		return QSO{}, err
	}
	var transmitter int
	if hasTransmitterColumn {
		transmitter, err = strconv.Atoi(columns[len(columns)-1])
		if err != nil {
			return QSO{}, err
		}
	}

	var result QSO
	result.Frequency = QSOFrequency(columns[0])
	result.Mode = QSOMode(columns[1])
	result.Timestamp = timestamp
	result.Sent = sentInfo
	result.Received = receivedInfo
	result.Transmitter = transmitter

	return result, nil
}

func parseQSOInfo(columns []string) (QSOInfo, error) {
	if len(columns) < 2 {
		return QSOInfo{}, fmt.Errorf("not enough QSO info columns: %d", len(columns))
	}
	var result QSOInfo
	var err error

	result.Call, err = callsign.Parse(columns[0])
	if err != nil {
		return QSOInfo{}, err
	}
	result.RST = columns[1]
	result.Exchange = append([]string{}, columns[2:]...)

	return result, nil
}
