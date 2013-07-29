package quality

import "fmt"

// Evaluator calcuates a quality metric given search results and a corpus.
type Evaluator func(results Results, corpus Corpus) (float64, error)

// Evaluators is a set of evaluators keyed by their metric name.
type Evaluators map[string]Evaluator

// MRR is the Mean Reciprocal Rank; that is, the inverse harmonic mean the
// positions of first relevant result across all information needs.
// Find the index of the first relevant document. Use reciprocal no relevant
// document? Use 0. Then average across all information needs.
func MRR(results Results, corpus Corpus) (float64, error) {
	count, acc := 0, 0.0
	for query, docs := range results {
		rrank := 0.0
		for i, doc := range docs {
			relevant, err := corpus.Relevant(query, doc)
			if err != nil {
				return 0.0, fmt.Errorf("could not test relevance: %s", err.Error())
			}

			if relevant {
				rrank = (1.0 / float64(i+1))
				break
			}
		}

		acc = acc + rrank
		count = count + 1
	}

	return acc / float64(count), nil
}

// MAP is the mean average precision. For a single information need, Average
// Precision is the average of the precision value obtained for the set of top
// k documents existing after each relevant document is retrieved, and this
// value is then averaged over information needs.
func MAP(results Results, corpus Corpus) (float64, error) {
	count, acc := 0, 0.0
	for query, docs := range results {
		tp, fp := 0, 0
		averagePrecision := 0.0
		averagePrecisionCount, averagePrecisionAcc := 0, 0.0
		for _, doc := range docs {
			relevant, err := corpus.Relevant(query, doc)
			if err != nil {
				return 0.0, fmt.Errorf("could not test relevance: %s", err.Error())
			}

			if relevant {
				tp = tp + 1
				precision := float64(tp) / float64(tp+fp)

				// set up to calulate average precision for this query
				averagePrecisionAcc = averagePrecisionAcc + precision
				averagePrecisionCount = averagePrecisionCount + 1
			} else {
				fp = fp + 1
			}
		}

		// set up to calculate mean of average precisions. count ap as 0 if
		// no relevant documents were found at all.
		if averagePrecisionCount > 0 {
			averagePrecision = averagePrecisionAcc / float64(averagePrecisionCount)
		}

		acc = acc + averagePrecision
		count = count + 1
	}

	return acc / float64(count), nil
}

// PrecisionAtK returns an evaluator function for the given K.
func PrecisionAtK(K int) Evaluator {
	return func(results Results, corpus Corpus) (float64, error) {
		count, acc := 0, 0.0
		tp, fp := 0, 0
		precision := 0.0
		for query, docs := range results {
			for i, doc := range docs {
				if i >= K {
					break
				}

				relevant, err := corpus.Relevant(query, doc)
				if err != nil {
					return 0.0, fmt.Errorf("could not test relevance: %s", err.Error())
				}

				if relevant {
					tp = tp + 1
				} else {
					fp = fp + 1
				}
			}

			precision = float64(tp) / float64(tp+fp)
			acc = acc + precision
			count = count + 1
		}

		return acc / float64(count), nil
	}
}
