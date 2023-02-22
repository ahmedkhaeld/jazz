package jazz

import (
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Validation struct {
	Data   url.Values
	Errors map[string]string
}

//Validator initiate a validation fields,
//init the Errors to empty map, so later we can append to it
func (j *Jazz) Validator(data url.Values) *Validation {
	return &Validation{
		Data:   data,
		Errors: make(map[string]string),
	}
}

//Valid check if there is no errors in the Errors
func (v *Validation) Valid() bool {
	return len(v.Errors) == 0
}

//AddError append key with its message to the Errors map
func (v *Validation) AddError(key, message string) {
	//check if key is not exists, if so, add key and message
	_, exists := v.Errors[key]
	if !exists {
		v.Errors[key] = message
	}
}

//Has check if particular field is in Form POST request
func (v *Validation) Has(r *http.Request, field string) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

//Required makes certain fields required
func (v *Validation) Required(r *http.Request, fields ...string) {
	for _, field := range fields {
		value := r.Form.Get(field)
		//value should not be empty
		if strings.TrimSpace(value) == "" {
			v.AddError(field, "This field cannot be blank")
		}
	}
}

func (v *Validation) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

//IsEmail make sure that field's value is an email address
func (v *Validation) IsEmail(field, value string) {
	if !govalidator.IsEmail(value) {
		v.AddError(field, "invalid email address")
	}
}

func (v *Validation) IsInt(field, value string) {
	_, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(field, "this field must be an integer")
	}
}

func (v *Validation) IsFloat(field, value string) {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		v.AddError(field, "this field must be an float point number")
	}
}

func (v *Validation) IsDateISO(field, value string) {
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		v.AddError(field, "this field must be a date in the form YYYY-MM-DD")
	}
}

func (v *Validation) NoSpaces(field, value string) {
	if govalidator.HasWhitespace(value) {
		v.AddError(field, "spaces are not permitted")
	}
}
