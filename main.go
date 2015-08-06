package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"time"

	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/codegangsta/cli"
)

//Regex for url validation
const urlValidatorRegex = "^((ftp|http|https):\\/\\/)?([a-zA-Z0-9]+(\\.[a-zA-Z0-9]+)+.*)$"

//Regex to remove all html from content
const htmlReplacerRegex = "<[^>]*>.*</[^>]*>|<[^>]*>"

//regex validator for URL
var urlValidator *regexp.Regexp = regexp.MustCompile(urlValidatorRegex)

//regex replacer for HTML
var htmlReplacer *regexp.Regexp = regexp.MustCompile(htmlReplacerRegex)

//Parsed template
var tmplt *template.Template = template.Must(template.ParseFiles("news.tmpl"))

func main() {

	app := cli.NewApp()
	app.Name = "rssreader"
	app.Usage = "read the rss feed to command line"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Usage: "url to read from",
		},
	}

	app.Action = func(c *cli.Context) {
		fmt.Printf("Url %s", c.String("url"))

		val := c.String("url")

		if val != "" {
			fmt.Println("Reading one URL")
			ReadUrl(val)
		} else {
			fmt.Println("Raeading all sources")
			ReadAll()
		}
	}

	app.Run(os.Args)
}

//Loads url strings from source file
func parseSourceFile(pathToFile string) ([]string, error) {
	sourcesFile, err := os.Open(pathToFile)
	defer sourcesFile.Close()
	if err != nil {
		log.Fatal("Cant read file with rss sources", err)
	}
	scanner := bufio.NewScanner(sourcesFile)
	urlList := []string{}
	for scanner.Scan() {
		url := scanner.Text()
		if urlValidator.MatchString(url) {
			urlList = append(urlList, url)
		}
	}
	return urlList, err
}

//Read one rss source
func ReadUrl(url string) {
	val, err := ReadNewsFrom(url)
	if err != nil {
		log.Fatalln(err)
	}
	renderToSTDOUT(val)
}

//Read all sources from rss.source file
func ReadAll() {

	urlList, err := parseSourceFile("rss.source")
	log.Println(urlList)
	if err != nil {
		log.Fatal(err)
	}

	output := make(chan *InfoChanel, 100)

	sync := make(chan bool)

	//asynchronous url reader
	go func(c chan *InfoChanel) {
		for _, url := range urlList {
			// url := scanner.Text()
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

	//asynchronous content consumer
	go consume(output)

	<-sync //wai until all done

	//Clean up
	close(output)
	close(sync)
}

//Asynchronously consumes the content comming from gorutine
//the content is then passed to template and written to command line
func consume(newschannel chan *InfoChanel) {
	for {
		val, ok := <-newschannel
		if !ok {
			return
		}
		renderToSTDOUT(val)
		time.Sleep(1)
	}
}

func ReadNewsFrom(url string) (*InfoChanel, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	result, err := ReadRss(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	newsChannel, err := ExtractInfo(result)
	return newsChannel, err
}

//Reads informations from parsed RSSDoc
func ExtractInfo(doc *RssDoc) (*InfoChanel, error) {
	output := InfoChanel{
		Name: string(doc.Channel.Titles[0]),
	}
	posts := make([]Post, 0)
	for _, item := range doc.Channel.Items {
		content := RemoveAllHtml(string(item.Descriptions[0]))
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

//Removes all html with its content so only plain
//text will survive.
func RemoveAllHtml(content string) string {
	return htmlReplacer.ReplaceAllString(content, "")
}

//Parse rss feed to RssDoc
func ReadRss(reader io.Reader) (*RssDoc, error) {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	result := RssDoc{}
	xml.Unmarshal(content, &result)
	return &result, nil
}

//Renders template to STDOUT
func renderToSTDOUT(post interface{}) {
	tmplt.ExecuteTemplate(os.Stdout, "NewsTemplate", post)
}
