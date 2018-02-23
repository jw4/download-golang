package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type file struct {
	Name string `goquery:".filename a,text"`
	URL  string `goquery:".filename a,[href]"`
	SHA  sha    `goquery:"td tt,text"`
}

func (f *file) check(root string) bool {
	name := filepath.Join(root, f.Name)
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	if !info.Mode().IsRegular() {
		return false
	}

	in, err := os.Open(name)
	if err != nil {
		return false
	}
	defer in.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, in)
	if err != nil {
		return false
	}
	return f.SHA.match(hasher.Sum(nil))
}

func (f *file) get(root string) error {
	tfile, err := ioutil.TempFile(root, f.Name)
	if err != nil {
		return err
	}
	defer os.Remove(tfile.Name())

	res, err := http.Get(f.URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
	default:
		return fmt.Errorf("unexpected response code: %d", res.StatusCode)
	}

	sh := sha256.New()
	tee := io.MultiWriter(tfile, sh)
	_, err = io.Copy(tee, res.Body)
	if err != nil {
		return err
	}

	sha2 := sh.Sum(nil)
	if !f.SHA.match(sha2) {
		return fmt.Errorf("sha256 differs; expected:\n%s got:\n%x", f.SHA, sha2)
	}

	name := filepath.Join(root, f.Name)
	return os.Rename(tfile.Name(), name)
}
