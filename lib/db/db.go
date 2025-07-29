package db

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/plasticgaming99/integra/internal/dirutil"
	"github.com/plasticgaming99/integra/lib/pkg"
)

type LocalDB struct {
	Path string
}

func (ldb LocalDB) AddFile(pkgname string, fname string, input io.Reader) error {
	if fname == ".PACKAGE" || fname == ".INTEGRITY" || fname == ".MTREE" {
		destpath := filepath.Join(ldb.Path, "local", pkgname, fname)
		dirutil.MdAllIfNeeded(filepath.Join(ldb.Path, "local", pkgname))
		dest, err := os.OpenFile(destpath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("error writing db")
		}
		io.Copy(dest, input)
		dest.Close()
	}
	return nil
}

func ReadLocalDB(dbdir string) ([]pkg.Packinfo, []pkg.Integrity, error) {
	returninfo := []pkg.Packinfo{}
	returntree := []pkg.Integrity{}
	err := filepath.WalkDir(dbdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error reading database")
		}
		if d.IsDir() {
			return nil
		}
		switch d.Name() {
		case ".PACKAGE":
			add := pkg.Packinfo{}
			iFile, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading database: package")
			}
			inputFile := strings.Split(string(iFile), "\n")
			//printdbg("db input length:", len(inputFile))

			// same way with build
			for i := 0; i < len(inputFile); i++ {
				inputVar := strings.Split(inputFile[i], " = ")
				switch inputVar[0] {
				case "package":
					//printdbg("packname", inputVar[1])
					add.Packagename = inputVar[1]
				case "version":
					add.Version = inputVar[1]
				case "release":
					a, err := strconv.Atoi(inputVar[1])
					if err != nil {
						fmt.Println("release number is not int")
					}
					add.Release = a
				case "license":
					add.License = inputVar[1]
				case "architecture":
					add.Architecture = inputVar[1]
				case "description":
					add.Description = inputVar[1]
				case "depends":
					add.Depends = append(add.Depends, inputVar[1])
				case "optdeps":
					add.Optdeps = append(add.Optdeps, inputVar[1])
				case "builddeps":
					add.Builddeps = append(add.Builddeps, inputVar[1])
				case "conflicts":
					add.Depends = append(add.Depends, inputVar[1])
				case "provides":
					add.Provides = append(add.Provides, inputVar[1])
				case "url":
					add.Url = inputVar[1]
				}
			}
			returninfo = append(returninfo, add)
			return nil
		case ".INTEGRITY":
			f, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("error opening file: " + fmt.Sprint(err))
			}
			b := bufio.NewReader(f)
			returntree = append(returntree, ReadIntegrity(b, filepath.Dir(path)))
			f.Close()
		case ".MTREE":
			file, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error during reading mtree")
			}
			strs := string(file)
			slcs := strings.Split(strs, "\n")
			rows := [][]string{}
			for i := 0; i < len(slcs); i++ {
				rowstmp := strings.Split(slcs[i], " ")
				rows = append(rows, rowstmp)
			}
			// convert to integrity from here
			rows[0][0] = strings.Split(filepath.ToSlash(path), "/")[len(strings.Split(filepath.ToSlash(path), "/"))-2]
			returntree = append(returntree, pkg.Integrity{Filepath: rows[0][0]})
			return nil
		}
		return fmt.Errorf("DB directory might be broken")
	})
	if err == nil {
		return returninfo, returntree, nil
	}
	errorpi := []pkg.Packinfo{}
	errorintg := []pkg.Integrity{}
	return errorpi, errorintg, nil
}

// read integrity file yay
func ReadIntegrity(input io.Reader, ppkg string) (intg pkg.Integrity) {
	intg.ParentPkg = ppkg
	// make a buffer
	sc := bufio.NewScanner(input)
	dirname := ""
	for sc.Scan() {
		text := sc.Text()
		splittedtxt := filepath.SplitList(text)
		if len(splittedtxt) == 1 {
			// perhaps root dir
			dirname = text
		} else {
			// perhaps a file
			f := strings.Fields(text)
			for _, key := range f {
				bef, aft, fou := strings.Cut(key, "=")
				if !fou {
					// yas filename
					intg.Filepath = filepath.Join(dirname, bef)
					continue
				}
				switch bef {
				case "perm":
					strconv.ParseUint(aft, 8, 32)
				case "uid":
					intg.FileUid = 0 //currently
				case "gid":
					intg.FileGid = 0 //currently
				case "blake3sum":
					intg.Blake3Sum = aft
				}
			}
		}
	}
	return
}

type RemoteDB struct {
	Path string
}

func ReadRemoteDB(st string) {
}
