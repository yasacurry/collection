package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

var lineFlag = flag.Bool("l", false, "use LINE Notify")
var streamFlag = flag.Bool("s", false, "show User Stream")
var jst *time.Location
var twitter *anaconda.TwitterApi
var line *LineApi
var words []string
var messageTemplate *template.Template

func main() {
	flag.Parse()
	setup()
	readStream()
}

func setup() {
	setupJST()
	setupTwitter()
	setupWords()
	if *lineFlag == true {
		setupLine()
	}
	setupMessageTemplate()
}

func setupJST() {
	l, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	jst = l
}

func setupTwitter() {
	anaconda.SetConsumerKey(os.Getenv("CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("CONSUMER_SECRET"))
	twitter = anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
}

func setupWords() {
	f, err := os.Open("words")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
}

func setupLine() {
	line = NewLineApi(os.Getenv("LINE_AUTHORIZATION"))
}

func setupMessageTemplate() {
	messageTemplate = template.Must(template.ParseFiles(filepath.Join("templates", "line")))
}

func readStream() {
	f, err := os.OpenFile(filepath.Join("tweets.csv"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)

	stream := twitter.UserStream(nil)
	for {
		x := <-stream.C
		switch x := x.(type) {
		case anaconda.Tweet:
			if x.RetweetedStatus == nil {
				w.Write(record(x))

				fmt.Printf("@%v / %v\n%v\n%v\n", x.User.ScreenName, x.User.Name,
					x.FullText, x.Entities)

				if *lineFlag == true && inspectTweet(x, words) == true {
					s := templateWrite(x)
					lineNotify(s)
				}
			} else {
				w.Write(record(*x.RetweetedStatus))

				fmt.Printf("@%v / %v\nRT @%v %v\n%v\n",
					x.User.ScreenName, x.User.Name, x.RetweetedStatus.User.ScreenName,
					x.RetweetedStatus.FullText, x.RetweetedStatus.Entities)
			}
		default:
		}
		w.Flush()
	}
}

func record(tweet anaconda.Tweet) []string {
	var r []string
	t, err := tweet.CreatedAtTime()
	if err != nil {
		log.Fatal(err)
	}
	jt := t.In(jst)
	r = append(r, jt.String())
	r = append(r, tweet.User.ScreenName)
	s := strings.Replace(tweet.FullText, "\n", " ", -1)
	r = append(r, s)
	return r
}

func inspectTweet(tweet anaconda.Tweet, words []string) bool {
	t := tweet.FullText
	for _, w := range words {
		if strings.Contains(t, w) {
			return true
		}
	}
	return false
}

func templateWrite(tweet anaconda.Tweet) string {
	var b bytes.Buffer
	messageTemplate.Execute(&b, tweet)
	return b.String()
}

func lineNotify(message string) {
	r, err := line.sendNotify(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("lineNotify: ", r)
}
