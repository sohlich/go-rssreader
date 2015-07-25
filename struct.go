package main

import (
	"encoding/xml"
	rss "github.com/metaleap/go-xsd-pkg/thearchitect.co.uk/schemas/rss-2_0.xsd_go"
)

type RssDoc struct {
	XMLName xml.Name `xml:"rss"`
	rss.TxsdRss
}

type InfoChanel struct {
	Name  string
	Posts []Post
}

type Post struct {
	Title   string
	Content string
	Link    string
}
