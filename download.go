package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func getGoDownloads() (download, error) {
	client := &http.Client{Timeout: time.Second * 30}
	res, err := client.Get("https://golang.org/dl/?mode=json&include=all")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var d download
	err = json.NewDecoder(res.Body).Decode(&d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

type version struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []file `json:"files"`
}

type download []version
