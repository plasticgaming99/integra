package op

import (
	"bufio"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/mholt/archives"
)

type handleArchive struct {
	rootDir string
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
		//localDB.AddFile(pkinfo.Packagename, f.Name(), reader)
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

// install installs package
func Install(filepath string, rootdir string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	buffile := bufio.NewReader(file)
	tzst := archives.Tar{}
	fh := handleArchive{
		rootDir: rootdir,
	}
	tzst.Extract(context.TODO(), buffile, fh.fileHandler)
	return nil
}
