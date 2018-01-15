package main

import (
	"flag"
	"os"
	"os/user"
	"path/filepath"
)

func main() {

	var dest = flag.String("dest", "", "...")
	flag.Parse()

	if *dest == "" {
		usr, err := user.Current()

		if err != nil {
			log.Fatal(err)
		}

		home := usr.HomeDir
		config := filepath.Join(home, ".config")
		root := filepath.Join(config, "secretbox")
		path = filepath.Join(root, "salt")

		*dest = path
	}
}
