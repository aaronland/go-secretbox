package main

// https://godoc.org/golang.org/x/crypto/nacl/secretbox

// please don't import anything that isn't part of the standard
// library or "golang.org/x/" unless there's a really good reason
// to (20171025/thisisaaronland)

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"github.com/aaronland/go-secretbox/config"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Secretbox struct {
	Key    [32]byte
	Suffix string
	Debug  bool
}

func (sb Secretbox) LockFile(abs_path string) (string, error) {

	root := filepath.Dir(abs_path)
	fname := filepath.Base(abs_path)

	var nonce [24]byte

	_, err := io.ReadFull(rand.Reader, nonce[:])

	if err != nil {
		return "", err
	}

	body, err := ReadFile(abs_path)

	if err != nil {
		return "", err
	}

	enc := secretbox.Seal(nonce[:], body, &nonce, &sb.Key)
	enc_hex := base64.StdEncoding.EncodeToString(enc)

	enc_fname := fmt.Sprintf("%s%s", fname, sb.Suffix)
	enc_path := filepath.Join(root, enc_fname)

	if sb.Debug {
		log.Printf("debugging is enabled so don't actually write %s\n", enc_path)
		return enc_path, nil
	}

	return WriteFile([]byte(enc_hex), enc_path)
}

func (sb Secretbox) UnlockFile(abs_path string) (string, error) {

	root := filepath.Dir(abs_path)
	fname := filepath.Base(abs_path)
	ext := filepath.Ext(abs_path)

	if ext != sb.Suffix {
		return "", errors.New("Unexpected suffix")
	}

	body_hex, err := ReadFile(abs_path)

	if err != nil {
		return "", err
	}

	body_str, err := base64.StdEncoding.DecodeString(string(body_hex))

	if err != nil {
		return "", err
	}

	body := []byte(body_str)

	var nonce [24]byte
	copy(nonce[:], body[:24])

	out, ok := secretbox.Open(nil, body[24:], &nonce, &sb.Key)

	if !ok {
		msg := fmt.Sprintf("Failed to unlock %s", abs_path)
		return "", errors.New(msg)
	}

	out_fname := strings.TrimRight(fname, ext)
	out_path := filepath.Join(root, out_fname)

	if sb.Debug {
		log.Printf("debugging is enabled so don't actually write %s\n", out_path)
		return out_path, nil
	}

	return WriteFile(out, out_path)
}

func ReadFile(path string) ([]byte, error) {

	fh, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(fh)
}

func WriteFile(body []byte, path string) (string, error) {

	fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return "", err
	}

	_, err = fh.Write(body)

	if err != nil {
		return "", err
	}

	err = fh.Close()

	if err != nil {
		return "", err
	}

	return path, nil
}

func main() {

	var suffix = flag.String("suffix", ".enc", "...")
	var unlock = flag.Bool("unlock", false, "Decrypt files.")
	var debug = flag.Bool("debug", false, "...")
	var salt = flag.String("salt", "config:", "...")

	flag.Parse()

	possible := flag.Args()
	files := make([]string, 0)

	if len(possible) == 0 {
		log.Println("No secrets to tell!")
		os.Exit(0)
	}

	for _, path := range possible {

		abs_path, err := filepath.Abs(path)

		if err != nil {
			log.Fatal(err)
		}

		files = append(files, abs_path)
	}

	if *salt == "env:" {
		*salt = os.Getenv("SECRETBOX_SALT")
	} else if strings.HasPrefix(*salt, "config:") {

		parts := strings.Split(*salt, ":")
		var path string

		if parts[1] == "" {
			usr, err := user.Current()

			if err != nil {
				log.Fatal(err)
			}

			path = config.DefaultPathForUser(usr)

		} else {
			path = parts[1]
		}

		abs_path, err := filepath.Abs(path)
		// log.Println("SALT", abs_path)

		if err != nil {
			log.Fatal(err)
		}

		log.Println(abs_path)

		_, err = os.Stat(abs_path)

		if err != nil {
			log.Fatal(err)
		}

		fh, err := os.Open(abs_path)

		if err != nil {
			log.Fatal(err)
		}

		defer fh.Close()

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			log.Fatal(err)
		}

		*salt = string(body)

	} else if *salt != "" {
		// pass
	} else {
		log.Fatal("missing salt")
	}

	if len(*salt) < 8 {
		log.Fatal("invalid salt")
	}

	fmt.Println("enter password: ")
	pswd, err := terminal.ReadPassword(0)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("enter password (again): ")
	pswd2, err := terminal.ReadPassword(0)

	if err != nil {
		log.Fatal(err)
	}

	if string(pswd) != string(pswd2) {
		log.Fatal("password mismatch")
	}

	// https://godoc.org/golang.org/x/crypto/scrypt

	N := 32768
	r := 8
	p := 1

	skey, err := scrypt.Key(pswd, []byte(*salt), N, r, p, 32)

	if err != nil {
		log.Fatal(err)
	}

	var key [32]byte
	copy(key[:], skey)

	sb := Secretbox{
		Key:    key,
		Suffix: *suffix,
		Debug:  *debug,
	}

	for _, abs_path := range files {

		var sb_path string
		var sb_err error

		if *unlock {
			sb_path, sb_err = sb.UnlockFile(abs_path)
		} else {
			sb_path, sb_err = sb.LockFile(abs_path)
		}

		if sb_err != nil {
			log.Fatal(sb_err)
		}

		log.Println(sb_path)
	}

	os.Exit(0)
}
