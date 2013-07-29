package quality

import "testing"

func TestSearch(t *testing.T) {
	corpus := testCorpus()
	api := MockAPI{}

	for query := range corpus {
		docs, err := search(api, query)
		if err != nil {
			t.Fatalf("search for '%v' failed: %v", query, err)
		}

		if got := len(docs); got != 1 {
			t.Fatalf("expected one result but got %v", got)
		}
	}
}
