package fs_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/kislerdm/github-repo-details/internal/fs"
)

func TestFSOps(t *testing.T) {
	data := []byte(`foo-bar`)
	filePath := "/tmp/test-fs-pkc.txt"

	err := fs.FileWrite(data, filePath)
	if err != nil {
		t.Fatal(err)
	}

	got := fs.FileRead(filePath)
	if !reflect.DeepEqual(got, data) {
		t.Errorf("\nwant: %+v,\ngot: %+v", data, got)
	}

	t.Cleanup(func() {
		os.Remove(filePath)
	})
}
