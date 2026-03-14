package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const timeOnlyLayout = "15:04"

type TimeOnly struct {
	Time *time.Time
}

func NewTimeOnly(hour, minute int) *TimeOnly {
	t := time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC)
	return &TimeOnly{Time: &t}
}

// MarshalJSON implements the json.Marshaler interface and will be called automatically
func (t *TimeOnly) MarshalJSON() ([]byte, error) {
	if t.Time == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(t.Time.Format(timeOnlyLayout))
}

// UnmarshalJSON implements the json.Marshaler interface and will be called automatically
// This method uses a pointer receiver because it modifies the internal state.
// Using a value receiver would only update a copy, leaving the original unchanged.
func (t *TimeOnly) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "null" || s == "" {
		t.Time = nil
		return nil
	}
	parsed, err := time.Parse(timeOnlyLayout, s)
	if err != nil {
		return err
	}
	// Normalize to today's date
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

func (t *TimeOnly) IsZero() bool {
	return t == nil || t.Time == nil || t.Time.IsZero()
}

func (t *TimeOnly) Equal(other *TimeOnly) bool {
	if t == nil || t.Time == nil || other == nil || other.Time == nil {
		return false
	}
	return t.Time.Hour() == other.Time.Hour() && t.Time.Minute() == other.Time.Minute()
}

func (t *TimeOnly) Before(other *TimeOnly) bool {
	if t == nil || t.Time == nil || other == nil || other.Time == nil {
		return false
	}
	return t.Time.Before(*other.Time)
}

func (t *TimeOnly) After(other *TimeOnly) bool {
	if t == nil || t.Time == nil || other == nil || other.Time == nil {
		return false
	}
	return t.Time.After(*other.Time)
}

func (t *TimeOnly) String() string {
	if t == nil || t.Time == nil {
		return "null"
	}
	return t.Time.Format(timeOnlyLayout)
}
