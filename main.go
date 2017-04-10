package main

import (
	"fmt"
	"log"
	"os"

	"github.com/maxmansfield/fixcsv/fixer"
)

func main() {
	usage := "USAGE: " + os.Args[0] + " [IN] [OUT]"

	if len(os.Args) != 3 {
		fmt.Println(usage)
		os.Exit(0)
	}
	f, err := fixer.NewFixer(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	f.Write(os.Args[2])
	return
}
