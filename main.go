package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"

	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

const (
	VERSION           = "v0.1.0"
	urlValidatorRegex = "^((ftp|http|https):\\/\\/)?([a-zA-Z0-9]+(\\.[a-zA-Z0-9]+)+.*)$"
	htmlReplacerRegex = "<[^>]*>.*</[^>]*>|<[^>]*>"
)

var (
	urlValidator *regexp.Regexp     = regexp.MustCompile(urlValidatorRegex)           //regex validator for URL
	htmlReplacer *regexp.Regexp     = regexp.MustCompile(htmlReplacerRegex)           //regex replacer for HTML
	tmplt        *template.Template = template.Must(template.ParseFiles("news.tmpl")) //Parsed template
)

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
		val := c.String("url")
		if val != "" {
			ReadUrl(val)
		} else {
			ReadAll()
		}
	}

	app.Run(os.Args)
}

//Loads url strings from source file
func parseSourceFile(pathToFile string) (<-chan string, error) {
	sourcesFile, err := os.Open(pathToFile)
	// defer sourcesFile.Close()
	if err != nil {
		log.Fatal("Cant read file with rss sources", err)
	}
	output := make(chan string)
	scanner := bufio.NewScanner(sourcesFile)
	go func(sourcesFile *os.File) {
		for scanner.Scan() {
			url := scanner.Text()
			if urlValidator.MatchString(url) {
				output <- url
			}
		}
		close(output)
		sourcesFile.Close()
	}(sourcesFile)
	return output, err
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
	urlchan, err := parseSourceFile("rss.source")
	if err != nil {
		log.Fatal(err)
	}

	for {
		val, ok := <-urlchan
		if ok {
			go func() { ReadUrl(val) }()
		} else {
			break
		}
	}

	var input string
	fmt.Scanln(&input)
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
