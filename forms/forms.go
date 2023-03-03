package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

type Form struct {
	url.Values
	Errors errors
}

// New initialize the form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		map[string][]string{},
	}
}

func (f *Form) Has(field string) bool {
	resp := f.Get(field) //  Get request associated to the form receiver
	if resp == "" {
		return false
	}
	return true
}

// Required checks for required form fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Valid validate our form fields, returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength check for entered form field minimum length
func (f *Form) MinLength(field string, length int) bool {
	submittedField := f.Get(field) //  Get request associated to the form receiver
	if len(submittedField) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters", length))
		return false
	}
	return true

}

// IsEmail checks for valid emails
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email")
	}
}
