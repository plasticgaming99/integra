package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/plasticgaming99/integra/lib/build"
	"github.com/plasticgaming99/integra/lib/pkg/op"
)

type intgOpts struct {
	Install bool `json:"install"`
	Sync    bool `json:"sync"`
	Upgrade bool `json:"upgrade"`
	Remove  bool `json:"remove"`
	Search  bool `json:"search"`
	Query   bool `json:"query"`
	Check   bool `json:"check"`
	AllYes  bool `json:"allyes"`

	RootDir string `json:"rootdir"`
	DbDir   string `json:"dbdir"`

	Quiet bool `json:"quiet"`
	Debug bool `json:"debug"`
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
				intg.Install = true
			case "S":
				intg.Sync = true
			case "U":
				intg.Upgrade = true
			case "R":
				intg.Remove = true
			case "Q":
				intg.Query = true
			case "C":
				intg.Check = true
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

	switch args[0] {
	case "build":
		build.BuildIntegra(args[1:])
		os.Exit(0)
	}

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
			intg.Install = true
		case "-s", "-sync", "--sync":
			intg.Sync = true
		case "-search", "--search":
			intg.Search = true
		case "-u", "-upgrade", "--upgrade":
			intg.Upgrade = true
		case "-r", "-remove", "--remove":
			intg.Remove = true
		case "-q", "-query", "--query":
			intg.Query = true
		case "-c", "-check", "--check":
			intg.Check = true
		case "-y", "-yes", "--yes":
			intg.AllYes = true

		case "-dbg", "--dbg":
			intg.Debug = true
		case "-quiet", "--quiet":
			intg.Quiet = true
		case "-override-root":
			if i < len(args) {
				i++
				if !strings.HasPrefix(args[i], "-") {
					return errors.New("-override-root: option must be a directory")
				}
				intg.RootDir = args[i]
			}
		case "--override-root":
			if ok {
				intg.RootDir = a
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
		Install: false,
		Sync:    false,
		Upgrade: false,
		Remove:  false,
		Search:  false,
		Query:   false,
		Check:   false,

		AllYes: false,

		RootDir: "/",
		DbDir:   "/var/lib/integra/db", // relative path from rootdir

		Quiet: false,
		Debug: false,
	}
	packs := make([]string, 0, 10) // tiny buffer
	err := parseArgs(intg, &packs, args)
	if err != nil {
		log.Fatal(err)
	}
	if !(intg.Install || intg.Sync || intg.Upgrade || intg.Remove || intg.Search) {
		CommandHelp(os.Stdout)
	}
	if intg.Debug {
		j, err := json.MarshalIndent(&intg, "", "    ")
		if err != nil {
			fmt.Println("error marshilizing option")
			fmt.Println(err)
		}
		fmt.Println(string(j))
	}

	if intg.Install {
		for _, s := range packs {
			err := op.Install(s, intg.RootDir)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
