# mdtoc
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/jar-b/mdtoc)
[![build](https://github.com/jar-b/mdtoc/actions/workflows/build.yml/badge.svg)](https://github.com/jar-b/mdtoc/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jar-b/mdtoc)](https://goreportcard.com/report/github.com/jar-b/mdtoc)
[![Go Reference](https://pkg.go.dev/badge/github.com/jar-b/mdtoc.svg)](https://pkg.go.dev/github.com/jar-b/mdtoc)

Generate a table of contents for an existing markdown document. The table of contents will link to anchor tags, and preserve the level of nesting. 


<!--mdtoc: begin-->
* [CLI](#cli)
  * [Installation](#installation)
  * [Usage](#usage)
  * [Examples](#examples)
* [Library](#library)
  * [Usage](#usage-1)
<!--mdtoc: end-->
## CLI

### Installation

Via `go install`:

```sh
go install github.com/jar-b/mdtoc/cmd/mdtoc@latest
```

### Usage

```
$ mdtoc -h
Generate a table of contents for an existing markdown document.

Usage: mdtoc [flags] [filename]

Flags:
  -dry-run
        print generated contents, but do not write to file (optional)
  -force
        force overwrite of existing contents (optional)
  -out string
        output file (optional, defaults to adding to source file)
  -toc-heading string
        contents heading (-with-toc-heading must be specified) (default "Table of Contents")
  -version
        display version
  -with-toc-heading
        include a heading with the generated contents (optional)
```

### Examples

```sh
# add new
mdtoc mydoc.md

# dry run
mdtoc -dry-run mydoc.md

# force overwrite of existing
mdtoc -force mydoc.md

# redirect output to new document
mdtoc -out other.md mydoc.md

# with custom heading
mdtoc -with-toc-heading -toc-heading "document stuff" mydoc.md
```

## Library

`import github.com/jar-b/mdtoc`

### Usage

```go
package main

import (
        "fmt"

        "github.com/jar-b/mdtoc"
)

func main() {
        b := []byte("# Title\ndescription text\n\n## Heading 1\ntext\n## Heading 2\nmore text")

        // extract just the proposed TOC ("dry-run")
        toc, _ := mdtoc.New(b)
        fmt.Println(toc.String())

        // OR insert TOC into an existing document
        out, _ := mdtoc.Insert(b, mdtoc.DefaultConfig)
        fmt.Println(string(out))
}
```
