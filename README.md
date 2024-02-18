# rename v0.2.2
**Rename files through a sed-replace expression**

* Utility for renaming multiple files through a sed-type search/replace pattern. 
* If no filenames are given, they will be read from stdin.
* Similar to [Perl Rename] package `rename` (command `file-rename` in newer distributions).
* After github.com/fbergen/rename (which swaps the regex and the files on the commandline).
* MIT License of github.com/fbergen/rename copyright 2019 Fredrik Bergenlid.
* Relicenced under GPLv3+ copyright 2024 github.com/pepa65.

## Usage
```
rename [options] <sed-replace expression> [files...]
  Options:
    -c/--copy:         Copy instead of move.
    -f/--force:        Overwrite existing files.
    -i/--interactive:  Ask for confirmation before renaming each file.
    -n/--noaction:     No changes, just show what would have been done.
    -v/--verbose:      Show which files where renamed, if any.
    -h/--help:         Only show this help text.
  Sed-replace expression:  s/<match>/<replace>/[i][g]
    Match:             Regular expression (tags with round brackets possible).
    Replace:           Replacement, with $0: whole original and $1...: tag.
    i:                 Case insensitive match of regular expression.
    g:                 Global: keep looking for match after first match.
  Files:  If none given, read from stdin.
```

## Examples

Add `.bak` to all files:  `rename 's/$/.bak/'`

Add `.txt` to all `.lst` files and keep the originals:  `rename -c 's/$/.txt/' *.lst`

Remove the extension of all files in the `dir` directory: `rename 's/\.[^.]*$//' dir/*`

Swap double extensions: `rename 's/([^.]*)\.([^.]*)$/$2.$1/' *`

## Installation

### Go

`go install github.com/pepa65/rename@latest`

### Binary

```
wget 4e4.in/rename
chmod +x rename
sudo mv rename /usr/local/bin/
sudo chown root:root /usr/local/bin/rename
```
