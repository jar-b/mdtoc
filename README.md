# mdtoc
[![Go Reference](https://pkg.go.dev/badge/github.com/jar-b/mdtoc.svg)](https://pkg.go.dev/github.com/jar-b/mdtoc)
[![Go Report Card](https://goreportcard.com/badge/github.com/jar-b/mdtoc)](https://goreportcard.com/report/github.com/jar-b/mdtoc)
[![go-build](https://github.com/jar-b/mdtoc/actions/workflows/go.yml/badge.svg)](https://github.com/jar-b/mdtoc/actions/workflows/go.yml)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/jar-b/mdtoc)

Generate a table of contents for an existing markdown document. The table of contents will link to anchor tags, and preserve the level of nesting. 


<!---mdtoc begin--->
* [CLI](#cli)
  * [Installation](#installation)
  * [Usage](#usage)
* [Library](#library)
  * [Usage](#usage-1)
<!---mdtoc end--->
## CLI

### Installation

Via `go get`:

```sh
go get github.com/jar-b/mdtoc/cmd/mdtoc
```

### Usage

```
$ mdtoc -h
Usage: mdtoc [flags] [filename]

Flags:
  -dry-run
    	print generated contents, but do not write to file (optional)
  -force
    	force overwrite of existing contents (optional)
  -out string
        output file (optional, defaults to adding to source file)
```

Examples:

```sh
# add new
mdtoc mydoc.md

# dry run
mdtoc -dry-run mydoc.md

# force overwrite of existing
mdtoc -force mydoc.md
```

## Library

`import github.com/jar-b/mdtoc`

### Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/jar-b/mdtoc"
)

func main() {
	b := []byte("# Title\ndescription text\n\n## Heading 1\ntext\n## Heading 2\nmore text")

	toc, err := mdtoc.Parse(b)
	if err != nil {
		log.Fatal(err)
	}

	new, err := mdtoc.Add(b, toc, false)
	if err != nil {
		log.Fatal(err)
	}

	// do something with `new`
	fmt.Println(string(new))
}
```
