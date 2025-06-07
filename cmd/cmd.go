package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	p "github.com/plasticgaming99/integra/internal/printdbg"
	"github.com/plasticgaming99/integra/lib/db"
	"github.com/plasticgaming99/integra/lib/pkg/op"

	"github.com/BurntSushi/toml"
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

// parse config with toml
type (
	configfile struct {
		repo []repo
	}
	config struct {
		parallel bool
	}
	repo struct {
		server  string
		include string
	}
)

func Exec(args []string) {
	fmt.Println("integra test")
	if len(os.Args) == 1 {
		fmt.Println("No arguments, abort.")
		os.Exit(0)
	}
	if upperAvailable(args[1]) {
		parseOp(args[1])
		parse(args[2:])
	} else {
		parse(args[1:])
	}
	if debug {
		p.Dbg()
	}
	p.Printdbg(install, sync, search, upgrade, allyes)
	p.Printdbg(pack2ins)
	p.Printdbg(rootdir)
	p.Printdbg(dbdir)

	localdb, localtrdb, err := db.ReadLocalDB(filepath.Join(dbdir, "local"))
	if err != nil {
		log.Fatal("failed to read db")
	}

	// debug?
	p.Printdbg("entries???:", len(localtrdb))
	for i := 0; i < len(localtrdb); i++ {
		p.Printdbg(localtrdb[i][0][0])
	}

	// remote repo
	configpath := filepath.Join(rootdir, "etc", "integra.conf")
	var config configfile
	_, err = toml.DecodeFile(configpath, &config)
	if err != nil {
		fmt.Println("error parsing config file")
		os.Exit(1)
	}

	for _, s := range config.repo {
		fmt.Println(s)
	}

	localdbdir := filepath.Join(dbdir, "local")

	for cnt := 0; cnt < len(pack2ins); cnt++ {
		if install {
			err := op.Install(localdb, db.LocalDB{Path: localdbdir}, pack2ins[cnt], rootdir)
			if err != nil {
				log.Fatal(err)
			}
		}
		if remove {
			p.Printdbg("start remove")
			var inttt int
			var found bool
			p.Printdbg("toremove:", pack2ins[cnt])
			for i := 0; i < len(localtrdb); i++ {
				p.Printdbg(localtrdb[i][0][0], " == ", pack2ins[cnt])
				if localtrdb[i][0][0] == pack2ins[cnt] {
					p.Printdbg("yes found in db")
					inttt = i
					found = true
					break
				}
			}
			if !found {
				log.Fatal("package ", pack2ins[cnt], " not found in db")
			}
			p.Printdbg("db id:", inttt)
			for i := 1; i < len(localtrdb[inttt])-1; i++ {
				path2remove := localtrdb[inttt][i][0]
				path2remove, err := filepath.Abs(filepath.Join(rootdir, path2remove))
				if err != nil {
					log.Fatal("error during removing")
				}
				if !(strings.HasPrefix(path2remove, filepath.Join(rootdir, "/set")) || strings.HasPrefix(path2remove, filepath.Join(rootdir, "./.PACKAGE"))) {
					p.Printdbg(path2remove)
					os.Remove(path2remove)
				}
			}
			p.Printdbg("delete db path", filepath.Join(dbdir, "local", pack2ins[cnt]))
			os.RemoveAll(filepath.Join(dbdir, "local", pack2ins[cnt]))
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
				dbdir, err = filepath.Abs(filepath.Join(rootdir, dbdir))
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

/*func preMergeCheck(fPath string, mList []string) (ok bool) {

}*/

func upperAvailable(st string) bool {
	for _, ru := range st {
		if unicode.IsUpper(ru) {
			return true
		}
	}
	return false
}
