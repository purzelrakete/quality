package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/peterbourgon/g2s"
	"github.com/purzelrakete/quality"
	"log"
)

var (
	concurrency    = flag.Int("concurrency", 2, "concurrency")
	inputFile      = flag.String("inputFile", "input.tsv", "labelled results")
	retries        = flag.Int("retries", 0, "retries before failing")
	searchPath     = flag.String("searchPath", "/search/%s?q=%s", "path")
	searchSite     = flag.String("searchSite", "http://localhost", "site")
	statsDEndpoint = flag.String("statsDEndpoint", "stats:8125", "statsD endpoint")
	statsDNS       = flag.String("statsDNS", "search.quality", "statsD ns")
)

func init() {
	flag.Parse()
}

func main() {
	s, err := g2s.Dial("udp", *statsDEndpoint)
	if err != nil {
		log.Fatalf("no transport to statsD endpoint %v: %v", *statsDEndpoint, err)
	}

	corpus, err := quality.ReadTsvCorpus(*inputFile)
	if err != nil {
		log.Fatalf("could not read corpus: %v", err)
	}

	type HTTPAPIResponse struct {
		Docs []quality.Doc `json:"docs"`
	}

	decoder := func(jsn []byte) ([]quality.Doc, error) {
		var response HTTPAPIResponse
		if err := json.Unmarshal(jsn, &response); err != nil {
			return []quality.Doc{}, fmt.Errorf("could not converted from json: %s")
		}

		return response.Docs, nil
	}

	api := quality.HTTPAPI{
		Site:    *searchSite,
		Path:    *searchPath,
		Decoder: decoder,
	}

	results, err := quality.Crawl(api, corpus, *concurrency, *retries)
	if err != nil {
		log.Fatalf("crawl aborted: %v", err)
	}

	fmt.Println()
	log.Printf("completed crawl with %v queries", len(results))

	es := quality.Evaluators{
		"mrr":            quality.MRR,
		"map":            quality.MAP,
		"precision-at-1": quality.PrecisionAtK(1),
		"precision-at-2": quality.PrecisionAtK(2),
		"precision-at-3": quality.PrecisionAtK(3),
		"precision-at-4": quality.PrecisionAtK(4),
		"precision-at-5": quality.PrecisionAtK(5),
	}

	for name, e := range es {
		summary, err := e(results, corpus)
		if err != nil {
			log.Fatalf("could not evaluate %v: %v", name, err)
		}

		key := fmt.Sprintf("%v.%v", *statsDNS, name)
		value := fmt.Sprintf("%.5f", summary)

		log.Printf("%v: %v", key, value)
		s.Gauge(1.0, key, value)
	}

	log.Printf("done.")
}
