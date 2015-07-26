package main

import (
	"testing"
	"fmt"
	"io/ioutil"
	"reflect"
	"bytes"
	"github.com/stretchr/testify/assert"
)


func TestRemoveAllHtml(t *testing.T) {
	testString := `Hello<ul><li>bla bla<li></ul>`
	expected := "Hello"
	result, err := RemoveAllHtml(testString);
	IfError(err, t)
	if result != expected {
		t.Error("Expected : "+expected+" but got "+result)
	}
}


func TestReadRss(t *testing.T) {
	fmt.Println("Test ReadRss")
	content, err := ioutil.ReadFile("test/rss_test.xml")
	IfError(err, t)
	result, err := ReadRss(bytes.NewBuffer(content))
	IfError(err, t)
	expectedTitle := "Business Insider India"
	assert.Equal(t, expectedTitle, string(result.Channel.Titles[0]))
}


func IfError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}



// Restorer holds a function that can be used
// to restore some previous state.
type Restorer func()

// Restore restores some previous state.
func (r Restorer) Restore() {
	r()
}

// Patch sets the value pointed to by the given destination to the given
// value, and returns a function to restore it to its original value.  The
// value must be assignable to the element type of the destination.
func Patch(dest, value interface{}) Restorer {
	destv := reflect.ValueOf(dest).Elem()
	oldv := reflect.New(destv.Type()).Elem()
	oldv.Set(destv)
	valuev := reflect.ValueOf(value)
	if !valuev.IsValid() {
		// This isn't quite right when the destination type is not
		// nilable, but it's better than the complex alternative.
		valuev = reflect.Zero(destv.Type())
	}
	destv.Set(valuev)
	return func() {
		destv.Set(oldv)
	}
}