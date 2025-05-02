package helper

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func StringToDuration(input string) (time.Duration, error) {
	re := regexp.MustCompile(`(?:(\d+)d)?(?:(\d+)h)?(?:(\d+)m)?`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return 0, fmt.Errorf("invalid format: %s", input)
	}

	var days, hours, minutes int64
	var err error

	if matches[1] != "" {
		days, err = strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return 0, err
		}
	}
	if matches[2] != "" {
		hours, err = strconv.ParseInt(matches[2], 10, 64)
		if err != nil {
			return 0, err
		}
	}
	if matches[3] != "" {
		minutes, err = strconv.ParseInt(matches[3], 10, 64)
		if err != nil {
			return 0, err
		}
	}

	duration := time.Duration(days)*24*time.Hour +
		time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute

	return duration, nil
}
