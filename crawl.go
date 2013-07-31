package quality

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

// Results from an api crawl. Holds all queries and their returned ordered []Doc.
type Results map[InformationNeed][]Doc

// CrawlMessage is Either([]Docs, Error). Exclusively for egress chan.
type CrawlMessage struct {
	IN    InformationNeed
	Docs  []Doc
	Error error
}

// crawl does a concurrent crawl of the api.
func Crawl(api API, corpus Corpus, concurrency, retries int) (Results, error) {
	results := Results{}
	ingress := make(chan InformationNeed)
	egress := make(chan CrawlMessage)

	go func() {
		for in := range corpus {
			ingress <- in
		}

		close(ingress)
	}()

	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for in := range ingress {
				var docs []Doc
				var err error

				for i := 0; i <= retries; i++ {
					docs, err = search(api, in)
					if err == nil {
						break
					}

					sleep := math.Pow(2.5, float64(i)) - 1
					time.Sleep(time.Second * time.Duration(sleep))
					log.Printf(err.Error())
				}

				fmt.Print(".")
				egress <- CrawlMessage{
					IN:    in,
					Docs:  docs,
					Error: err,
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(egress)
	}()

	for msg := range egress {
		need := msg.IN
		if msg.Error != nil {
			return Results{}, fmt.Errorf("%v '%v': %v", need.Kind, need.Query, msg.Error)
		}

		results[need] = msg.Docs
	}

	return results, nil
}
