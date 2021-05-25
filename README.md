# mdtoc
Generate a table of contents for an existing markdown document. The table of contents will link to anchor tags, and preserve the level of nesting. 

<!---mdtoc begin--->
* [Installation](#installation)
* [Usage](#usage)
  * [Examples](#examples)
<!---mdtoc end--->
## Installation

Via `go get`:

```sh
go get github.com/jar-b/mdtoc/cmd/mdtoc
```

## Usage

```
$ mdtoc -h
Usage of mdtoc:
  -dry-run
        print generated contents, but do not write to file (optional)
  -force
        force overwrite of existing contents (optional)
```

### Examples

```sh
# add new
mdtoc mydoc.md

# dry run
mdtoc -dry-run mydoc.md

# force overwrite of existing
mdtoc -force mydoc.md
```

