package quality

import (
	"flag"
	"fmt"
	"github.com/peterbourgon/g2s"
	"log"
)

// InformationNeed is the Information Need the user wishes to satisfy.
type InformationNeed struct {
	Kind  string
	Query string
}

// Doc is a query result. Kind can be used by a universal search returning
// multiple types of objects.
type Doc struct {
	Kind string
	ID   int
}

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

	corpus, err := readTsvCorpus(*inputFile)
	if err != nil {
		log.Fatalf("could not read corpus: %v", err)
	}

	decoder := func(json []byte) ([]Doc, error) {
		return []Doc{}, nil
	}

	api := HTTPAPI{
		Site:    *searchSite,
		Path:    *searchPath,
		Decoder: decoder,
	}

	results, err := crawl(api, corpus, *concurrency, *retries)
	if err != nil {
		log.Fatalf("crawl aborted: %v", err)
	}

	fmt.Println()
	log.Printf("completed crawl with %v queries", len(results))

	es := Evaluators{
		"mrr":            MRR,
		"map":            MAP,
		"precision-at-1": PrecisionAtK(1),
		"precision-at-2": PrecisionAtK(2),
		"precision-at-3": PrecisionAtK(3),
		"precision-at-4": PrecisionAtK(4),
		"precision-at-5": PrecisionAtK(5),
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
