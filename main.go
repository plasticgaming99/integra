package main

import (
	"os"

	"github.com/plasticgaming99/integra/cmd"
)

func main() {
	cmd.Execute(os.Args[1:])
}
