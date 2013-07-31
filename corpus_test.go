package quality

import "testing"

func TestTsvCorpusLength(t *testing.T) {
	corpus, err := ReadTsvCorpus("labels.tsv")
	if err != nil {
		t.Fatalf("could not read from tsv repo: %s", err)
	}

	expected := 4
	if got := len(corpus); got != expected {
		t.Fatalf("expected %v docs but got %v", expected, got)
	}
}

func TestTsvRepoLabelGrouping(t *testing.T) {
	corpus, err := ReadTsvCorpus("labels.tsv")
	if err != nil {
		t.Fatalf("could not read from tsv repo: %s", err)
	}

	expected := []Label{
		Label{
			Doc:      testDoc(15),
			Relevant: false,
		},
		Label{
			Doc:      testDoc(16),
			Relevant: false,
		},
		Label{
			Doc:      testDoc(17),
			Relevant: true,
		},
	}

	if got, ok := corpus[testIN("bach")]; ok {
		if l := len(got); l != 3 {
			t.Fatalf("did not get 3 bach labels")
		}

		for i, label := range expected {
			if got[i] != label {
				t.Fatalf("expected %v to equal %v.", label, got[i])
			}
		}
	} else {
		t.Fatalf("expected 'bach' labels but found none.")
	}
}

func TestRelevance(t *testing.T) {
	corpus := testCorpus()
	doc := testDoc(1)

	relevant, err := corpus.Relevant(testIN("deadmau5"), doc)
	if err != nil {
		t.Fatalf("error testing relevance: %v", err)
	}

	if !relevant {
		t.Fatalf("document should have been relevant but was not: %v", doc)
	}
}
