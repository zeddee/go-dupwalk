# Simple duplicate file finder

`dewalk` walks the current directory tree
starting to find duplicate files.

It hashes the first 5MB/5120 bytes of each file
in the directory tree and compares it to subsequent
files it finds.
When it finishes walking the directory tree,
the application, it prints to STDOUT a list
of all found duplicates as JSON in this format:

```
[{"Original":"test/data/dir.log","Duplicate":"test/data/dir_duplicate.log"}]
```
## Usage

```
Usage:
  -exclude string
    	Filter out files that match these filenames
  -ext string
    	Filter out these file extensions.
  -list
    	List files only.
  -out string
    	Write list of duplicates as JSON to this file.
  -root string
    	Start walking from this directory. (default ".")
  -size int
    	Minimum size of files to display.
  -v	Prints to output each time an item is found
```

