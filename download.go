package main

import (
	"net/http"

	"astuart.co/goq"
)

type version struct {
	Version string `goquery:"h2,text"`
	Files   []file `goquery:"table.codetable tbody tr"`
}

type download struct {
	Versions []version `goquery:"div[id^='go1.'] .expanded"`
}

func allVersions() (*download, error) {
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
