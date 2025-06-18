package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const timeOnlyLayout = "15:04"

type TimeOnly struct {
	Time *time.Time
}

// MarshalJSON implements the json.Marshaler interface and will be called automatically
func (t *TimeOnly) MarshalJSON() ([]byte, error) {
	if t.Time == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(t.Time.Format(timeOnlyLayout))
}

// UnmarshalJSON implements the json.Marshaler interface and will be called automatically
func (t *TimeOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		t.Time = nil
		return nil
	}
	parsed, err := time.Parse(timeOnlyLayout, s)
	if err != nil {
		return err
	}
	// Store with today's date to maintain consistency internally
	now := time.Now()
	parsed = time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, time.UTC)
	t.Time = &parsed
	return nil
}

// Value implements the driver.Value interface and will be called automatically
func (t *TimeOnly) Value() (driver.Value, error) {
	if t.Time == nil {
		return nil, nil
	}
	return t.Time.Format(timeOnlyLayout), nil
}

// Scan implements the sql.Scanner interface and will be called automatically
func (t *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		t.Time = nil
		return nil
	}
	switch v := value.(type) {
	case string:
		parsed, err := time.Parse(timeOnlyLayout, v)
		if err != nil {
			return err
		}
		// Normalize to today
		now := time.Now()
		parsed = time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, time.UTC)
		t.Time = &parsed
	case time.Time:
		// Strip date part
		h, m, _ := v.Clock()
		parsed := time.Date(0, 1, 1, h, m, 0, 0, time.UTC)
		t.Time = &parsed
	default:
		return fmt.Errorf("unsupported type for TimeOnly: %T", value)
	}
	return nil
}
