# Quality

[![Build Status](https://travis-ci.org/purzelrakete/quality.png)](https://travis-ci.org/purzelrakete/quality)

Ranked retrieval quality metrics, using a labelled corpus.

## Corpus

Corpuus is stored as a line oriented TSV. Each individual quality judgement
takes on the following form:

    intent kind, query, entity kind, entity id, relevance

where relevance is in {0, 1}.

    person	"""deadma5"""	person	11	0
    person	skrillex	person	12	1
    person	asura	person	13	1
    person	asura	person	14	1
    person	bach	person	15	0
    person	bach	person	16	0
    person	bach	person	17	1

## Metrics

The following metrics are currently implemented.

### Mean Reciprocal Rank

MRR is the Mean Reciprocal Rank; that is, the inverse harmonic mean the
positions of first relevant result across all information needs.

### Mean Average Precision

MAP is the mean average precision. For a single information need, Average
Precision is the average of the precision value obtained for the set of top
k documents existing after each relevant document is retrieved, and this value
s then averaged over information needs.

### Precision at K

Precision at K is the precision (true positives / (true positives + false
positives) of all search results up to position K.
