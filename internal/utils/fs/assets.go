package fs

import (
	"fmt"
	io "io/fs"
)

type assets struct {
	fs io.FS
}

// Assets provides access to the assets filesystem.
var Assets *assets

// read reads the file on the Assets filesystem from the supplied path on the Assets filesystem.
//
// Returns the content of the file and any error encountered.
func (a *assets) read(path string) ([]byte, error) {
	file, err := io.ReadFile(a.fs, path)
	if err != nil {
		return nil, fmt.Errorf("reading from assets filesystem: %v", err)
	}
	return file, err
}

// NewFS initializes the Assets filesystem.
func (a *assets) NewFS(fs io.FS) {
	a.fs = fs
}

// ReadAsset reads the supplied asset file on the Assets filesystem.
//
// Returns the content of the file and any error encountered.
func (a *assets) Read(asset string) ([]byte, error) {
	return a.read(asset)
}
