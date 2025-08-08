package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const dateFormat = time.DateOnly

type DateOnly struct {
	Time *time.Time
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
	s := strings.Trim(string(b), `"`)
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
	return d.Time == nil || d.Time.IsZero()
}

func (d *DateOnly) After(other DateOnly) bool {
	if d.Time == nil || other.Time == nil {
		return false // or panic/log depending on your domain
	}
	return d.Time.After(*other.Time)
}

func (d *DateOnly) Before(other DateOnly) bool {
	if d.Time == nil || other.Time == nil {
		return false
	}
	return d.Time.Before(*other.Time)
}

func (d *DateOnly) ToTime() time.Time {
	if d.Time == nil {
		return time.Time{}
	}
	return *d.Time
}

func (d *DateOnly) Add(dur time.Duration) DateOnly {
	if d.Time == nil {
		return DateOnly{}
	}
	t := d.Time.Add(dur)
	return DateOnly{Time: &t}
}
