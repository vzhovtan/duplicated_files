package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type fentry struct {
	fpath string
	fhash string
}

// fList type and assosiated methods created to use sort.Sort based on file MD5 hash
type fList []fentry

func (f fList) Len() int {
	return len(f)
}

func (f fList) Less(i, j int) bool {
	return f[i].fhash > f[j].fhash
}

func (f fList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func main() {
	dirStart := flag.String("dir", ".", "directory to start search for duplicated files recursively")
	flag.Parse()

	root := *dirStart
	fileList, err := getStat(root)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", *fileList)

	dupList, err := findDup(fileList)
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(fList(*dupList))
	printOut(os.Stdout, dupList)
}

//func printOut sends the output to writer
func printOut(w io.Writer, d *[]fentry) {
	for _, file := range *d {
		fmt.Fprintf(w, "Path: %s ==> MD5 %s\n", file.fpath, file.fhash)
	}
}

// findDup function keep only duplicated entries based on MD5 hash for the file
func findDup(flist *[]fentry) (*[]fentry, error) {
	origList := &[]fentry{}
	dupList := &[]fentry{}
	key := map[string]struct{}{}
	for _, v := range *flist {
		if _, ok := key[v.fhash]; !ok {
			key[v.fhash] = struct{}{}
			*origList = append(*origList, v)
		} else {
			*dupList = append(*dupList, v)
		}
	}
	for _, or := range *origList {
		for _, dup := range *dupList {
			if or.fhash == dup.fhash {
				*dupList = append(*dupList, or)
				break
			}
		}
	}
	return dupList, nil
}

// getStat function build slice of file entries starting from provided root directory
func getStat(root string) (*[]fentry, error) {
	flist := &[]fentry{}
	fileSystem := os.DirFS(root)
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasPrefix(path, ".") {
			info, err := d.Info()
			if err != nil {
				return err
			}
			if !info.IsDir() {
				entry := fentry{
					fpath: path,
					fhash: getHash(root, path),
				}
				*flist = append(*flist, entry)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return flist, nil
}

// getHash function calculate MD5 hash sum for the file and return it as string
func getHash(root, fname string) string {
	absPath := filepath.Join(root, fname)
	f, err := os.Open(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := md5.New()

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
