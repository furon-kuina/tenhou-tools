package cmd

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(downloadCmd)
}

type Client struct {
	baseUrl string
	c       *http.Client
}

func NewClient() *Client {
	return &Client{
		baseUrl: "https://tenhou.net/sc/raw/",
		c:       http.DefaultClient,
	}
}

func (c *Client) GetFileIndex() (*http.Response, error) {
	return c.c.Get(c.baseUrl + "list.cgi")
}

// GetDownloadList fetches FileIndex and compares size with the files in
// the specified directory. It returns files 1. that don't exist in the
// directory, or 2. whose size don't match the size in the directory
func (c *Client) GetDownloadList(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir(%q): %v", dir, err)
	}
	sizeMap := make(map[string]int)
	for _, e := range entries {
		if !e.Type().IsRegular() {
			continue
		}
		fi, err := e.Info()
		if err != nil {
			return nil, fmt.Errorf("%s info: %v", e.Name(), err)
		}
		sizeMap[e.Name()] = int(fi.Size())
	}
	resp, err := c.GetFileIndex()
	if err != nil {
		return nil, fmt.Errorf("GetFileIndex: %v", err)
	}
	files, err := parseFileIndex(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parseFileIndex: %v", err)
	}
	res := make([]string, 0)
	for _, logFile := range files {
		size, ok := sizeMap[logFile.FileName]
		if !ok {
			res = append(res, logFile.FileName)
			continue
		}
		if size != logFile.Size {
			res = append(res, logFile.FileName)
		}
	}
	return res, nil
}

func (c *Client) GetLogFile(fileName string) (io.Reader, error) {
	resp, err := c.c.Get(c.baseUrl + "dat/" + fileName)
	if err != nil {
		return nil, err
	}
	return gzip.NewReader(resp.Body)
}

var downloadCmd = &cobra.Command{
	Use:   "tenhou-tool download path",
	Short: "tenhou-tool download downloads Tenhou history file listed in a file in the given path",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires 1 argument")
		}
		_, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("could not open %s: %w", args[0], err)
		}

		return nil
	},
}

// DownloadLogs donwloads log files listed in the io.Reader.
// Refer to https://tenhou.net/sc/raw/ for the details.

type LogFile struct {
	FileName string
	Size     int
}

func parseFileIndex(r io.Reader) ([]*LogFile, error) {
	sc := bufio.NewScanner(r)
	logFiles := make([]*LogFile, 0)
	for sc.Scan() {
		line := sc.Text()
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		fileName := strings.Split(line, "'")[1]
		tmp1 := strings.Split(line, ":")
		tmp2 := tmp1[len(tmp1)-1]
		sizeStr := tmp2[0 : len(tmp2)-2]
		// no trailing comma
		if line[len(line)-1] == '}' {
			sizeStr = tmp2[0 : len(tmp2)-1]
		}
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return nil, fmt.Errorf("strconv.Atoi(%q): %w", sizeStr, err)
		}
		logFile := &LogFile{
			FileName: fileName,
			Size:     size,
		}
		logFiles = append(logFiles, logFile)
	}
	return logFiles, nil
}
