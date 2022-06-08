package fs

import (
	"fmt"
	io "io/fs"
)

type templates struct {
	fs io.FS
}

// Templates provides access to the templates filesystem.
var Templates *templates

// read reads the file on the Templates filesystem from the supplied path on the Templates filesystem.
//
// Returns the content of the file and any error encountered.
func (t *templates) read(path string) ([]byte, error) {
	file, err := io.ReadFile(t.fs, path)
	if err != nil {
		return nil, fmt.Errorf("reading on templates filesystem: %v", err)
	}
	return file, err
}

// NewFS initializes the Templates filesystem.
func (t *templates) NewFS(fs io.FS) {
	t.fs = fs
}

// FS returns the Templates filesystem.
func (t *templates) FS() io.FS {
	return t.fs
}

// ReadObsidianFile reads the obsidian file from the supplied obsidian path on the Templates filesystem.
//
// Returns the content of the file and any error encountered.
func (t *templates) ReadObsidianFile(path string) ([]byte, error) {
	dir := fmt.Sprintf("obsidian/%s", path)
	return t.read(dir)
}

// ReadHtmlFile reads the html file from the supplied html path on the Templates filesystem.
//
// Returns the content of the file and any error encountered.
func (t *templates) ReadHtmlFile(path string) ([]byte, error) {
	dir := fmt.Sprintf("html/%s", path)
	return t.read(dir)
}

// ReadTemplateTheme reads the yaml theme file from the supplied theme name on the Templates filesystem.
//
// Returns the content of the file and any error encountered.
func (t *templates) ReadTemplateTheme(theme string) ([]byte, error) {
	dir := fmt.Sprintf("themes/%s.yaml", theme)
	return t.read(dir)
}
