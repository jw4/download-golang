package main

import (
	"crypto"
	"net/http"

	"astuart.co/goq"
)

type version struct {
	Version string `goquery:"h2,text"`
	Hash    string `goquery:"table.codetable thead tr.first th:nth-child(6),text"`
	Files   []file `goquery:"table.codetable tbody tr"`
}

func (v *version) ver() string { return v.Version[:len(v.Version)-4] }

func (v *version) hash() crypto.Hash {
	switch v.Hash {
	case "SHA1 Checksum":
		return crypto.SHA1
	case "SHA256 Checksum":
		return crypto.SHA256
	default:
		return crypto.SHA256
	}
}

type download struct {
	Versions []version `goquery:"div[id^='go1.'] .expanded"`
}

func getGoDownloads() (*download, error) {
	res, err := http.Get("https://golang.org/dl/")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	d := &download{}

	err = goq.NewDecoder(res.Body).Decode(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
