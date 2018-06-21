package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

const FullDirectory int = 0

type ByName []os.FileInfo

func (b ByName) Len() int { return len(b) }

func (b ByName) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

func (b ByName) Less(i, j int) bool { return b[i].Name() < b[j].Name() }

func isDir(path string) bool {
	stats, e := os.Stat(path)
	if e != nil {
		fmt.Println(e)
	}
	return stats.IsDir()
}

func formatEntryString(info os.FileInfo, last bool, l []bool) (entry string) {
	for i := 0; i < len(l); i++ {
		if l[i] {
			entry += "\t"
		} else {
			entry += "│\t"
		}
	}
	if last {
		entry += "└───"
	} else {
		entry += "├───"
	}
	entry += info.Name()
	if !info.IsDir() {
		switch info.Size() == 0 {
		case true:
			entry += " (empty)"
		case false:
			entry += fmt.Sprintf(" (%db)", info.Size())
		}
	}
	entry += "\n"
	return entry
}

func treeRecursive(out io.Writer, path string, printFiles bool, lastitudes []bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	entries, err := f.Readdir(FullDirectory)
	if err != nil {
		return err
	}

	if !printFiles {
		for i := len(entries) - 1; i >= 0; i-- {
			if !entries[i].IsDir() {
				entries = append(entries[:i], entries[i+1:]...)
			}
		}
	}

	sort.Sort(ByName(entries))

	eNum := len(entries)
	for i := 0; i < eNum; i++ {
		newL := make([]bool, len(lastitudes))
		copy(newL, lastitudes)
		subPath := filepath.Join(path, entries[i].Name())
		b := entries[i].IsDir()
		last := i == eNum-1
		name := formatEntryString(entries[i], last, newL)
		//out.Write([]byte(name))
		fmt.Fprint(out, name)
		if b {
			if last {
				newL = append(newL, true)
			} else {
				newL = append(newL, false)
			}
			treeRecursive(out, subPath, printFiles, newL)
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) (e error) {
	// lastitudes[i] == true when parent i of a child i+1 is the last folder on i level
	lastitudes := make([]bool, 0)
	if !filepath.IsAbs(path) {
		path, e = filepath.Abs(path)
		if e != nil {
			return e
		}
	}
	b := isDir(path)
	if !b {
		return errors.New("Specify directory")
	}
	return treeRecursive(out, path, printFiles, lastitudes)
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
