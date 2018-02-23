package main

import (
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
	Sum  sha    `goquery:"td tt,text"`
}

func (f *file) check(v *version) bool {
	name := filepath.Join(v.ver(), f.Name)
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

	hasher := v.hash().New()
	if _, err = io.Copy(hasher, in); err != nil {
		return false
	}
	return f.Sum.match(hasher.Sum(nil))
}

func (f *file) get(v *version) error {
	temp, err := ioutil.TempFile(v.ver(), f.Name)
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	defer temp.Close()

	response, err := http.Get(f.URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	default:
		return fmt.Errorf("unexpected response code: %d", response.StatusCode)
	}

	hasher := v.hash().New()
	if _, err = io.Copy(io.MultiWriter(temp, hasher), response.Body); err != nil {
		return err
	}

	sum := hasher.Sum(nil)
	if !f.Sum.match(sum) {
		return fmt.Errorf("sha256 differs; expected:\n%s got:\n%x", f.Sum, sum)
	}

	if err = temp.Chmod(0644); err != nil {
		return err
	}

	return os.Rename(temp.Name(), filepath.Join(v.ver(), f.Name))
}
