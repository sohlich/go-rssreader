package main

import (
	"io/ioutil"
	"bytes"
	"fmt"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	rss "github.com/metaleap/go-xsd-pkg/thearchitect.co.uk/schemas/rss-2_0.xsd_go"
)

type RssDoc struct {
    XMLName xml.Name `xml:"rss"`
    rss.TxsdRss
}


type InfoChanel struct {
	Name string
	Posts []Post
}


type Post struct{
	Title string
	Content string
	Link string
}

func main(){
		
	result,err := ReadRss("http://servis.idnes.cz/rss.aspx?c=zpravodaj")
	
	if err != nil {log.Fatal(err)}	
	buf := new(bytes.Buffer)
	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.Encode(result)
	
	output := InfoChanel{
		Name: string(result.Channel.Titles[0]),
	}
	posts := make([]Post,0)
	for _,item := range result.Channel.Items{
		
		newPost := Post{
			string(item.Titles[0]),
			string(item.Descriptions[0]),
			string(item.Links[0]),
		}
		posts = append(posts,newPost)
	}
	
	output.Posts = posts
	
	
}


func ExtractInfo(doc RssDoc)(*InfoChanel,error){
	
}


func ReadRss(url string) (*RssDoc, error){
	resp, err := http.Get(url)
	if err != nil {return nil,err}	
	content, err :=ioutil.ReadAll(resp.Body)
	if err != nil {log.Fatal(err)}	
	result := RssDoc{};
	xml.Unmarshal(content,&result)
	return &result, nil
}




