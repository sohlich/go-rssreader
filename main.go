package main

import (
	"io/ioutil"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"fmt"
	"bufio"
	"time"
	"text/template"
//	"bytes"
)


func main() {
	argcount := len(os.Args)

	if argcount < 2 {
		fmt.Println(help)
		return
	}

	sourcesFile, err := os.Open("rss.source", )
	defer sourcesFile.Close();
	if err != nil {log.Fatal("Cant read file with rss sources", err)}


	output := make(chan *InfoChanel, 100)

	sync := make(chan bool)


	//sources := make([]string,0)
	scanner := bufio.NewScanner(sourcesFile)



	//start asynchronous reading
	go func(c chan *InfoChanel) {
		for scanner.Scan() {
			url := scanner.Text();
			if url == "" {return }
			channelInfo, err := ReadNewsFrom(url)
			if err != nil {
				//c <- fmt.Sprintln("Failed to load news from url: "+url)
				close(c)
				return
			}
			c <- channelInfo
		}
		sync <- true
	}(output)

	go consume(output)


	<-sync

	//Clen up
	close(output)
	close(sync)
}

func consume(newschannel chan *InfoChanel) {
	tmplt := template.Must(template.ParseFiles("news.tmpl"))
	for {
		val, ok := <- newschannel
		if !ok {return }
//		buf := new(bytes.Buffer)
		tmplt.ExecuteTemplate(os.Stdout,"NewsTemplate",val)
		time.Sleep(1)
	}
}


func ReadNewsFrom(url string) (*InfoChanel, error) {
	result, err := ReadRss(url)
	if err != nil {log.Fatal(err)}
	newsChannel, err := ExtractInfo(result)
	return newsChannel, err
}

func ExtractInfo(doc *RssDoc) (*InfoChanel, error) {
	output := InfoChanel{
		Name: string(doc.Channel.Titles[0]),
	}
	posts := make([]Post, 0)
	for _, item := range doc.Channel.Items {

		newPost := Post{
			string(item.Titles[0]),
			string(item.Descriptions[0]),
			string(item.Links[0]),
		}
		posts = append(posts, newPost)
	}

	output.Posts = posts

	return &output, nil
}


func ReadRss(url string) (*RssDoc, error) {
	resp, err := http.Get(url)
	if err != nil {return nil, err}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {log.Fatal(err)}
	result := RssDoc{};
	xml.Unmarshal(content, &result)
	return &result, nil
}




