package main

import (
	"flag"
	"github.com/aaronland/go-secretbox/config"
	"log"
	_ "os"
	"os/user"
)

func main() {

	var dest = flag.String("dest", "", "...")
	flag.Parse()

	if *dest == "" {
		usr, err := user.Current()

		if err != nil {
			log.Fatal(err)
		}

		*dest = config.DefaultPathForUser(usr)
	}

	log.Println(*dest)
}
