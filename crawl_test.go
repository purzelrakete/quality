package quality

import "testing"

func TestCrawl(t *testing.T) {
	corpus := testCorpus()
	api := MockAPI{}

	results, err := crawl(api, corpus, 1, 0)
	if err != nil {
		t.Fatalf("could not crawl: %v", err)
	}

	expectedSearches := 1
	if got := len(results); got != expectedSearches {
		t.Fatalf("crawled %v instead of %v queries", got, expectedSearches)
	}

	expectedDoc := Doc{
		Kind: "person",
		ID:   1,
	}

	docs, ok := results[testIN("deadmau5")]

	if !ok {
		t.Fatalf("did not crawl 'deadmau5'")
	}

	expectedDocuments := 1
	if got := len(docs); got != expectedDocuments {
		t.Fatalf("%v 'deadmau5' results instead of %v", got, expectedDocuments)
	}

	if got := docs[0]; got != expectedDoc {
		t.Fatalf("got %v instead of %v", got, expectedDoc)
	}
}

func TestFailedRetriesCrawl(t *testing.T) {
	corpus := testCorpus()
	api := FailAPI{}

	results, err := crawl(api, corpus, 1, 1)

	expectedSearches := 0
	if got := len(results); got != expectedSearches {
		t.Fatalf("crawled %v instead of %v queries", got, expectedSearches)
	}

	if err == nil {
		t.Fatalf("expected error but got no error")
	}
}
