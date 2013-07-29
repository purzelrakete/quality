package quality

import "fmt"

type MockAPI struct{}

func (api MockAPI) Search(query InformationNeed) ([]Doc, error) {
	return []Doc{testDoc(1)}, nil
}

type FailAPI struct{}

func (api FailAPI) Search(query InformationNeed) ([]Doc, error) {
	return []Doc{}, fmt.Errorf("FailAPI is fail.")
}

func testCorpus() Corpus {
	return Corpus{
		testIN("deadmau5"): []Label{
			Label{
				Doc:      testDoc(1),
				Relevant: true,
			},
		},
	}
}

func testIN(name string) InformationNeed {
	return InformationNeed{
		Kind:  "person",
		Query: name,
	}
}

func testDoc(id int) Doc {
	return Doc{
		Kind: "person",
		ID:   id,
	}
}
