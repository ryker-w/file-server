package filesystem

import (
	"testing"
)

func TestName(t *testing.T) {
	resp, err := listDir(`c:\`)
	if err != nil {
		t.Fatal(err)
	}
	for _, folder := range resp.Folders {
		t.Logf("%s", folder)
	}
	for _, file := range resp.Files {
		t.Logf("%s[%s]", file.Name, file.Ext)
	}

}
