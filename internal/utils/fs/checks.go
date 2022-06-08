package fs

import (
	"fmt"
	io "io/fs"
	"runtime"
	"strings"

	"github.com/quimera-project/quimera/internal/env"
)

type checks struct {
	fs io.FS
}

// Checks provides access to the checks filesystem.
var Checks *checks

// read reads the file on the Checks filesystem from the supplied path on the Checks filesystem.
//
// Returns the content of the file and any error encountered.
func (c *checks) read(path string) ([]byte, error) {
	file, err := io.ReadFile(c.fs, path)
	if err != nil {
		return nil, fmt.Errorf("reading from checks filesystem: %v", err)
	}
	return file, err
}

// NewFS initializes the Checks filesystem.
func (c *checks) NewFS(fs io.FS) {
	c.fs = fs
}

// FS returns the Templates filesystem.
func (c *checks) FS() io.FS {
	return c.fs
}

// GetCategories returns a string slice with all the categories names on the Checks filesystem and any error encountered.
func (c *checks) GetCategories() ([]string, error) {
	var categories []string
	dir, err := io.ReadDir(c.fs, ".")
	if err != nil {
		return nil, fmt.Errorf("reading checks filesystem directory: %v", err)
	}
	for _, c := range dir {
		categories = append(categories, c.Name())
	}
	return categories, nil
}

// ReadBinaryByCategory reads the binary file from the supplied category and binary name on the Checks filesystem.
//
// Returns the content of the file and any error encountered.
func (c *checks) ReadBinaryByCategory(cat, bin string) ([]byte, error) {
	var (
		dir string
	)
	if strings.HasSuffix(bin, ".sh") {
		dir = fmt.Sprintf("%s/bin/%s", cat, bin)
	} else {
		dir = fmt.Sprintf("%s/bin/%s_%s", cat, bin, runtime.GOARCH)
	}
	return c.read(dir)
}

// ReadDocByCategory reads the documentation file from the supplied category and filename on the Checks filesystem.
//
// Returns the content of the file and any error encountered.
func (c *checks) ReadDocByCategory(cat, doc string) ([]byte, error) {
	dir := fmt.Sprintf("%s/doc/%s/%s.md", cat, env.Config.Language, doc)
	file, err := c.read(dir)
	if err != nil {
		dir := fmt.Sprintf("%s/doc/en/%s.md", cat, doc)
		return c.read(dir)
	}
	return file, nil
}

// ReadMeta reads the meta file from the supplied category on the Checks filesystem.
//
// Returns the content of the file and any error encountered.
func (c *checks) ReadMeta(cat string) ([]byte, error) {
	dir := fmt.Sprintf("%s/meta.yaml", cat)
	return c.read(dir)
}

// ReadCategoryLang reads the lang file from the supplied category and language on the Checks filesystem.
//
// Returns the content of the file and any error encountered.
func (c *checks) ReadCategoryLang(cat, lang string) ([]byte, error) {
	dir := fmt.Sprintf("%s/lang/%s.yaml", cat, lang)
	return c.read(dir)
}
