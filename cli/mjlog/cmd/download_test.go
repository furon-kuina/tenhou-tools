package cmd

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFetchFileInfo(t *testing.T) {
	c := NewClient()
	resp, err := c.GetFileIndex()
	if err != nil {
		t.Fatalf("GetFileIndex: %v", err)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(bytes) == 0 {
		t.Errorf("empty response")
	}
}

func TestParseFileIndex(t *testing.T) {
	give := `list([
{file:'sca20240715.log.gz',size:69273},
{file:'scb2024070600.log.gz',size:64226},
{file:'scf20240715.html.gz',size:8695}
]);`
	want := []*LogFile{
		{FileName: "sca20240715.log.gz", Size: 69273},
		{FileName: "scb2024070600.log.gz", Size: 64226},
		{FileName: "scf20240715.html.gz", Size: 8695},
	}
	r := strings.NewReader(give)
	res, err := parseFileIndex(r)
	if err != nil {
		t.Fatalf("ParseFileIndex: %v", err)
	}
	if !cmp.Equal(res, want) {
		t.Errorf("diff:\n %s", cmp.Diff(res, want))
	}
}

func TestGetLogFile(t *testing.T) {
	give := `list([
{file:'scb2024071810.log.gz',size:64226},
]);`
	c := NewClient()
	logFiles, err := parseFileIndex(strings.NewReader(give))
	if err != nil {
		t.Fatalf("parseFileIndex: %v", err)
	}
	for _, logFile := range logFiles {
		r, err := c.GetLogFile(logFile.FileName)
		if err != nil {
			t.Errorf("GetLogFile(%q): %v", logFile.FileName, err)
		}
		respByte, err := io.ReadAll(r)
		if err != nil {
			t.Errorf("reading body: %v", err)
		}
		fmt.Println(string(respByte))
	}

}
