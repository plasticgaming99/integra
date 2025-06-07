package printdbg

import (
	"fmt"
)

var (
	dbg = false
)

type PrintDbg struct {
	Debug bool
}

func Dbg() {
	dbg = true
}

func Printdbg(a ...any) {
	if dbg {
		fmt.Println(a...)
	}
}
