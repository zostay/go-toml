package toml

import (
	"fmt"
	"time"

	"github.com/pelletier/go-toml/v2/internal/parser"
)

// LocalDate represents a calendar day in no specific timezone.
type LocalDate struct {
	Year  int
	Month int
	Day   int
}

// AsTime converts d into a specific time instance at midnight in zone.
func (d LocalDate) AsTime(zone *time.Location) time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, zone)
}

// String returns RFC 3339 representation of d.
func (d LocalDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// MarshalText returns RFC 3339 representation of d.
func (d LocalDate) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText parses b using RFC 3339 to fill d.
func (d *LocalDate) UnmarshalText(b []byte) error {
	res, err := decodeLocalDate(b)
	if err != nil {
		return err
	}
	*d = res
	return nil
}

// LocalTime represents a time of day of no specific day in no specific
// timezone.
type LocalTime struct {
	Hour       int
	Minute     int
	Second     int
	Nanosecond int
}

// String returns RFC 3339 representation of d.
func (d LocalTime) String() string {
	s := fmt.Sprintf("%02d:%02d:%02d", d.Hour, d.Minute, d.Second)
	if d.Nanosecond == 0 {
		return s
	}
	return s + fmt.Sprintf(".%09d", d.Nanosecond)
}

// MarshalText returns RFC 3339 representation of d.
func (d LocalTime) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText parses b using RFC 3339 to fill d.
func (d *LocalTime) UnmarshalText(b []byte) error {
	res, left, err := decodeLocalTime(b)
	if err == nil && len(left) != 0 {
		err = parser.NewDecodeError(left, "extra characters")
	}
	if err != nil {
		return err
	}
	*d = res
	return nil
}

// LocalDateTime represents a time of a specific day in no specific timezone.
type LocalDateTime struct {
	LocalDate
	LocalTime
}

// AsTime converts d into a specific time instance in zone.
func (d LocalDateTime) AsTime(zone *time.Location) time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, d.Hour, d.Minute, d.Second, d.Nanosecond, zone)
}

// String returns RFC 3339 representation of d.
func (d LocalDateTime) String() string {
	return d.LocalDate.String() + " " + d.LocalTime.String()
}

// MarshalText returns RFC 3339 representation of d.
func (d LocalDateTime) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText parses b using RFC 3339 to fill d.
func (d *LocalDateTime) UnmarshalText(data []byte) error {
	res, left, err := decodeLocalDateTime(data)
	if err == nil && len(left) != 0 {
		err = parser.NewDecodeError(left, "extra characters")
	}
	if err != nil {
		return err
	}

	*d = res
	return nil
}

func decodeLocalDate(b []byte) (LocalDate, error) {
	// full-date      = date-fullyear "-" date-month "-" date-mday
	// date-fullyear  = 4DIGIT
	// date-month     = 2DIGIT  ; 01-12
	// date-mday      = 2DIGIT  ; 01-28, 01-29, 01-30, 01-31 based on month/year
	var date LocalDate

	if len(b) != 10 || b[4] != '-' || b[7] != '-' {
		return date, parser.NewDecodeError(b, "dates are expected to have the format YYYY-MM-DD")
	}

	date.Year = parseDecimalDigits(b[0:4])

	v := parseDecimalDigits(b[5:7])

	date.Month = v

	date.Day = parseDecimalDigits(b[8:10])

	return date, nil
}

// decodeLocalTime is a bit different because it also returns the remaining
// []byte that is didn't need. This is to allow darseDateTime to parse those
// remaining bytes as a timezone.
func decodeLocalTime(b []byte) (LocalTime, []byte, error) {
	var (
		nspow = [10]int{0, 1e8, 1e7, 1e6, 1e5, 1e4, 1e3, 1e2, 1e1, 1e0}
		t     LocalTime
	)

	const localTimeByteLen = 8
	if len(b) < localTimeByteLen {
		return t, nil, parser.NewDecodeError(b, "times are expected to have the format HH:MM:SS[.NNNNNN]")
	}

	t.Hour = parseDecimalDigits(b[0:2])
	if b[2] != ':' {
		return t, nil, parser.NewDecodeError(b[2:3], "expecting colon between hours and minutes")
	}

	t.Minute = parseDecimalDigits(b[3:5])
	if b[5] != ':' {
		return t, nil, parser.NewDecodeError(b[5:6], "expecting colon between minutes and seconds")
	}

	t.Second = parseDecimalDigits(b[6:8])

	const minLengthWithFrac = 9
	if len(b) >= minLengthWithFrac && b[minLengthWithFrac-1] == '.' {
		frac := 0
		digits := 0

		for i, c := range b[minLengthWithFrac:] {
			if !parser.IsDigit(c) {
				if i == 0 {
					return t, nil, parser.NewDecodeError(b[i:i+1], "need at least one digit after fraction point")
				}

				break
			}

			const maxFracPrecision = 9
			if i >= maxFracPrecision {
				return t, nil, parser.NewDecodeError(b[i:i+1], "maximum precision for date time is nanosecond")
			}

			frac *= 10
			frac += int(c - '0')
			digits++
		}

		t.Nanosecond = frac * nspow[digits]

		return t, b[9+digits:], nil
	}

	return t, b[8:], nil
}

func decodeLocalDateTime(b []byte) (LocalDateTime, []byte, error) {
	var dt LocalDateTime

	const localDateTimeByteMinLen = 11
	if len(b) < localDateTimeByteMinLen {
		return dt, nil, parser.NewDecodeError(b, "local datetimes are expected to have the format YYYY-MM-DDTHH:MM:SS[.NNNNNNNNN]")
	}

	date, err := decodeLocalDate(b[:10])
	if err != nil {
		return dt, nil, err
	}
	dt.LocalDate = date

	sep := b[10]
	if sep != 'T' && sep != ' ' {
		return dt, nil, parser.NewDecodeError(b[10:11], "datetime separator is expected to be T or a space")
	}

	t, rest, err := decodeLocalTime(b[11:])
	if err != nil {
		return dt, nil, err
	}
	dt.LocalTime = t

	return dt, rest, nil
}

func darseDateTime(b []byte) (time.Time, error) {
	// offset-date-time = full-date time-delim full-time
	// full-time      = partial-time time-offset
	// time-offset    = "Z" / time-numoffset
	// time-numoffset = ( "+" / "-" ) time-hour ":" time-minute

	dt, b, err := decodeLocalDateTime(b)
	if err != nil {
		return time.Time{}, err
	}

	var zone *time.Location

	if len(b) == 0 {
		// parser should have checked that when assigning the date time node
		panic("date time should have a timezone")
	}

	if b[0] == 'Z' {
		b = b[1:]
		zone = time.UTC
	} else {
		const dateTimeByteLen = 6
		if len(b) != dateTimeByteLen {
			return time.Time{}, parser.NewDecodeError(b, "invalid date-time timezone")
		}
		direction := 1
		if b[0] == '-' {
			direction = -1
		}

		hours := parser.DigitsToInt(b[1:3])
		minutes := parser.DigitsToInt(b[4:6])
		seconds := direction * (hours*3600 + minutes*60)
		zone = time.FixedZone("", seconds)
		b = b[dateTimeByteLen:]
	}

	if len(b) > 0 {
		return time.Time{}, parser.NewDecodeError(b, "extra bytes at the end of the timezone")
	}

	t := time.Date(
		dt.Year,
		time.Month(dt.Month),
		dt.Day,
		dt.Hour,
		dt.Minute,
		dt.Second,
		dt.Nanosecond,
		zone)

	return t, nil
}

func parseDecimalDigits(b []byte) int {
	v := 0

	for _, c := range b {
		v *= 10
		v += int(c - '0')
	}

	return v
}
