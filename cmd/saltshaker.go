package main

import (
	"flag"
	"github.com/aaronland/go-secretbox/config"
	"github.com/aaronland/go-secretbox/salt"
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

	opts := salt.DefaultSaltOptions()
	s, err := salt.NewRandomSalt(opts)

	log.Println(s, err)
}
