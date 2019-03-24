package replaceunzip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

var garbage = regexp.MustCompile(`(?:^__MACOSX/|/\.DS_Store$)`)

func hasRoot(ms []*zip.File) bool {
	names := make([]string, 0, len(ms))
	for _, f := range ms {
		if garbage.MatchString(f.Name) {
			continue
		}
		names = append(names, f.Name)
	}
	sort.Slice(names, func(i, j int) bool { return len(names[i]) < len(names[j]) })

	if !strings.HasSuffix(names[0], "/") {
		return false
	}
	root := names[0]

	for _, member := range names[1:] {
		if !strings.HasPrefix(member, root) {
			return false
		}
	}
	return true
}

func genRoot(zipName string) string {
	file := filepath.Base(zipName)

	ext := filepath.Ext(file)
	if ext == "" {
		return file
	}
	return strings.TrimSuffix(file, ext)
}

// Run acts like unzip, but leaves out .DS_Store and __MACOSX files,
// and puts all of the zip contents in a single root directory if they were not
// already.
func Run(args []string, _ io.Reader) error {
	if len(args) != 2 {
		fmt.Println("Usage:", args[0], "some-zip-file.zip")
		os.Exit(1)
	}

	zipName := args[1]

	fmt.Println("Archive:", zipName)
	r, err := zip.OpenReader(zipName)
	if err != nil {
		return errors.Wrap(err, "Couldn't open zip file")
	}
	defer r.Close()
	var root string
	if !hasRoot(r.File) {
		root = genRoot(zipName)
	}

	ms, err := sanitize(root, r.File)
	if err != nil {
		return err
	}

	for _, f := range ms {
		if err := extractMember(f); err != nil {
			fmt.Printf("  inflating: %s\n", f.Name)

			return errors.Wrap(err, "extractMember")
		}
	}
	return nil
}

// sanitize injects a root into all members if they do not share one, filters
// out garbage files, and errors if any of the members have ".." in the name
func sanitize(root string, ms []*zip.File) ([]*zip.File, error) {
	ret := make([]*zip.File, 0, len(ms))

	for _, m := range ms {
		if garbage.MatchString(m.Name) {
			continue
		}
		segments := strings.Split(m.Name, "/")
		for _, s := range segments {
			fmt.Println(s)
			if s == ".." {
				return nil, errors.New(".. not allowed in member name (Name=" + m.Name + ")")
			}
		}
		m.Name = filepath.Join(append([]string{root}, segments...)...)
		ret = append(ret, m)
	}

	return ret, nil
}

func extractMember(f *zip.File) error {
	if f.FileInfo().IsDir() {
		return errors.Wrap(os.Mkdir(f.Name, os.FileMode(0755)), "os.Mkdir")
	}

	rc, err := f.Open()
	if err != nil {
		return errors.Wrap(err, "Couldn't open zip file member")
	}
	defer rc.Close()

	dir := filepath.Dir(f.Name)
	err = os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't create directory to extract to: %s", err)
		return nil
	}

	file, err := os.Create(f.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't create file to extract to: %s", err)
		return nil
	}

	_, err = io.Copy(file, rc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't copy zip file member (%s): %s", f.Name, err)
	}

	err = file.Chmod(f.FileInfo().Mode())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't chown extracted file: %s", err)
	}

	err = file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't close extracted file: %s", err)
	}

	return nil
}
