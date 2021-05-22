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

	f, err := os.OpenFile(file, os.O_RDWR, 0664)
	if err != nil {
		log.Fatalf("opening file: %v", err)
	}
	defer f.Close()

	if dryRun {
		toc, err := mdtoc.Parse(f)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(toc.String())
	} else {
		b, err := mdtoc.Add(f, force)
		if err != nil {
			log.Fatal(err)
		}

		err = overwrite(f, b)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// overrwrite trucates an existing file and replaces with the content of b
func overwrite(f *os.File, b []byte) error {
	_, err := f.Seek(0, 0)
	if err != nil {
		return err
	}

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
