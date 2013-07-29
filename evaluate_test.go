package quality

import "testing"

// Mean Reciprocal Rank

func TestMRR(t *testing.T) {
	corpus := Corpus{
		testIN("deadmau5"): []Label{
			Label{
				Doc:      testDoc(1),
				Relevant: true,
			},
		},
		testIN("beethoven"): []Label{
			Label{
				Doc:      testDoc(7),
				Relevant: true,
			},
		},
	}

	results := Results{
		testIN("deadmau5"): []Doc{
			testDoc(9),
			testDoc(1),
		},
		testIN("beethoven"): []Doc{
			testDoc(9),
			testDoc(1),
			testDoc(7),
		},
	}

	got, err := MRR(results, corpus)
	if err != nil {
		t.Fatalf("could evaluate mrr: %v", err)
	}

	expected := 0.4166

	if int32(got*1000) != int32(expected*1000) {
		t.Fatalf("expected %v to equal %v.", expected, got)
	}
}

func TestPerfectMRR(t *testing.T) {
	corpus := testCorpus()
	results := Results{
		testIN("deadmau5"): []Doc{
			testDoc(1),
		},
	}

	got, err := MRR(results, corpus)
	if err != nil {
		t.Fatalf("could evaluate mrr: %v", err)
	}

	expected := 1.0

	if got != expected {
		t.Fatalf("expected %v to equal %v.", expected, got)
	}
}

func TestNoRankMRR(t *testing.T) {
	corpus := testCorpus()
	results := Results{
		testIN("deadmau5"): []Doc{},
	}

	got, err := MRR(results, corpus)
	if err != nil {
		t.Fatalf("could evaluate mrr: %v", err)
	}

	expected := 0.0

	if got != expected {
		t.Fatalf("expected %v to equal %v.", expected, got)
	}
}

// MAP

func TestMAP(t *testing.T) {
	corpus := Corpus{
		testIN("deadmau5"): []Label{
			Label{
				Doc:      testDoc(2),
				Relevant: true,
			},
			Label{
				Doc:      testDoc(4),
				Relevant: true,
			},
		},
		testIN("beethoven"): []Label{
			Label{
				Doc:      testDoc(5),
				Relevant: true,
			},
		},
		testIN("can"): []Label{
			Label{
				Doc:      testDoc(100),
				Relevant: true,
			},
		},
	}

	results := Results{
		testIN("deadmau5"): []Doc{
			testDoc(1),
			testDoc(2),
			testDoc(3),
			testDoc(4),
		},
		testIN("beethoven"): []Doc{
			testDoc(5),
			testDoc(6),
			testDoc(7),
		},
		testIN("can"): []Doc{},
	}

	got, err := MAP(results, corpus)
	if err != nil {
		t.Fatalf("could evaluate map: %v", err)
	}

	expected := 0.5

	if int32(got*1000) != int32(expected*1000) {
		t.Fatalf("expected %v to equal %v.", expected, got)
	}
}

// Precision@K

func TestPrecisionAtK(t *testing.T) {
	corpus := Corpus{
		testIN("deadmau5"): []Label{
			Label{
				Doc:      testDoc(2),
				Relevant: true,
			},
			Label{
				Doc:      testDoc(4),
				Relevant: true,
			},
		},
		testIN("beethoven"): []Label{
			Label{
				Doc:      testDoc(5),
				Relevant: true,
			},
		},
		testIN("can"): []Label{
			Label{
				Doc:      testDoc(100),
				Relevant: true,
			},
		},
	}

	results := Results{
		testIN("deadmau5"): []Doc{
			testDoc(1),
			testDoc(2),
			testDoc(3),
			testDoc(4),
		},
		testIN("beethoven"): []Doc{
			testDoc(5),
			testDoc(6),
			testDoc(7),
		},
		testIN("can"): []Doc{},
	}

	// test
	got, err := PrecisionAtK(2)(results, corpus)
	if err != nil {
		t.Fatalf("could evaluate map: %v", err)
	}

	expected := 0.5

	if int32(got*1000) != int32(expected*1000) {
		t.Fatalf("expected %v to equal %v.", expected, got)
	}
}
