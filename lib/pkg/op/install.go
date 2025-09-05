package op

import (
	"bufio"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/mholt/archives"
	"github.com/plasticgaming99/integra/lib/db/localdb"
	"github.com/plasticgaming99/integra/lib/pkg/types"
)

type handleArchive struct {
	pkg     types.Pkg
	rootDir string
	localdb *localdb.LocalDB
}

func (h *handleArchive) fileHandler(ctx context.Context, f archives.FileInfo) error {
	destpath := filepath.Join(h.rootDir, f.NameInArchive)
	if f.IsDir() {
		err := os.MkdirAll(destpath, f.Mode().Perm())
		return err
	}

	if f.LinkTarget != "" {
		targ := f.LinkTarget
		nameinarchiveABS, err := filepath.Abs(filepath.Join(h.rootDir, f.NameInArchive))
		if err != nil {
			return err
		}
		err = os.Symlink(targ, nameinarchiveABS)
		return err
	}

	reader, err := f.Open()
	if err != nil {
		return err
	}
	bufread := bufio.NewReader(reader)
	defer reader.Close()

	if f.Name() == ".PACKAGE" || f.Name() == ".INTEGRITY" || f.Name() == ".MTREE" {
		h.localdb.AddFile(h.pkg, f.Name(), bufread)
		return nil
	}

	destfile, err := os.OpenFile(destpath, os.O_CREATE|os.O_WRONLY, f.Mode().Perm())
	if err != nil {
		return err
	}
	bufdest := bufio.NewWriter(destfile)
	defer destfile.Close()

	_, err = io.Copy(bufdest, bufread)
	if err != nil {
		return err
	}
	return nil
}

// install installs package, do not care about dependency
func Install(filepath string, rootdir string, ldb localdb.LocalDB) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	afs := archives.ArchiveFS{
		Path:   filepath,
		Format: archives.Tar{},
	}
	pkinfofile, err := afs.Open(".PACKAGE")
	if err != nil {
		return err
	}
	defer pkinfofile.Close()
	pkinfofilebuf := bufio.NewReader(pkinfofile)
	pkinfo := types.ReadPackinfo(pkinfofilebuf)

	buffile := bufio.NewReader(file)
	tzst := archives.Tar{}
	fh := handleArchive{
		pkg:     types.PackInfoToPkg(pkinfo),
		rootDir: rootdir,
		localdb: &ldb,
	}
	tzst.Extract(context.TODO(), buffile, fh.fileHandler)
	return nil
}
