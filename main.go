// artex package manager
// integra (much unstable such wow)
// goal: full dependency-based parallel
//   package management
// now:
//   just simple package tool

package main

import (
	"os"

	"github.com/plasticgaming99/integra/cmd"
)

func main() {
	cmd.Exec(os.Args)
}
