package rename

import (
	"bufio"
	"fmt"
	"github.com/manifoldco/promptui"
	flag "github.com/ogier/pflag"
	"io"
	"os"
	"path"
)

type Args struct {
	Files       []string
	Expression  string
	NoAct       bool
	Verbose     bool
	Interactive bool
	Force       bool
	Copy        bool
}

type FromTo struct {
	From string
	To   string
}

func ParseArgs() *Args {
	verbosePtr := flag.BoolP("verbose", "v", false, "Show which files where renamed, if any")
	noActPtr := flag.BoolP("noaction", "n", false, "No renaming, just show what would have been done")
	forcePtr := flag.BoolP("force", "f", false, "Overwrite existing files")
	copyPtr := flag.BoolP("copy", "c", false, "Copy instead of move.")
	helpPtr := flag.BoolP("help", "h", false, "Only show this help text.")
	interactivePtr := flag.BoolP("interactive", "i", false, "Ask for confirmation before renaming each file.")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr,
			"rename - Rename files through a sed-replace expression\n" +
			"Usage:  rename [options] <sed-replace expression> [files...]\n" +
			"  Options:\n" +
			"    -c/--copy:         Copy instead of move.\n" +
			"    -f/--force:        Overwrite existing files.\n" +
			"    -i/--interactive:  Ask for confirmation before renaming each file.\n" +
			"    -n/--noaction:     No changes, just show what would have been done.\n" +
			"    -v/--verbose:      Show which files where renamed, if any.\n" +
			"    -h/--help:         Only show this help text.\n" +
			"  Sed-replace expression:  s/<match>/<replace>/[i][g]\n" +
			"    Match:             Regular expression (tags with round brackets possible).\n" +
			"    Replace:           Replacement, with $0: whole original and $1...: tag.\n" +
			"    i:                 Case insensitive match of regular expression.\n" +
			"    g:                 Global: keep looking for match after first match.\n" +
			"  Files:  If none given, read from stdin.\n")
	}

	flag.Parse()
	l := flag.NArg()
	if l < 1 || *helpPtr { // No arguments
		flag.Usage()
		os.Exit(2)
	}

	var files []string
	if l < 2 { // Only expression: Read files from stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			files = append(files, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil
		}

	} else { // Expression with files
		files = flag.Args()[1:]
	}
	expression := flag.Args()[0]
	return &Args{
		Files:       files,
		Expression:  expression,
		NoAct:       *noActPtr,
		Verbose:     *verbosePtr,
		Copy:        *copyPtr,
		Force:       *forcePtr,
		Interactive: *interactivePtr,
	}
}

func GetReplacements(engine *Engine, args *Args) ([]FromTo, error) {
	destinations := make(map[string]bool)
	var replacements []FromTo
	for _, file := range args.Files {
		dir := path.Dir(file)
		filename := path.Base(file)
		dest, err := engine.Run(filename)
		if err != nil {
			return nil, err
		}

		dest = path.Join(dir, dest)
		if destinations[dest] {
			return nil, fmt.Errorf("Conflicting rename pattern, multiple files will be renamed to the same destination  '%s'", dest)
		}

		destinations[dest] = true
		replacements = append(replacements, FromTo{path.Join(dir, filename), dest})
	}
	return replacements, nil
}

func PrintRename(engine *Engine, fromto FromTo) { // Color match
	from, to, _ := engine.Highlight(fromto.From)
	if from != to {
		fmt.Printf("%s\t-> %s\n", from, to)
	}
}

func Run(args *Args) error {
	engine, err := NewEngine(args.Expression)
	if err != nil {
		return err
	}

	replacements, err := GetReplacements(engine, args)
	if err != nil {
		return err
	}

	if args.Interactive || args.Verbose || args.NoAct {
		for _, fromto := range replacements {
			PrintRename(engine, fromto)
		}
	}

	if args.Interactive {
		prompt := promptui.Prompt{
			Label:     "Continue?",
			IsConfirm: true,
		}
		_, err = prompt.Run()
		if err != nil {
			return nil
		}
	}

	for _, fromto := range replacements {
		if !args.NoAct {
			act := true
			if _, err := os.Stat(fromto.To); err == nil { // File exists
				if !args.Force {
					if fromto.To != fromto.From { // From == To means: no rename
						fmt.Printf("Not overwriting file: '%s'\n", fromto.To)
					}
					act = false
				}
			}
			if act {
				if args.Copy {
					err := copy(fromto.From, fromto.To)
					if err != nil {
						fmt.Printf("Failed to copy file '%s'\n", err)
					}
				} else {
					os.Rename(fromto.From, fromto.To)
					if err != nil {
						fmt.Printf("Failed to rename file '%s'\n", err)
					}
				}
			}
		}
	}
	return nil
}

func copy(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}

	defer src.Close()
	dest, err := os.Create(destination)
	if err != nil {
		return err
	}

	defer dest.Close()
	_, err = io.Copy(dest, src)
	return err
}
