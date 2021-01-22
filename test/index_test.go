package test

import (
	"testing"

	"github.com/Kua-Fu/gsearch/core"
)

func TestDoc(t *testing.T) {
	var (
		err      error
		fPtr     *core.File
		indexDir string
	)

	indexDir = "/Users/yz/work/github/gsearch/test/index/"

	fPtr, err = core.CreateFile(indexDir, true, true)
	if err != nil {
		t.Error(err)
	}

	analyzer := new(core.Analyzer)
	writer := new(core.DocumentWriter)
	writer.Init(fPtr.FilePath, *analyzer, int64(1000))

	doc := new(core.Document)
	f1, _ := core.Keyword("path", "/etc/test.txt")
	doc.Add(f1)

	segment := "s1"
	writer.AddDocument(segment, *doc)

}
