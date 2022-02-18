package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var block = NewSet()
var format string
var tab string

var formatMap = map[string]string{
	"zip":   "zip",
	"apk":   "zip",
	"docx":  "zip",
	"xlsx":  "zip",
	"ppsx":  "zip",
	"pptx":  "zip",
	"thmx":  "zip",
	"pk3":   "zip",
	"pk4":   "zip",
	"usdz":  "zip",
	"xpi":   "zip",
	"mgz":   "zip",
	"smzip": "zip",

	"gzip":     "gzip",
	"gz":       "gzip",
	"tgz":      "gzip",
	"gnumeric": "gzip",
	"adz":      "gzip",
	"maf":      "gzip",

	"dir": "dir",
}

func run(path string) (err error) {
	if format == "" {
		ext := strings.TrimSpace(filepath.Ext(path))
		if ext == "" || ext == "." {
			format = "dir"
		} else {
			format = ext[1:]
		}
	}
	var tree *DirTree
	f, ok := formatMap[format]
	if !ok {
		return fmt.Errorf("unsupported format: %s", format)
	}
	switch f {
	case "zip":
		tree, err = doZip(path)
	default:
		tree, err = doDirectory(path)
	}
	if err != nil {
		return
	}
	return printTree(os.Stdout, tree)
}

func doDirectory(path string) (tree *DirTree, err error) {
	tree = new(DirTree)
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		tree.AddPath(path)
		return nil
	})
	return
}

func doZip(path string) (tree *DirTree, err error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return
	}
	tree = new(DirTree)
	for _, f := range r.File {
		tree.AddPath(f.Name)
	}
	return
}

func printTree(w io.Writer, tree *DirTree) (err error) {
	tree.Sort()
	return tree.WalkDepth(func(name string, parts []string, depth int) (err error) {
		padding := strings.Repeat(tab, depth)
		_, err = fmt.Fprintf(w, "%s%s\n", padding, name)
		return
	})
}

func main() {
	flag.StringVar(&format, "f", "", "specify the format for the input (dir,zip,gzip)")
	flag.Var(block, "b", "specify paths to block")
	flag.StringVar(&tab, "tab", "  ", "what characters to use as a tab for each level")
	flag.Parse()
	for _, f := range flag.Args() {
		if err := run(f); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
}
