package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	colourize = false
)

func main() {
	//Get query
	query := strings.Join(os.Args[1:], " ")

	query = parseColourizeFlag(query)
	//Perform request to Google searching for a query in StackOverflow
	questions, err := performSearch(query)
	if err != nil {
		log.Fatal("sorry, I couldn't find what you're looking for :(")
	}
	//Get the first link of the response. Perform a request to that link
	answer := performRequest(questions[0])
	//Display the answer propperly (http://misc.flogisoft.com/bash/tip_colors_and_formatting)
	printAnswer(answer, questions[0])
}

func printAnswer(answer, url string) {
	if colourize {
		answer, url = colourizeOutput(answer, url)
	}
	fmt.Println(answer)
	fmt.Println(url)
}

func colourizeOutput(answer, url string) (string, string) {
	for _, v := range os.Args[1:] {
		answer = strings.Replace(answer, v, "\033[32m" + v + "\033[39m", -1)
		if(len(v) >= 3){
			answer = strings.Replace(answer, strings.Title(v), "\033[32m" + strings.Title(v) + "\033[39m", -1)
			answer = strings.Replace(answer, strings.ToUpper(v), "\033[32m" + strings.ToUpper(v) + "\033[39m", -1)
		}
	}
	url = "\033[32mUrl: \033[33m" + url + "\033[39m"
	return answer, url
}

func parseColourizeFlag(query string) string {
	if strings.Index(query, "-c") >= 0 {
		colourize = true
		query = strings.Replace(query, "-c", "", 1)
	}
	return query
}

func performRequest(url string) string {
	doc, err := goquery.NewDocument(url + "?answertab=votes")
	if err != nil {
		log.Fatal(err)
	}
	return doc.Find(".answercell .post-text").First().Text()
}

func performSearch(query string) ([]string, error) {
	searchURL := "http://www.google.com/search?q=site:stackoverflow.com/questions%20" + url.QueryEscape(query)

	/* userAgents := [...]string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; rv:11.0) Gecko/20100101 Firefox/11.0",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:22.0) Gecko/20100 101 Firefox/22.0",
		"Mozilla/5.0 (Windows NT 6.1; rv:11.0) Gecko/20100101 Firefox/11.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_4) AppleWebKit/536.5 (KHTML, like Gecko) ' 'Chrome/19.0.1084.46 Safari/536.5",
		"Mozilla/5.0 (Windows; Windows NT 6.1) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.46' 'Safari/536.5",
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		log.Fatal("error performing request")
	}
	req.Header.Add("User-Agent", userAgents[rand.Intn(len(userAgents))])
	resp, err := client.Do(req) */
	resp, err := http.Get(searchURL)
	if err != nil {
		log.Fatal(err)
	}

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

	if len(links) == 0 {
		return nil, errors.New("search failed")
	}
	return links, nil
}

// returns true if the link match with the regex. This indicates that it's a real question
// and not a tagged question.
func isQuestion(link string) bool {
	r := regexp.MustCompile(`questions/\d+/`)
	return r.Find([]byte(link)) != nil
}
