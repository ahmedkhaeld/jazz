package session

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"reflect"
	"testing"
)

func TestSession_New(t *testing.T) {
	//to call the New() method we need to create a Session struct instance first
	//with the required fields values to be passed to the New() method
	s := &Session{
		Cookie: Cookie{
			LifeTime: "100",
			Persist:  "true",
			Name:     "jazz",
			Domain:   "localhost",
		},
		SessionType: "cookie",
	}
	// mySM should be a scs.SessionManager type
	mySM := s.New()
	//we need to make sure that we've got the reasonable return value from the New() method

	var mySMKind reflect.Kind
	var mySMType reflect.Type
	refValue := reflect.ValueOf(mySM)

	var sm *scs.SessionManager

	// loop through all the kinds
	for refValue.Kind() == reflect.Ptr || refValue.Kind() == reflect.Interface {
		fmt.Println("For loop:", "kind:", refValue.Kind(), "type:", refValue.Type(), "value:", refValue)
		mySMKind = refValue.Kind()
		mySMType = refValue.Type()
		refValue = refValue.Elem()
	}
	// check if the value is valid
	if !refValue.IsValid() {
		t.Error("invalid type or kind; kind:", refValue.Kind(),
			"type:", refValue.Type())
	}
	// compare the kind and type of the returned value with the expected value (sm) kind and type
	if mySMKind != reflect.ValueOf(sm).Kind() {
		t.Error("wrong kind returned testing cookie session. "+
			"Expected", reflect.ValueOf(sm).Kind(), "and got", mySMKind)
	}
	if mySMType != reflect.ValueOf(sm).Type() {
		t.Error("wrong kind returned testing cookie session."+
			" Expected", reflect.ValueOf(sm).Kind(), "and got", mySMKind)
	}

}
