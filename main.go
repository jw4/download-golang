package main

import (
	"log"
	"os"
)

func main() {
	d, err := allVersions()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range d.Versions {
		ver := v.Version[:len(v.Version)-4]
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
			if f.Name == "" {
				continue
			}
			if !f.check(ver) {
				log.Printf("missing %q; downloading", f.Name)
				if err = f.get(ver); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Printf("already have %q with matching hash; skipping", f.Name)
			}
		}
	}
}
