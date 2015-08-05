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