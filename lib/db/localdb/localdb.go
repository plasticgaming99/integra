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

const (
	PackageFile   = ".PACKAGE"
	IntegrityFile = ".INTEGRITY"

	lDBDirPerm  = 0755
	lDBFilePerm = 0644
)

type LocalDB struct {
	root os.Root
}

func InitializeDBDir(path string) error {
	return os.MkdirAll(path, lDBDirPerm)
}

// enter db dir for path
func InitializeLocalDB(path string) error {
	return os.Mkdir(filepath.Join(path, "localdb"), lDBDirPerm)
}

// enter db dir for path
func OpenLocalDB(path string) (LocalDB, error) {
	ldb := LocalDB{}
	root, err := os.OpenRoot(filepath.Join(path, "localdb"))
	if err != nil {
		return LocalDB{}, fmt.Errorf("error opening db: %w", err)
	}
	ldb.root = *root
	return ldb, nil
}

func PkgToDirname(pkg types.Pkg) string {
	dirname := pkg.PkgName + "-" + pkg.Version + "-" + strconv.Itoa(pkg.Release)
	return dirname
}

// packagename-version-release
func DirnameToPkg(s string) (pkg types.Pkg) {
	rIndex := strings.LastIndexByte(s, '-')
	if rIndex == -1 {
		return types.Pkg{}
	}
	vIndex := strings.LastIndexByte(s[:rIndex], '-')

	rel, err := strconv.Atoi(s[rIndex+1:])
	if err != nil {
		fmt.Println("db error: release field incorrect")
		rel = 0
	}
	pkg = types.Pkg{
		PkgName: s[:vIndex],
		Version: s[vIndex+1 : rIndex],
		Release: rel,
	}
	return
}

// generate merged filepath from package info
func getFileName(pkg types.Pkg, filename string) string {
	dirname := PkgToDirname(pkg)
	return filepath.Join(dirname, filename)
}

// just initialize pkg dir
func (ldb *LocalDB) InitLocalDBPkgDir(pkg types.Pkg) error {
	return ldb.root.Mkdir(PkgToDirname(pkg), lDBDirPerm)
}

// register a file named fname linked to pkg
func (ldb *LocalDB) AddFile(pkg types.Pkg, fname string, reader io.Reader) error {
	filePath := getFileName(pkg, fname)
	ldb.InitLocalDBPkgDir(pkg)
	f, err := ldb.root.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file in DB: %w", err)
	}
	defer f.Close()
	f.Chmod(lDBFilePerm)
	wr := bufio.NewWriter(f)
	defer wr.Flush()
	_, err2 := io.Copy(wr, reader)
	if err2 != nil {
		return fmt.Errorf("error writing file to DB: %w", err2)
	}
	return nil
}

// not so fuzzy but it selects similar version and release
// ordered like pkgname -> version -> release
// pkgname is not searched fuzzy
// finish earlier if it find perfect match
func (ldb *LocalDB) GetKeyPkgFuzzy(pkg types.Pkg) types.Pkg {
	rtpkg := types.Pkg{}
	rfs := ldb.root.FS()
	fs.WalkDir(rfs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			dirpkg := DirnameToPkg(path)
			if dirpkg.PkgName == pkg.PkgName {
				rtpkg = dirpkg
			}
		}
		return nil
	})
	return rtpkg
}

// returns reader from file
func (ldb *LocalDB) GetFile(pkg types.Pkg, fname string) (rd io.Reader, er error) {
	return ldb.root.Open(getFileName(pkg, fname))
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
	err := ldb.root.RemoveAll(PkgToDirname(pkg))
	if err != nil {
		return fmt.Errorf("error failed unregistering pkg: %w", err)
	}
	return nil
}

// aa
// pkg, filename
// notnow
func (ldb *LocalDB) CollectGarbageFunc(fn func(types.Pkg, string) error) {
	fs.WalkDir(ldb.root.FS(), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return nil
	})
}
