package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jar-b/mdtoc"
)

var (
	force, dryRun, withTocHeading, version bool

	out        string
	tocHeading string = mdtoc.DefaultTocHeading
)

func main() {
	// slightly better usage output
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "Generate a table of contents for an existing markdown document.\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] [filename]\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.BoolVar(&force, "force", false, "force overwrite of existing contents (optional)")
	flag.BoolVar(&dryRun, "dry-run", false, "print generated contents, but do not write to file (optional)")
	flag.StringVar(&out, "out", "", "output file (optional, defaults to adding to source file)")
	flag.BoolVar(&withTocHeading, "with-toc-heading", false, "include a heading with the generated contents (optional)")
	flag.StringVar(&tocHeading, "toc-heading", tocHeading, "contents heading (-with-toc-heading must be specified)")
	flag.BoolVar(&version, "version", false, "display version")
	flag.Parse()

	if version {
		fmt.Print(mdtoc.Version)
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		log.Fatal("unexpected number of args")
	}
	in := flag.Arg(0)

	b, err := os.ReadFile(in)
	if err != nil {
		log.Fatalf("reading file: %v", err)
	}

	if dryRun {
		toc, err := mdtoc.New(b)
		if err != nil {
			log.Fatalf("parsing file: %v", err)
		}
		fmt.Println(toc.String())
		os.Exit(0)
	}

	cfg := mdtoc.Config{
		Force:          force,
		WithTocHeading: withTocHeading,
		TocHeading:     tocHeading,
	}
	withToc, err := mdtoc.Insert(b, &cfg)
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

	err = os.WriteFile(target, withToc, 0644)
	if err != nil {
		log.Fatalf("writing file: %v", err)
	}
}
