package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	b, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("reading file: %v", err)
	}

	toc, err := mdtoc.Parse(b)
	if err != nil {
		log.Fatalf("parsing file: %v", err)
	}

	if dryRun {
		fmt.Println(toc.String())
		os.Exit(0)
	}

	new, err := mdtoc.Add(b, toc, force)
	if err != nil {
		log.Fatalf("adding toc: %v", err)
	}

	err = overwrite(file, new)
	if err != nil {
		log.Fatalf("writing file: %v", err)
	}
}

// overrwrite trucates an existing file and replaces with the content of b
func overwrite(file string, b []byte) error {
	f, err := os.OpenFile(file, os.O_WRONLY, 0664)
	if err != nil {
		log.Fatalf("opening file: %v", err)
	}
	defer f.Close()

	f.Truncate(0)
	f.Seek(0, 0)
	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
