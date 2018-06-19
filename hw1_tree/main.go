package main

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
)

const FullDirectory int = 0

type ByName []os.FileInfo

func (b ByName) Len() int { return len(b) }

func (b ByName) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

func (b ByName) Less(i, j int) bool { return b[i].Name() < b[j].Name() }

func isDir(name string) (bool, error) {
	stats, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	return stats.IsDir(), nil
}

func writeLevelChars(lvl int, out *os.File, last bool, h []bool) error {
	var char string
	if last {
		char = "└───"
	} else {
		char = "├───"
	}
	for i := 0; i < len(h); i++ {
		if h[i] {
			out.Write([]byte("|   "))
		} else {
			out.Write([]byte("   "))
		}
	}
	out.Write([]byte(char))
	return nil
}

func dirTreeLevel(lvl int, out *os.File, path string, pf bool, hack []bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	hack = append(hack, true)
	b, _ := isDir(path)
	if !b {
		return nil
	}
	entries, err := f.Readdir(FullDirectory)
	if err != nil {
		return err
	}
	eNum := len(entries)

	sort.Sort(ByName(entries))

	for i := 0; i < eNum; i++ {
		subPath := filepath.Join(path, entries[i].Name())
		b, _ := isDir(subPath)
		last := i == eNum-1
		if last {
			hack[lvl] = false
		}
		writeLevelChars(lvl, out, last, hack)
		out.Write([]byte(entries[i].Name() + "\n"))

		if b {
			dirTreeLevel(lvl+1, out, subPath, pf, hack)
		}
	}
	return nil
}

func dirTree(out *os.File, path string, printFiles bool) (e error) {
	hack := make([]bool, 0)
	startingLevel := 0
	if !filepath.IsAbs(path) {
		path, e = filepath.Abs(path)
		if e != nil {
			return e
		}
	}
	b, _ := isDir(path)
	if !b {
		return errors.New("Specify directory")
	}
	if e = os.Chdir(path); e != nil {
		return e
	}
	return dirTreeLevel(startingLevel, out, path, printFiles, hack)
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
