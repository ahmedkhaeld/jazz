package render

import (
	"html/template"
	"time"
)

// var functions a map of functions that can be used in templates e.g. format a date
// some time we will create our own functions and pass them to the template
var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

// HumanDate returns time in YYYY-MM-DD format
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDate formats time to string
func FormatDate(t time.Time, f string) string {
	return t.Format(f)

}

// Iterate returns a slice of ints, starting from i , end to count
func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

func Add(a, b int) int {
	return a + b
}
