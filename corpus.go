package quality

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// Label represents a single relevance judgement for a document. Labels can
// only be used in connection with a query (see Corpus).
type Label struct {
	Doc      Doc
	Relevant bool
}

// Corpus connects relevance judgements (Label) to queries.
type Corpus map[InformationNeed][]Label

// Relevant checks if a document is relevant to a query in the receiving corpus.
// Will fail if the query is not in the corpus; we should never be crawling
// queries that are not in the corpus.
func (corpus Corpus) Relevant(query InformationNeed, doc Doc) (bool, error) {
	labels, ok := corpus[query]
	if !ok {
		return false, fmt.Errorf("query not in corpus: %s", query)
	}

	for _, label := range labels {
		if label.Doc.ID == doc.ID {
			return label.Relevant, nil
		}
	}

	return false, nil
}

// ReadTsvCorpus reads the corpus from disk into a Corpus struct
func ReadTsvCorpus(filename string) (Corpus, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Corpus{}, fmt.Errorf("need a valid input file: %v", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		return Corpus{}, fmt.Errorf("could not read csv: %s ", err)
	}

	corpus := Corpus{}
	for i, record := range records {
		if l := len(record); l != 5 {
			return Corpus{}, fmt.Errorf("record is not %v long: %v", l, record)
		}

		id, err := strconv.Atoi(record[3])
		if err != nil {
			return Corpus{}, fmt.Errorf("bad id on line %n: %s", i, err)
		}

		relevance, err := strconv.ParseBool(record[4])
		if err != nil {
			return Corpus{}, fmt.Errorf("bad relevance on line %v: %s", i, err)
		}

		query := InformationNeed{
			Kind:  record[0],
			Query: record[1],
		}

		label := Label{
			Doc: Doc{
				Kind: record[2],
				ID:   id,
			},
			Relevant: relevance,
		}

		corpus[query] = append(corpus[query], label)
	}

	return corpus, nil
}
