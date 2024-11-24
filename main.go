// artex package manager
// integra (much unstable such wow)

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"
)

var (
	install  = false
	sync     = false
	search   = false
	refresh  = false
	refforce = false
	upgrade  = false
	allyes   = false

	rootdir = "/"

	debug = false

	pack2ins []string
)

func main() {
	fmt.Println("integra test")
	parse(os.Args[1:])
	printdbg(sync, search, refresh, refforce, upgrade, allyes)
	printdbg(pack2ins)
	printdbg(rootdir)

	var intgpack archives.Tar
	abs, err := filepath.Abs(pack2ins[0])
	in, err := os.Open(abs)
	if err != nil {
		fmt.Println("Error opening archive file")
	}
	defer in.Close()
	input := bufio.NewReader(in)

	handler := func(ctx context.Context, f archives.FileInfo) error {
		destpath := rootdir + f.NameInArchive
		if f.IsDir() {
			err := os.MkdirAll(destpath, f.Mode().Perm())
			return err
		}

		if f.LinkTarget != "" {
			targAbs, err := filepath.Abs(rootdir + f.LinkTarget)
			if err != nil {
				return err
			}
			nameAbs, err := filepath.Abs(rootdir + f.NameInArchive)
			if err != nil {
				return err
			}
			err = os.Symlink(targAbs, nameAbs)
			return err
		}

		reader, err := f.Open()
		if err != nil {
			return err
		}
		bufread := bufio.NewReader(reader)
		defer reader.Close()

		destfile, err := os.OpenFile(destpath, os.O_CREATE|os.O_WRONLY, f.Mode())
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
	intgpack.Extract(context.Background(), input, handler)

	//err := intgpack.Extract(nil, nil, nil)
}

func parse(in []string) {
	for iv := 0; len(in) > iv; iv++ {
		if strings.Contains(in[iv], "--") && !strings.Contains(in[iv], "=") {
			fmt.Println("i found the option yay")
			// check
			switch in[iv][2:] {
			case "dbg":
				debug = true
			case "install":
				install = true
			case "sync":
				sync = true
			case "search":
				search = true
			case "refresh":
				if refresh {
					refforce = true
				}
				refresh = true
			case "upgrade":
				upgrade = true
			case "yes":
				allyes = true
			default:
				fmt.Println("Error when parsing:", in[iv])
			}
		} else if strings.Contains(in[iv], "--") && strings.Contains(in[iv], "=") {
			i := strings.Split(in[iv][2:], "=")
			switch i[0] {
			case "override-root":
				dir, err := filepath.Abs(i[1])
				if err != nil {
					fmt.Println("errrrr")
				}
				rootdir = dir + "/"
			}
		} else
		// check combined option like pacman
		if strings.HasPrefix(in[iv], "-") {
		} else {
			pack2ins = append(pack2ins, in[iv])
		}
	}
}

func printdbg(a ...any) {
	if debug {
		fmt.Println(a...)
	}
}
