package op

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/mholt/archives"
	p "github.com/plasticgaming99/integra/internal/printdbg"
	"github.com/plasticgaming99/integra/lib/db"
	"github.com/plasticgaming99/integra/lib/pkg"
)

func Install(packDB []pkg.Packinfo, localDB db.LocalDB, archivePath string, rootDir string) error {
	var pkinfo pkg.Packinfo
	var intgpack archives.Tar
	abs, err := filepath.Abs(archivePath)
	if err != nil {
		fmt.Println(err)
	}
	in, err := os.Open(abs)
	if err != nil {
		fmt.Println("Error opening archive file")
		return err
	}
	defer in.Close()
	input := bufio.NewReader(in)

	{
		fsys := archives.ArchiveFS{
			Path:   archivePath,
			Format: archives.Tar{},
		}
		f, err := fsys.Open(".PACKAGE")
		if err != nil {
			log.Fatal(err)

		}
		read := bufio.NewReader(f)
		pkg.ReadPackinfo(read)
	}

	fmt.Println("checking conflicts...")
	for i := range packDB {
		p.Printdbg("checking", packDB[i].Packagename)
		if slices.Contains(packDB[i].Conflicts, pkinfo.Packagename) {
			log.Fatal(pkinfo.Packagename, " and ", packDB[i].Packagename, " is conflicting")
		}
		if slices.Contains(pkinfo.Conflicts, packDB[i].Packagename) {
			log.Fatal(pkinfo.Packagename, " and ", packDB[i].Packagename, " is conflicting")
		}
	}

	installArchive := func(ctx context.Context, f archives.FileInfo) error {
		p.Printdbg(f.Name())
		destpath := rootDir + f.NameInArchive
		if f.IsDir() {
			err := os.MkdirAll(destpath, f.Mode().Perm())
			p.Printdbg(err)
			return err
		}

		if f.LinkTarget != "" {
			targ := f.LinkTarget
			p.Printdbg("symlink target:", targ)
			if err != nil {
				p.Printdbg(err)
				return err
			}
			nameinarchiveABS, err := filepath.Abs(filepath.Join(rootDir, f.NameInArchive))
			if err != nil {
				p.Printdbg(err)
				return err
			}
			err = os.Symlink(targ, nameinarchiveABS)
			p.Printdbg(err)
			return err
		}

		reader, err := f.Open()
		if err != nil {
			p.Printdbg(err)
			return err
		}
		bufread := bufio.NewReader(reader)
		defer reader.Close()

		if f.Name() == ".PACKAGE" || f.Name() == ".MTREE" {

		}

		destfile, err := os.OpenFile(destpath, os.O_CREATE|os.O_WRONLY, f.Mode().Perm())
		if err != nil {
			p.Printdbg(err)
			return err
		}
		bufdest := bufio.NewWriter(destfile)
		defer destfile.Close()

		_, err = io.Copy(bufdest, bufread)
		if err != nil {
			p.Printdbg(err)
			return err
		}
		return nil
	}

	p.Printdbg("installing ", pkinfo.Packagename)

	return intgpack.Extract(context.Background(), input, installArchive)
}
