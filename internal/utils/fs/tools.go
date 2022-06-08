package fs

import (
	"fmt"
	io "io/fs"
)

type tools struct {
	fs io.FS
}

// Tools provides access to the tools filesystem.
var Tools *tools

// read reads the file on the Tools filesystem from the supplied path on the Tools filesystem.
//
// Returns the content of the file and any error encountered.
func (t *tools) read(path string) ([]byte, error) {
	file, err := io.ReadFile(t.fs, path)
	if err != nil {
		return nil, fmt.Errorf("reading from tools filesystem: %v", err)
	}
	return file, err
}

// NewFS initializes the Tools filesystem.
func (t *tools) NewFS(fs io.FS) {
	t.fs = fs
}

// ReadAsset reads the supplied asset file on the Tools filesystem.
//
// Returns the content of the file and any error encountered.
func (t *tools) Read(asset string) ([]byte, error) {
	return t.read(asset)
}
