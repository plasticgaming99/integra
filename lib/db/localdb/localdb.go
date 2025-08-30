package localdb

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/plasticgaming99/integra/lib/pkg/types"
)

type LocalDB struct {
	root os.Root
}

func (ldb *LocalDB) OpenDB(path string) error {
	root, err := os.OpenRoot(path)
	if err != nil {
		return fmt.Errorf("error opening db: %w", err)
	}
	ldb.root = *root
	return nil
}

func GetPkgDirname(pkg types.Pkg) string {
	dirname := pkg.PkgName + "-" + strconv.Itoa(pkg.Release) + "-" + pkg.Version
	return dirname
}

// packagename-version-release
func DirnameToPkg(s string) (pkg types.Pkg) {
	rIndex := strings.LastIndexByte(s, '-')
	vIndex := strings.LastIndexByte(s[:rIndex], '-')

	rel, err := strconv.Atoi(s[rIndex+1:])
	if err != nil {
		fmt.Println("db error: release field incorrect")
		rel = 0
	}
	pkg = types.Pkg{
		PkgName: s[:vIndex-1],
		Version: s[vIndex+1 : rIndex-1],
		Release: rel,
	}
	return
}

// generate merged filepath from package info
func getFileName(pkg types.Pkg, filename string) string {
	dirname := GetPkgDirname(pkg)
	return filepath.Join(dirname, filename)
}

// register a file named fname linked to pkg
func (ldb *LocalDB) AddFile(pkg types.Pkg, fname string, reader io.Reader) error {
	filePath := getFileName(pkg, fname)
	f, err := ldb.root.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file in DB: %w", err)
	}
	defer f.Close()
	wr := bufio.NewWriter(f)
	defer wr.Flush()
	_, err2 := io.Copy(wr, reader)
	if err2 != nil {
		return fmt.Errorf("error writing file to DB: %w", err2)
	}
	return nil
}

// delete file named fname that linked to the pkg
func (ldb *LocalDB) DelFile(pkg types.Pkg, fname string) error {
	filePath := getFileName(pkg, fname)
	err := ldb.root.Remove(filePath)
	if err != nil {
		return fmt.Errorf("error removing file from DB: %w", err)
	}
	return nil
}

// purge db linked with db
func (ldb *LocalDB) UnregisterPkg(pkg types.Pkg) error {
	err := ldb.root.RemoveAll(GetPkgDirname(pkg))
	if err != nil {
		return fmt.Errorf("error failed unregistering pkg: %w", err)
	}
	return nil
}

// aa
// pkg, filename
// notnow
func (ldb *LocalDB) CollectGarbageFunc(fn func(types.Pkg, string) error) {
	fs.WalkDir(ldb.root.FS(), "/", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return nil
	})
}
