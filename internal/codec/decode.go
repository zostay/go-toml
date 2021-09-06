package codec

import (
	"fmt"
	"math"
	"strconv"

	"github.com/pelletier/go-toml/v2/internal/parser"
)

func DecodeInteger(b []byte) (int64, error) {
	if len(b) > 2 && b[0] == '0' {
		switch b[1] {
		case 'x':
			return decodeIntHex(b)
		case 'b':
			return decodeIntBin(b)
		case 'o':
			return decodeIntOct(b)
		default:
			panic(fmt.Errorf("invalid base '%c', should have been checked by scanIntOrFloat", b[1]))
		}
	}

	return decodeIntDec(b)
}

//nolint:cyclop
func DecodeFloat(b []byte) (float64, error) {
	if len(b) == 4 && (b[0] == '+' || b[0] == '-') && b[1] == 'n' && b[2] == 'a' && b[3] == 'n' {
		return math.NaN(), nil
	}

	cleaned, err := checkAndRemoveUnderscores(b)
	if err != nil {
		return 0, err
	}

	if cleaned[0] == '.' {
		return 0, parser.NewDecodeError(b, "float cannot start with a dot")
	}

	if cleaned[len(cleaned)-1] == '.' {
		return 0, parser.NewDecodeError(b, "float cannot end with a dot")
	}

	f, err := strconv.ParseFloat(string(cleaned), 64)
	if err != nil {
		return 0, parser.NewDecodeError(b, "unable to parse float: %w", err)
	}

	return f, nil
}

func decodeIntHex(b []byte) (int64, error) {
	cleaned, err := checkAndRemoveUnderscores(b[2:])
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(string(cleaned), 16, 64)
	if err != nil {
		return 0, parser.NewDecodeError(b, "couldn't parse hexadecimal number: %w", err)
	}

	return i, nil
}

func decodeIntOct(b []byte) (int64, error) {
	cleaned, err := checkAndRemoveUnderscores(b[2:])
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(string(cleaned), 8, 64)
	if err != nil {
		return 0, parser.NewDecodeError(b, "couldn't parse octal number: %w", err)
	}

	return i, nil
}

func decodeIntBin(b []byte) (int64, error) {
	cleaned, err := checkAndRemoveUnderscores(b[2:])
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(string(cleaned), 2, 64)
	if err != nil {
		return 0, parser.NewDecodeError(b, "couldn't parse binary number: %w", err)
	}

	return i, nil
}

func decodeIntDec(b []byte) (int64, error) {
	cleaned, err := checkAndRemoveUnderscores(b)
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(string(cleaned), 10, 64)
	if err != nil {
		return 0, parser.NewDecodeError(b, "couldn't parse decimal number: %w", err)
	}

	return i, nil
}

func checkAndRemoveUnderscores(b []byte) ([]byte, error) {
	if b[0] == '_' {
		return nil, parser.NewDecodeError(b[0:1], "number cannot start with underscore")
	}

	if b[len(b)-1] == '_' {
		return nil, parser.NewDecodeError(b[len(b)-1:], "number cannot end with underscore")
	}

	// fast path
	i := 0
	for ; i < len(b); i++ {
		if b[i] == '_' {
			break
		}
	}
	if i == len(b) {
		return b, nil
	}

	before := false
	cleaned := make([]byte, i, len(b))
	copy(cleaned, b)

	for i++; i < len(b); i++ {
		c := b[i]
		if c == '_' {
			if !before {
				return nil, parser.NewDecodeError(b[i-1:i+1], "number must have at least one digit between underscores")
			}
			before = false
		} else {
			before = true
			cleaned = append(cleaned, c)
		}
	}

	return cleaned, nil
}
