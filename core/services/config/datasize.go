package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var dataSizePattern = regexp.MustCompile(`^"?(\d+)(b|kb|mb)?"?$`)

var (
	ErrUnmatched       = errors.New("data size string does not match expected pattern")
	ErrUnsupportedUnit = errors.New("unsupported unit")
)

const (
	KiloByte DataSize = 1024
	MegaByte          = KiloByte * 1024
)

func ParseDataSize(raw string) (DataSize, error) {
	matches := dataSizePattern.FindStringSubmatch(raw)
	if numMatches := len(matches); numMatches != 3 {
		return 0, fmt.Errorf("%w: %s", ErrUnmatched, raw)
	}

	multiplicator, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing multiplicator %s: %w", matches[1], err)
	}

	switch matches[2] {
	case "", "b":
		return DataSize(multiplicator), nil
	case "kb":
		return DataSize(multiplicator) * KiloByte, nil
	case "mb":
		return DataSize(multiplicator) * MegaByte, nil
	default:
		return 0, fmt.Errorf("%w: %s", ErrUnsupportedUnit, matches[2])
	}
}

var _ json.Unmarshaler = (*DataSize)(nil)

type DataSize int64

func (d DataSize) Bytes() int64 {
	return int64(d)
}

func (d *DataSize) UnmarshalJSON(bytes []byte) error {
	parsed, err := ParseDataSize(string(bytes))
	if err != nil {
		return err
	}

	*d = parsed
	return nil
}
