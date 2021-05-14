package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jar-b/mdtoc"
)

var (
	file          string
	force, dryRun bool
)

func main() {
	flag.StringVar(&file, "file", "", "file to add contents (required)")
	flag.BoolVar(&force, "force", false, "force overwrite of existing contents (optional)")
	flag.BoolVar(&dryRun, "dry-run", false, "print generated contents, but do not write to file (optional)")

	flag.Parse()
	if file == "" {
		log.Fatal("-file flag is required")
	}

	if dryRun {
		toc, err := mdtoc.ParseToc(file)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(toc.String())
	}
}
