package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"time"

	"encoding/xml"
	"io/ioutil"
	"net/http"
	"text/template"
	"fmt"
)

func main() {
	argcount := len(os.Args)

	var command string
	if argcount > 1 {
		command = os.Args[1]
	} else {
		command = "readall"
	}

	switch command {
	case "readall":
		ReadAll()
	case "url":
		ReadUrl(os.Args[2])
		default:
		fmt.Println("Unknown command")
		fmt.Println(help)
	}

}

func ReadUrl(url string) {
	val, err := ReadNewsFrom(url)
	tmplt := template.Must(template.ParseFiles("news.tmpl"))
	if err != nil {
		log.Fatal(err)
	}
	tmplt.ExecuteTemplate(os.Stdout, "NewsTemplate", val)
}

func ReadAll() {
	sourcesFile, err := os.Open("rss.source")
	defer sourcesFile.Close()
	if err != nil {
		log.Fatal("Cant read file with rss sources", err)
	}

	output := make(chan *InfoChanel, 100)

	sync := make(chan bool)
	scanner := bufio.NewScanner(sourcesFile)

	//start asynchronous reading
	go func(c chan *InfoChanel) {
		for scanner.Scan() {
			url := scanner.Text()
			if url == "" {
				return
			}
			channelInfo, err := ReadNewsFrom(url)
			if err != nil {
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
		val, ok := <-newschannel
		if !ok {
			return
		}
		tmplt.ExecuteTemplate(os.Stdout, "NewsTemplate", val)
		time.Sleep(1)
	}
}

func ReadNewsFrom(url string) (*InfoChanel, error) {
	result, err := ReadRss(url)
	if err != nil {
		log.Fatal(err)
	}
	newsChannel, err := ExtractInfo(result)
	return newsChannel, err
}

func ExtractInfo(doc *RssDoc) (*InfoChanel, error) {
	output := InfoChanel{
		Name: string(doc.Channel.Titles[0]),
	}
	posts := make([]Post, 0)
	for _, item := range doc.Channel.Items {

		content := string(item.Descriptions[0])
				regex,err := regexp.Compile("<ul>.*</ul>|<br.*>.*<img.*")
				if err != nil {continue}
				content = regex.ReplaceAllString(content,"")

		newPost := Post{
			string(item.Titles[0]),
			content,
			string(item.Links[0]),
		}
		posts = append(posts, newPost)
	}

	output.Posts = posts

	return &output, nil
}

func ReadRss(url string) (*RssDoc, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	result := RssDoc{}
	xml.Unmarshal(content, &result)
	return &result, nil
}
