package domo

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestCreatePageCollectionWillSerializeCorrectly(t *testing.T) {
	page := PageCollection{Title: "test title", Description: "test description"}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(page)
	if err != nil {
		t.Error("Error in JSON encoding")
	}
	actual := buf.String()
	expected := "{\"title\":\"test title\",\"description\":\"test description\"}\n"
	if actual != expected {
		t.Errorf("Json didn't serialize as expected \n\nactual:  \t%s\nexpected:\t%s", actual, expected)
	}
}
