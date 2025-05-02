package helper

import (
	"errors"
	"time"
)

func StringToTtl(input string) (time.Time, error) {
	if input == "" {
		return time.Time{}, errors.New("empty string")
	}

	if input == "0" {
		return time.Time{}, nil
	}

	dur, err := time.ParseDuration(input)
	if err != nil {
		return time.Time{}, err
	}

	return time.Now().Add(dur), nil
}

func TtlToString(input time.Time) string {
	if input.IsZero() {
		return " - "
	}

	return input.Format(time.DateTime)
}
