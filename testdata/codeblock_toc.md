# Title
An example document with code blocks containing `#` characters to verify these are not included as TOC items.

<!---mdtoc begin--->
* [Heading 1](#heading-1)
* [Heading 2](#heading-2)
* [Heading 3](#heading-3)
<!---mdtoc end--->
## Heading 1

```
## plaintext
### lines
```

## Heading 2

```sh
## shell
### lines
#### with
##### comments
echo "skip this"
```

## Heading 3

```python
## python
### lines
#### with
##### comments

if __name__ == '__main__':
    print('skip this')
```

