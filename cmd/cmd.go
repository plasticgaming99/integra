package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

type intgOpts struct {
	install bool
	sync    bool
	upgrade bool
	remove  bool
	search  bool
	allyes  bool

	rootdir string
	dbdir   string

	quiet bool
	debug bool
}

func upperAvailable(st string) bool {
	for _, ru := range st {
		if unicode.IsUpper(ru) {
			return true
		}
	}
	return false
}

// parse integra style opt
func parseOp(input string, intg *intgOpts) {
	run := []rune(input)
	for i := 0; i < len(run); i++ {
		if unicode.IsUpper(run[i]) {
			switch string(run[i]) {
			case "I":
				intg.install = true
			case "S":
				intg.sync = true
			case "U":
				intg.upgrade = true
			case "R":
				intg.remove = true
			}
		}
	}
}

func parseArgs(intg *intgOpts, packs *[]string, args []string) error {
	// yeah i can't parse nothing
	if len(args) == 0 {
		return nil
	}

	i := 0
	if upperAvailable(args[i]) {
		parseOp(args[i], intg)
		i++
	}

	var b, a, ok = "", "", false
	for ; i < len(args); i++ {
		if args[i] == "" {
			// sad
			continue
		}

		if !strings.HasPrefix(args[i], "-") {
			*packs = append(*packs, args[i])
			continue
		}

		b, a, ok = strings.Cut(args[i], "=")
		switch b {
		case "-i", "-install", "--install":
			intg.install = true
		case "-s", "-sync", "--sync":
			intg.sync = true
		case "-search", "--search":
			intg.search = true
		case "-u", "-upgrade", "--upgrade":
			intg.upgrade = true
		case "-r", "-remove", "--remove":
			intg.remove = true
		case "-y", "-yes", "--yes":
			intg.allyes = true

		case "-dbg", "--dbg":
			intg.debug = true
		case "-quiet", "--quiet":
			intg.quiet = true
		case "-override-root":
			if i < len(args) {
				i++
				if !strings.HasPrefix(args[i], "-") {
					return errors.New("-override-root: option must be a directory")
				}
				intg.rootdir = args[i]
			}
		case "--override-root":
			if ok {
				intg.rootdir = a
			} else {
				return errors.New("--override-root: option was not passed")
			}
		}
	}
	return nil
}

func Execute(args []string) {
	fmt.Println("integra")
	intg := &intgOpts{
		install: false,
		sync:    false,
		upgrade: false,
		remove:  false,
		search:  false,

		allyes: false,

		rootdir: "/",
		dbdir:   "/var/lib/integra/db",

		quiet: false,
		debug: false,
	}
	packs := make([]string, 0, 10) // tiny buffer
	err := parseArgs(intg, &packs, args)
	if err != nil {
		log.Fatal(err)
	}
	if !(intg.install || intg.sync || intg.upgrade || intg.remove || intg.search) {
		CommandHelp(os.Stdout)
	}
}
