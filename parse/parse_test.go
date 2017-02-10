package parse

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseDir(t *testing.T) {
	p := New(true)

	if err := p.ParseDir("../testdata"); err != nil {
		t.Log("error:", err)
	}
	bs, _ := json.MarshalIndent(p.Mappings, "", "\t")
	fmt.Println(string(bs))
}
