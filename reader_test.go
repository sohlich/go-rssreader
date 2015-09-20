package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestRemoveAllHtml(t *testing.T) {
	testString := `Hello<ul><li>bla bla<li></ul>`
	expected := "Hello"
	result := RemoveAllHtml(testString)
	if result != expected {
		t.Error("Expected : " + expected + " but got " + result)
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
