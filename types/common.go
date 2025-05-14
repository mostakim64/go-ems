package types

import "time"

func parseTime(timeStr string, format string) (*time.Time, error) {
	t, err := time.Parse(format, timeStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
