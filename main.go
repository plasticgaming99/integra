// artex package manager
// integra (much unstable such wow)

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/mholt/archives"
)

var (
	install = false
	sync    = false
	upgrade = false
	remove  = false
	search  = false
	allyes  = false

	rootdir = "/"
	dbdir   = "/var/lib/integra/db"

	debug = false

	pack2ins []string
)

type packinfo struct {
	packagename  string
	version      string
	release      int
	license      string
	architecture string
	description  string
	url          string
	depends      []string
	optdeps      []string
	builddeps    []string
	conflicts    []string
	provides     []string
}

func main() {
	fmt.Println("integra test")
	if len(os.Args) == 1 {
		fmt.Println("No arguments, abort.")
		os.Exit(0)
	}
	if upperAvailable(os.Args[1]) {
		parseOp(os.Args[1])
		parse(os.Args[2:])
	} else {
		parse(os.Args[1:])
	}
	printdbg(sync, search, upgrade, allyes)
	printdbg(pack2ins)
	printdbg(rootdir)
	printdbg(dbdir)

	_, localtrdb, err := readLocalDB(filepath.Join(dbdir, "local"))
	if err != nil {
		log.Fatal("failed to read db")
	}

	printdbg("entries???:", len(localtrdb))
	for i := 0; i < len(localtrdb); i++ {
		printdbg(localtrdb[i][0][0])
	}

	for cnt := 1; cnt < len(pack2ins); cnt++ {
		if install {
			var packagename string
			var intgpack archives.Tar
			abs, err := filepath.Abs(pack2ins[cnt])
			if err != nil {
				fmt.Println(err)
			}
			in, err := os.Open(abs)
			if err != nil {
				fmt.Println("Error opening archive file")
			}
			defer in.Close()
			input := bufio.NewReader(in)

			getPackName := func(ctx context.Context, f archives.FileInfo) error {
				if f.NameInArchive == ".PACKAGE" {
					fssss, err := f.Open()
					if err != nil {
						return err
					}
					defer fssss.Close()
					toreadbyte, err := io.ReadAll(fssss)
					if err != nil {
						return err
					}
					stringToOperate := string(toreadbyte)
					slice2op := strings.Split(stringToOperate, "\n")
					for i := 0; i < len(slice2op); i++ {
						if s := strings.Split(slice2op[i], " = "); s[0] == "package" {
							packagename = s[1]
						}
					}
				}
				if f.IsDir() {
					return nil
				}
				return nil
			}

			installArchive := func(ctx context.Context, f archives.FileInfo) error {
				printdbg(f.Name())
				destpath := rootdir + f.NameInArchive
				if f.IsDir() {
					err := os.MkdirAll(destpath, f.Mode().Perm())
					printdbg(err)
					return err
				}

				if f.LinkTarget != "" {
					targAbs, err := filepath.Abs(filepath.Join(rootdir, f.LinkTarget))
					printdbg("symlink target:", targAbs)
					if err != nil {
						printdbg(err)
						return err
					}
					nameinarchiveABS, err := filepath.Abs(filepath.Join(rootdir, f.NameInArchive))
					if err != nil {
						printdbg(err)
						return err
					}
					err = os.Symlink(targAbs, nameinarchiveABS)
					printdbg(err)
					return err
				}

				reader, err := f.Open()
				if err != nil {
					printdbg(err)
					return err
				}
				bufread := bufio.NewReader(reader)
				defer reader.Close()

				if f.Name() == ".PACKAGE" || f.Name() == ".MTREE" {
					destpath = filepath.Join(dbdir, "local", packagename, f.NameInArchive)
					mdAllIfNeeded(filepath.Join(dbdir, "local", packagename))
					dest, err := os.OpenFile(destpath, os.O_CREATE|os.O_WRONLY, f.Mode())
					if err != nil {
						log.Fatal("error writing db")
					}
					io.Copy(dest, bufread)
					return nil
				}

				destfile, err := os.OpenFile(destpath, os.O_CREATE|os.O_WRONLY, f.Mode().Perm())
				if err != nil {
					printdbg(err)
					return err
				}
				bufdest := bufio.NewWriter(destfile)
				defer destfile.Close()

				_, err = io.Copy(bufdest, bufread)
				if err != nil {
					printdbg(err)
					return err
				}
				return nil
			}

			{
				abs, err := filepath.Abs(pack2ins[cnt])
				if err != nil {
					fmt.Println(err)
				}
				in, err := os.Open(abs)
				if err != nil {
					fmt.Println("Error opening archive file")
				}
				defer in.Close()
				tInput := bufio.NewReader(in)

				intgpack.Extract(context.Background(), tInput, getPackName)
			}

			intgpack.Extract(context.Background(), input, installArchive)
			// err := intgpack.Extract(nil, nil, nil)
		}

		if remove {
			printdbg("start remove")
			var inttt int
			printdbg("toremove:", pack2ins[cnt])
			for i := 0; i < len(pack2ins); i++ {
				if localtrdb[i][0][0] == pack2ins[cnt] {
					printdbg("yes")
					inttt = i
				}
			}
			printdbg("db id:", inttt)
			for i := 2; i < len(localtrdb[inttt])-2; i++ {
				path2remove := localtrdb[inttt][i][0]
				path2remove, err := filepath.Abs(filepath.Join(rootdir, path2remove))
				if err != nil {
					log.Fatal("error during removing")
				}
				if !(strings.HasPrefix(path2remove, "/set") || strings.HasPrefix(path2remove, "./.PACKAGE")) {
					printdbg(path2remove)
					os.Remove(path2remove)
				}
			}
			os.RemoveAll(filepath.Join(dbdir, pack2ins[cnt]))
		}
	}
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
			case "upgrade":
				upgrade = true
			case "remove":
				remove = true
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
				dbdir, err = filepath.Abs(rootdir + dbdir)
				if err != nil {
					fmt.Println("err")
				}
			}
		} else
		// check combined option like pacman
		if strings.HasPrefix(in[iv], "-") {
		} else {
			pack2ins = append(pack2ins, in[iv])
		}
	}
}

func parseOp(in string) {
	run := []rune(in)
	for i := 0; i < len(run); i++ {
		if unicode.IsUpper(run[i]) {
			switch string(run[i]) {
			case "I":
				install = true
			case "S":
				sync = true
			case "U":
				upgrade = true
			case "R":
				remove = true
			}
		}
	}
}

func readLocalDB(dbdir string) ([]packinfo, [][][]string, error) {
	returninfo := []packinfo{}
	returntree := [][][]string{}
	err := filepath.WalkDir(dbdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error reading database")
		}
		if d.IsDir() {
			return nil
		}
		switch d.Name() {
		case ".PACKAGE":
			add := packinfo{}
			iFile, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading database: package")
			}
			inputFile := strings.Split(string(iFile), "\n")
			printdbg("db input length:", len(inputFile))

			// same way with build
			for i := 0; i < len(inputFile); i++ {
				inputVar := strings.Split(inputFile[i], " = ")
				switch inputVar[0] {
				case "package":
					printdbg("packname", inputVar[1])
					add.packagename = inputVar[1]
				case "version":
					add.version = inputVar[1]
				case "release":
					a, err := strconv.Atoi(inputVar[1])
					if err != nil {
						fmt.Println("release number is not int")
					}
					add.release = a
				case "license":
					add.license = inputVar[1]
				case "architecture":
					add.architecture = inputVar[1]
				case "description":
					add.description = inputVar[1]
				case "depends":
					add.depends = append(add.depends, inputVar[1])
				case "optdeps":
					add.optdeps = append(add.optdeps, inputVar[1])
				case "builddeps":
					add.builddeps = append(add.builddeps, inputVar[1])
				case "conflicts":
					add.depends = append(add.depends, inputVar[1])
				case "provides":
					add.provides = append(add.provides, inputVar[1])
				case "url":
					add.url = inputVar[1]
				}
			}
			returninfo = append(returninfo, add)
			return nil
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
			/*for i := 0; i < len(rows); i++ {
				fmt.Println(rows[i][0])
			}*/
			rows[0][0] = strings.Split(filepath.ToSlash(path), "/")[len(strings.Split(filepath.ToSlash(path), "/"))-2]
			printdbg("waa", strings.Split(filepath.ToSlash(path), "/")[len(strings.Split(filepath.ToSlash(path), "/"))-2])
			returntree = append(returntree, rows)
			return nil
		}
		return fmt.Errorf("DB directory might be broken")
	})
	if err == nil {
		return returninfo, returntree, nil
	}
	errorpi := []packinfo{}
	return errorpi, nil, nil
}

func readRemoteDB(st string) {
}

func printdbg(a ...any) {
	if debug {
		fmt.Println(a...)
	}
}

func mdAllIfNeeded(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		err := os.MkdirAll(path, os.ModePerm)
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("file exists")
	}
	return nil
}

func upperAvailable(st string) bool {
	for _, ru := range st {
		if unicode.IsUpper(ru) {
			return true
		}
	}
	return false
}
