package main

import (
	"io/ioutil"
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
	newsChannel, err := ExtractInfo(result)
	
	log.Println(newsChannel)
	
}


func ExtractInfo(doc *RssDoc)(*InfoChanel,error){
	output := InfoChanel{
		Name: string(doc.Channel.Titles[0]),
	}
	posts := make([]Post,0)
	for _,item := range doc.Channel.Items{
		
		newPost := Post{
			string(item.Titles[0]),
			string(item.Descriptions[0]),
			string(item.Links[0]),
		}
		posts = append(posts,newPost)
	}
	
	output.Posts = posts
	
	return &output, nil
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




