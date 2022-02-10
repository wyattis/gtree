package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type WalkFunc = func(name string, parts []string, depth int) error

type DirTree struct {
	Name    string
	IsChild bool
	Roots   []*DirTree
}

func (t *DirTree) AddParts(parts []string) {
	if len(parts) == 0 {
		return
	}
	var branch *DirTree
	for _, b := range t.Roots {
		if b.Name == parts[0] {
			branch = b
			break
		}
	}
	if branch == nil {
		branch = new(DirTree)
		branch.Name = parts[0]
		branch.IsChild = true
		t.Roots = append(t.Roots, branch)
	}
	branch.AddParts(parts[1:])
}

func (t *DirTree) AddPath(path string) {
	parts := strings.Split(filepath.ToSlash(filepath.Clean(path)), "/")
	if len(parts) == 0 {
		return
	}
	t.AddParts(parts)
}

func (t *DirTree) walkDepthRecursive(hand WalkFunc, parts []string, depth int) (err error) {
	for _, b := range t.Roots {
		subParts := append([]string{}, parts...)
		subParts = append(subParts, b.Name)
		if err = hand(b.Name, subParts, depth); err != nil {
			return
		}
		if err = b.walkDepthRecursive(hand, subParts, depth+1); err != nil {
			return
		}
	}
	return
}

func (t *DirTree) WalkDepth(hand WalkFunc) error {
	parts := []string{}
	return t.walkDepthRecursive(hand, parts, 0)
}

func (t *DirTree) Sort() {
	sort.Slice(t.Roots, func(a, b int) bool {
		return strings.Compare(t.Roots[a].Name, t.Roots[b].Name) < 0
	})
	for i := range t.Roots {
		t.Roots[i].Sort()
	}
}

func run(path string) (err error) {
	ext := filepath.Ext(path)
	var tree *DirTree
	switch ext {
	case ".zip":
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
		padding := strings.Repeat(" ", depth*2)
		_, err = fmt.Fprintf(w, "%s%s\n", padding, name)
		return
	})
}

func main() {
	flag.Parse()
	for _, f := range flag.Args() {
		if err := run(f); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
}
