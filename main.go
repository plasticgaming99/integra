package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("integra test")
	parse(os.Args[1:])
}

func parse(in []string) {
	for iv := 0; len(in) > iv; iv++ {
		if strings.Contains(in[iv], "--") {
			fmt.Println("i found the option")
		}
	}
}
