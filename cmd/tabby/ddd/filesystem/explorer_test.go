package filesystem

import (
	"encoding/json"
	"testing"
)

func TestName(t *testing.T) {
	var req = ExplorerReq{Path: []string{"/"}}
	bs, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}
