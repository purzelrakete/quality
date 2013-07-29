package quality

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// API is a thing that can return []Docs for a given information need.
type API interface {
	Search(query InformationNeed) ([]Doc, error)
}

// HTTPAPI concrete Api implementation.
type HTTPAPI struct {
	Site    string
	Path    string
	Decoder HTTPAPIDecode
}

// HTTPAPIDecode
type HTTPAPIDecode func(json []byte) ([]Doc, error)

// Search an api with a given query.
func (api HTTPAPI) Search(q InformationNeed) ([]Doc, error) {
	url := api.Site + fmt.Sprintf(api.Path, q.Kind, url.QueryEscape(q.Query))

	start := time.Now()
	resp, err := http.Get(url)
	ms := time.Since(start)
	if err != nil {
		return []Doc{}, fmt.Errorf("query (%v, %v) failed: %v ", url, ms, err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Doc{}, fmt.Errorf("could not read body: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return []Doc{}, fmt.Errorf("%s (%v, %v): %s", url, resp.Status, ms, body)
	}

	docs, err := api.Decoder(body)
	if err != nil {
		return []Doc{}, fmt.Errorf("could not decode: %s")
	}

	return docs, nil
}

// search delegates to api.Search.
func search(api API, q InformationNeed) ([]Doc, error) {
	docs, err := api.Search(q)
	if err != nil {
		return []Doc{}, fmt.Errorf("failed to search api: %s", err.Error())
	}

	return docs, err
}
