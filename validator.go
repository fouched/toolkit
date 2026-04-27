package toolkit

import (
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"
)

type Validation struct {
	Errors map[string]string
}

type Field struct {
	Name  string
	Label string
	Value string
}

// ValidationError contains a map structure for validation errors
type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

func (e ValidationError) Error() string {
	return "validation failed"
}

// Validator creates an instance of a Validation based on form values.
// You can pass nil to this constructor to use the validation functions without an Http form
func Validator() *Validation {
	return &Validation{
		Errors: make(map[string]string),
	}
}

// Valid returns true if the current Validator contains no errors, otherwise false
func (v *Validation) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error to the current Validator instance
func (v *Validation) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Has checks if an Http form contains a field
func (v *Validation) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	return x == ""
}

// Required checks a variadic list of Field and adds an error if a field is blank
func (v *Validation) Required(fields ...Field) {
	for _, field := range fields {
		if strings.TrimSpace(field.Value) == "" {
			v.AddError(field.Name, fmt.Sprintf("%s cannot be blank", field.Label))
		}
	}
}

// Check takes any expression that can be evaluated to a bool
// and adds an error if result is false
func (v *Validation) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// IsLength checks if a field is at least a specific length and adds an error if result is false
func (v *Validation) IsLength(field Field, length int) {
	if len(strings.TrimSpace(field.Value)) < length {
		v.AddError(field.Name, fmt.Sprintf("%s must be at least %d characters", field.Label, length))
	}
}

// IsInt checks if a Field is an integer
func (v *Validation) IsInt(field Field) {
	_, err := strconv.Atoi(field.Value)
	if err != nil {
		v.AddError(field.Name, fmt.Sprintf("%s must be an integer", field.Label))
	}
}

// IsFloat checks if a Field is a floating point number and adds an error if result is false
func (v *Validation) IsFloat(field Field) {
	_, err := strconv.ParseFloat(field.Value, 64)
	if err != nil {
		v.AddError(field.Name, fmt.Sprintf("%s must contain decimal values", field.Label))
	}
}

// IsDateISO checks if a form field is an ISO date (YYYY-MM-DD) and adds an error if result is false
func (v *Validation) IsDateISO(field Field) {
	_, err := time.Parse(time.DateOnly, field.Value)
	if err != nil {
		v.AddError(field.Name, fmt.Sprintf("%s must be a date in YYYY-MM-DD format", field.Label))
	}
}

// IsEmail checks if a Field contains a valid email and adds an error if result is false
func (v *Validation) IsEmail(field Field) {
	_, err := mail.ParseAddress(field.Value)
	if err != nil {
		v.AddError(field.Name, fmt.Sprintf("%s must be a valid email address", field.Label))
	}
}

// NoSpaces checks if a Field contains spaces and adds an error if result is false
func (v *Validation) NoSpaces(field Field) {
	if strings.Contains(field.Value, " ") {
		v.AddError(field.Name, fmt.Sprintf("%s does not allow any spaces", field.Label))
	}
}
