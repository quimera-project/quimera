package fs

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/justincormack/go-memfd"
	"github.com/quimera-project/quimera/internal/utils/live"
)

func init() {
	Assets = &assets{}
	Checks = &checks{}
	Lang = &lang{}
	Templates = &templates{}
	Tools = &tools{}
}

type Memfd struct {
	Fd   *memfd.Memfd
	Name string
}

func Elf(buf []byte) bool {
	return len(buf) > 52 &&
		buf[0] == 0x7F && buf[1] == 0x45 &&
		buf[2] == 0x4C && buf[3] == 0x46
}

// NewMemfdAsset returns a new Memfd asset struct pointer.
func NewMemfdAsset(fd *memfd.Memfd, asset string) *Memfd {
	return &Memfd{Fd: fd, Name: strings.ToUpper(strings.Split(asset, ".")[0])}
}

// Anonymous create an anonymous file using memfd_create and writes the supplied data into the file.
//
// Returns the Memfd structure and any error encountered.
func Anonymous(data []byte) (*memfd.Memfd, error) {
	mfd, err := memfd.CreateNameFlags("", memfd.AllowSealing)
	if err != nil {
		return nil, fmt.Errorf("creating anonymous file: %v", err)
	}
	if _, err = mfd.Write(data); err != nil {
		return nil, fmt.Errorf("writting data to anonymous file: %v", err)
	}
	return mfd, nil
}

// SaveFile saves a file in the supplied path with the supplied content.
func SaveFile(path string, content []byte) {
	saveFile(path, "", content)
}

// SaveFile saves a json file in the supplied path with the supplied content.
func SaveJSON(path string, content []byte) {
	saveFile(path, ".json", content)
}

// SaveFile saves a markdown in the supplied path with the supplied content.
func SaveMarkdown(path string, content []byte) {
	saveFile(path, ".md", content)
}

// SaveFile saves a file in the supplied path with the supplied content and extension.
//
// It panics if something goes wrong.
func saveFile(path, ext string, content []byte) {
	if err := ioutil.WriteFile(fmt.Sprintf("%s%s", path, ext), content, 0644); err != nil {
		live.Printer.Fatalf("writing \"%s\": %v", fmt.Sprintf("%s.%s", path, ext), err)
	}
}
