package fs

import (
	"fmt"
	io "io/fs"
)

type lang struct {
	fs io.FS
}

// Lang provides access to the lang filesystem.
var Lang *lang

// read reads the file on the Lang filesystem from the supplied path on the Lang filesystem.
//
// Returns the content of the file and any error encountered.
func (l *lang) read(path string) ([]byte, error) {
	file, err := io.ReadFile(l.fs, path)
	if err != nil {
		return nil, fmt.Errorf("reading from lang filesystem: %v", err)
	}
	return file, err
}

// NewFS initializes the Lang filesystem.
func (l *lang) NewFS(fs io.FS) {
	l.fs = fs
}

// ReadMessages reads the messages file from the supplied language on the Lang filesystem.
//
// Returns the content of the file and any error encountered.
func (l *lang) ReadMessages(lang string) ([]byte, error) {
	dir := fmt.Sprintf("%s.yaml", lang)
	return l.read(dir)
}
