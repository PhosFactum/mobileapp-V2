package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
)

const (
	DATE_LAYOUT = "2006-01-02"
	TIME_LAYOUT = "15:04:05"
)

var (
	datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)    // YYYY-MM-DD
	timePattern = regexp.MustCompile(`^\d{2}:\d{2}(:\d{2})?$`) // HH:MM or HH:MM:SS
)

type ParamsParser struct {
}

func NewParamsParser() interfaces.ParamsParserService {
	return &ParamsParser{}
}

func (s *ParamsParser) ParseDateString(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse(DATE_LAYOUT, strings.TrimSpace(dateStr))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected '%s': %v", DATE_LAYOUT, err)
	}
	return parsedDate, nil
}

func (s *ParamsParser) ParseTimeString(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(TIME_LAYOUT, strings.TrimSpace(timeStr))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format, expected '%s': %v", TIME_LAYOUT, err)
	}
	return parsedTime, nil
}

func (s *ParamsParser) ParseUintString(uintStr string) (uint, error) {
	uintStr = strings.TrimSpace(uintStr)
	if uintStr == "" {
		return 0, fmt.Errorf("empty string provided, expected unsigned integer")
	}

	value, err := strconv.ParseUint(uintStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid uint format, expected unsigned integer: %v", err)
	}

	return uint(value), nil
}

func (s *ParamsParser) ParseIntString(intStr string) (int, error) {
	intStr = strings.TrimSpace(intStr)
	if intStr == "" {
		return 0, fmt.Errorf("empty string provided, expected integer")
	}

	value, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid int format, expected signed integer: %v", err)
	}

	return int(value), nil
}

func (s *ParamsParser) FormatDateToString(t time.Time) string {
	return t.Format(DATE_LAYOUT)
}

func (s *ParamsParser) FormatTimeToString(t time.Time) string {
	return t.Format(TIME_LAYOUT)
}

func ParseUint(value interface{}) (uint, error) {
	switch v := value.(type) {
	case uint:
		return v, nil
	case int:
		if v >= 0 {
			return uint(v), nil
		}
	case int64:
		if v >= 0 {
			return uint(v), nil
		}
	case float64:
		if v >= 0 && v == float64(uint(v)) {
			return uint(v), nil
		}
	case string:
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			return uint(id), nil
		}
	}
	return 0, fmt.Errorf("cannot convert %v (type %T) to uint", value, value)
}
