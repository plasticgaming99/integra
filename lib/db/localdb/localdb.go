package localdb

import (
	"io"
	"io/fs"
)

type LocalDB struct {
	fsystem fs.FS
}

func (ldb *LocalDB) AddFile(reader io.Reader, pkgname string, fname string) error {
	return nil
}

func (ldb *LocalDB) DelFile(pkgname string, fname string) error {
	return nil
}
