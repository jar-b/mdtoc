package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jar-b/mdtoc"
)

var (
	force, dryRun bool
	out           string
)

func init() {
	// slightly better usage output
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] [filename]\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.BoolVar(&force, "force", false, "force overwrite of existing contents (optional)")
	flag.BoolVar(&dryRun, "dry-run", false, "print generated contents, but do not write to file (optional)")
	flag.StringVar(&out, "out", "", "output file (optional, defaults to adding to source file)")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("unexpected number of args")
	}
	in := flag.Arg(0)

	b, err := os.ReadFile(in)
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
		if err == mdtoc.ErrExistingToc {
			log.Fatalf("%s. Use the -force flag to force overwrite.\n", err.Error())
		}
		log.Fatalf("adding toc: %v", err)
	}

	target := in
	if out != "" {
		target = out
	}

	err = os.WriteFile(target, new, 0644)
	if err != nil {
		log.Fatalf("writing file: %v", err)
	}
}
