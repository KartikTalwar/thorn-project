package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "backpage.com"

var (
	phoneNumberRegex = regexp.MustCompile(`(\d{3})\D*(\d{3})\D*(\d{4})`)
	ageRegex         = regexp.MustCompile(`[Aa][Gg][Ee].?\s*\d\d?`)
)

type result struct {
	err      error
	duration time.Duration
	listing  *Listing
}

type Scraper struct {
	WorkerCount       int
	RequestsPerSecond int

	results chan *result
}

func (s *Scraper) Run() {
	s.results = make(chan *result)

	var throttle <-chan time.Time
	if s.RequestsPerSecond > 0 {
		throttle =
			time.Tick(time.Duration(1e6/s.RequestsPerSecond) * time.Microsecond)
	}

	jobs := make(chan *http.Request)

	// Start workers
	for i := 0; i < s.WorkerCount; i++ {
		go s.worker(jobs)
	}

	// Process results
	go func() {
		for r := range s.results {
			if r.err != nil {
				fmt.Println(r.err)
			}

			json.NewEncoder(os.Stdout).Encode(r.listing)
		}
	}()

	// Start sending jobs to the workers
	for {
		if s.RequestsPerSecond > 0 {
			<-throttle
		}

		doc, err := goquery.NewDocument("http://sf.backpage.com/FemaleEscorts/")
		if err != nil {
			fmt.Println(err)
			continue
		}

		doc.Find(".cat").Each(func(i int, s *goquery.Selection) {
			link, _ := s.Find("a").Attr("href")

			req, err := http.NewRequest("GET", link, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			jobs <- req
		})
	}
	close(jobs)
}

func (s *Scraper) worker(c chan *http.Request) {
	for req := range c {
		start := time.Now()

		listing, err := parseListing(req.URL)
		if err != nil {
			s.results <- &result{err: err}
			continue
		}

		s.results <- &result{
			duration: time.Now().Sub(start),
			listing:  listing,
		}
	}
}

func parseListing(u *url.URL) (*Listing, error) {
	escort := Escort{}
	listing := Listing{URL: u, Escort: &escort}

	doc, err := goquery.NewDocument(u.String())
	if err != nil {
		return nil, err
	}

	listing.Title = doc.Find("#postingTitle .h1link").Text()
	listing.Description = strings.TrimSpace(doc.Find(".postingBody").Text())

	escort.PhoneNumber = phoneNumberRegex.FindString(listing.Description)
	escort.Age = ageRegex.FindString(listing.Description)

	return &listing, nil
}

type Listing struct {
	URL         *url.URL
	Title       string
	Description string
	Escort      *Escort
}

type Escort struct {
	Age         string
	PhoneNumber string
}

func main() {
	scraper := Scraper{
		WorkerCount: 8,
	}

	scraper.Run()
}
