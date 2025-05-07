package helper

import (
	"fmt"
	"time"
)

func StringToTtl(input string) (time.Time, error) {
	const op = "helper.StringToTtl"

	if input == "" {
		return time.Time{}, fmt.Errorf("%s: empty input", op)
	}

	if input == "0" {
		return time.Time{}, nil
	}

	dur, err := time.ParseDuration(input)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: %w", op, err)
	}

	return time.Now().Add(dur), nil
}

func TtlToString(input time.Time) string {
	if input.IsZero() {
		return " - "
	}

	return input.Format(time.DateTime)
}
