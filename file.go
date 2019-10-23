package main

import (
	"crypto"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"time"
)

var (
	staticBuffer = make([]byte, 1<<16)
	jar, _       = cookiejar.New(nil)
	client       = &http.Client{Timeout: time.Second * 600, Jar: jar}
)

type file struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	SHA256   sha    `json:"sha256"`
	Size     int    `json:"size"`
	Kind     string `json:"kind"`
}

func (f file) check() bool {
	name := filepath.Join(f.Version, f.Filename)
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	if !info.Mode().IsRegular() {
		return false
	}
	if f.SHA256.match([]byte{}) {
		return true
	}
	in, err := os.Open(name)
	if err != nil {
		return false
	}
	defer in.Close()

	hasher := crypto.SHA256.New()
	if _, err = io.CopyBuffer(hasher, in, staticBuffer); err != nil {
		return false
	}
	return f.SHA256.match(hasher.Sum(nil))
}

func (f file) get() error {
	temp, err := ioutil.TempFile(f.Version, f.Filename)
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	defer temp.Close()

	response, err := client.Get(fmt.Sprintf("https://golang.org/dl/%s", f.Filename))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
	default:
		return fmt.Errorf("unexpected response code: %d", response.StatusCode)
	}

	hasher := crypto.SHA256.New()
	if _, err = io.CopyBuffer(io.MultiWriter(temp, hasher), response.Body, staticBuffer); err != nil {
		return err
	}

	if !f.SHA256.match([]byte{}) {
		sum := hasher.Sum(nil)
		if !f.SHA256.match(sum) {
			return fmt.Errorf("sha256 differs; expected:\n%s got:\n%x", f.SHA256, sum)
		}
	}

	if err = temp.Chmod(0644); err != nil {
		return err
	}

	return os.Rename(temp.Name(), filepath.Join(f.Version, f.Filename))
}
