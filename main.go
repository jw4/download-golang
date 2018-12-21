package main

import (
	"log"
	"os"
)

func main() {
	d, err := getGoDownloads()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range d {
		ver := v.Version
		info, err := os.Stat(ver)
		if err != nil {
			if os.IsNotExist(err) {
				if err = os.Mkdir(ver, 0755); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
		if info == nil {
			continue
		}
		if !info.IsDir() {
			log.Fatalf("%q exists and is not a directory", ver)
		}
		for _, f := range v.Files {
			if f.Filename == "" {
				continue
			}
			if !f.check() {
				log.Printf("missing %q; downloading", f.Filename)
				if err = f.get(); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Printf("already have %q with matching hash; skipping", f.Filename)
			}
		}
	}
}
