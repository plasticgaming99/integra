package main

import (
	"os"

	cmd "github.com/plasticgaming99/integra/cmd/integra"
)

func main() {
	cmd.Execute(os.Args[1:])
}
