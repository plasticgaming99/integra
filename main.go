package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	sync     = false
	search   = false
	refresh  = false
	refforce = false
	upgrade  = false
	allyes   = false

	debug = false

	pack2ins []string
)

func main() {
	fmt.Println("integra test")
	parse(os.Args[1:])
	printdbg(sync, search, refresh, refforce, upgrade, allyes)
	printdbg(pack2ins)
}

func parse(in []string) {
	for iv := 0; len(in) > iv; iv++ {
		if strings.Contains(in[iv], "--") {
			fmt.Println("i found the option yay")
			//check
			switch in[iv][2:] {
			case "dbg":
				debug = true
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
		} else
		//check combined option like pacman
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
