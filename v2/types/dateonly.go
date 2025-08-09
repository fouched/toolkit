package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const dateFormat = time.DateOnly

type DateOnly struct {
	Time *time.Time
}

// NewDateOnly returns a DateOnly value normalized to midnight UTC.
func NewDateOnly(t time.Time) *DateOnly {
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return &DateOnly{Time: &date}
}

// MarshalJSON implements the json.Marshaler interface and will be called automatically
func (d *DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(d.Time.Format(dateFormat))
}

// UnmarshalJSON implements the json.Marshaler interface and will be called automatically
// This method uses a pointer receiver because it modifies the internal state.
// Using a value receiver would only update a copy, leaving the original unchanged.
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "null" || s == "" {
		d.Time = nil
		return nil
	}
	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return err
	}
	d.Time = &t
	return nil
}

// Value implements the driver.Value interface and will be called automatically
func (d *DateOnly) Value() (driver.Value, error) {
	if d.Time == nil {
		return nil, nil
	}
	return d.Time.Format(dateFormat), nil
}

// Scan implements the sql.Scanner interface and will be called automatically
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		d.Time = nil
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		d.Time = &v
	case string:
		t, err := time.Parse(dateFormat, v)
		if err != nil {
			return err
		}
		d.Time = &t
	default:
		return fmt.Errorf("unsupported type for DateOnly: %T", value)
	}
	return nil
}

func (d *DateOnly) IsZero() bool {
	if d == nil {
		return true
	}
	return d.Time == nil || d.Time.IsZero()
}

func (d *DateOnly) After(other *DateOnly) bool {
	if d == nil || d.Time == nil || other == nil || other.Time == nil {
		return false
	}
	return d.Time.After(*other.Time)
}

func (d *DateOnly) Before(other *DateOnly) bool {
	if d == nil || d.Time == nil || other == nil || other.Time == nil {
		return false
	}
	return d.Time.Before(*other.Time)
}

func (d *DateOnly) ToTime() time.Time {
	if d == nil || d.Time == nil {
		return time.Time{}
	}
	return *d.Time
}

func (d *DateOnly) Add(dur time.Duration) *DateOnly {
	if d == nil || d.Time == nil {
		return nil
	}
	t := d.Time.Add(dur)
	return &DateOnly{Time: &t}
}

func (d *DateOnly) StartOfDay() time.Time {
	if d == nil || d.Time == nil {
		return time.Time{}
	}
	return time.Date(d.Time.Year(), d.Time.Month(), d.Time.Day(), 0, 0, 0, 0, time.UTC)
}

func (d *DateOnly) EndOfDay() time.Time {
	if d == nil || d.Time == nil {
		return time.Time{}
	}
	return time.Date(d.Time.Year(), d.Time.Month(), d.Time.Day(), 23, 59, 59, 999, time.UTC)
}

func (d *DateOnly) String() string {
	if d == nil || d.Time == nil {
		return "null"
	}
	return d.Time.Format(dateFormat)
}

func (d *DateOnly) Equal(other *DateOnly) bool {
	if d == nil || d.Time == nil || other == nil || other.Time == nil {
		return false
	}
	return d.Time.Equal(*other.Time)
}
