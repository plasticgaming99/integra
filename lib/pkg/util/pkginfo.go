// utilities for getting package info
package util

import (
	"fmt"
	"io/fs"

	"github.com/mholt/archives"
	"github.com/plasticgaming99/integra/lib/integrity"
	"github.com/plasticgaming99/integra/lib/pkg"
	"github.com/plasticgaming99/integra/lib/pkg/types"
)

func getFileFromArchive(path string, fname string) (fs.File, error) {
	a := archives.ArchiveFS{
		Path:   path,
		Format: archives.Tar{}, // i'm sure you using tarball
	}
	f, err := a.Open(fname)
	return f, err
}

func GetPackinfo(path string) (types.Packinfo, error) {
	f, err := getFileFromArchive(path, ".PACKAGE")
	if err != nil {
		return types.Packinfo{}, fmt.Errorf("failed to open .PACKAGE in archive %s: %w", path, err)
	}
	i := pkg.ReadPackinfo(f)
	f.Close()
	return i, nil
}

func GetIntegrity(path string) (integrity.Integrity, error) {
	f, err := getFileFromArchive(path, ".INTEGRITY")
	if err != nil {
		return integrity.Integrity{}, fmt.Errorf("failed to open .INTEGRITY in archive %s: %w", path, err)
	}
	intg := integrity.Parse(f)
	f.Close()
	return intg, nil
}
