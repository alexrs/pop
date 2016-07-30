package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	//Get query
	query := strings.Join(os.Args[1:], " ")
	//Perform request to Google searching for a query in StackOverflow
	questions := performSearch(query)
	//Get the first link of the response. Perform a request to that link
	answer := performRequest(questions[0])
	//Display the answer propperly (http://misc.flogisoft.com/bash/tip_colors_and_formatting)
	printAnswer(answer, questions[0])
}

func printAnswer(answer, url string) {
	fmt.Println(answer)
	fmt.Println("Url:", url)
}

func performRequest(url string) string {
	doc, err := goquery.NewDocument(url + "?answertab=votes")
	if err != nil {
		log.Fatal(err)
	}
	return doc.Find(".answercell .post-text").First().Text()
}

func performSearch(query string) []string {
	userAgents := [...]string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; rv:11.0) Gecko/20100101 Firefox/11.0",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:22.0) Gecko/20100 101 Firefox/22.0",
		"Mozilla/5.0 (Windows NT 6.1; rv:11.0) Gecko/20100101 Firefox/11.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_4) AppleWebKit/536.5 (KHTML, like Gecko) ' 'Chrome/19.0.1084.46 Safari/536.5",
		"Mozilla/5.0 (Windows; Windows NT 6.1) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.46' 'Safari/536.5",
	}

	searchURL := "http://www.google.com/search?q=site:stackoverflow.com/questions%20" + url.QueryEscape(query)
	client := &http.Client{}
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		log.Fatal("error performing request")
	}
	req.Header.Add("User-Agent", userAgents[rand.Intn(len(userAgents))])
	resp, err := client.Do(req)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(io.Reader(resp.Body))
	if err != nil {
		log.Fatal("error reading document", err)
	}
	var links []string
	doc.Find("h3.r a").Each(func(i int, s *goquery.Selection) {
		str, exists := s.Attr("href")
		if exists {
			u, err := url.Parse(str)
			if err != nil {
				log.Fatal(err)
			}
			m, _ := url.ParseQuery(u.RawQuery)
			link := m["q"][0]
			if isQuestion(link) {
				links = append(links, link)
			}
		}
	})
	return links
}

// returns true if the link match with the regex. This indicates that it's a real question
// and not a tagged question.
func isQuestion(link string) bool {
	r := regexp.MustCompile(`questions/\d+/`)
	return r.Find([]byte(link)) != nil
}
