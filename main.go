package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/skanehira/clipboard-image/v2"
	"github.com/tomohiro/go-gyazo/gyazo"
)

var version = "0.0.1"

var token string

func getToken() (string, error) {
	token := os.Getenv("GYAZO_TOKEN")
	if token != "" {
		return token, nil
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFile := filepath.Join(dir, ".gyazo_token")
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return "", err
	}

	token = strings.Trim(string(b), "\r\n")

	if token == "" {
		return "", errors.New("gyazo token is empty")
	}

	return token, nil
}

func run(args []string) error {
	var (
		r   io.Reader
		err error
	)

	token, err = getToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(args) == 0 {
		if useClipboard {
			r, err = clipboard.Read()
			if err != nil {
				return err
			}
		} else {
			r = os.Stdin
			if err != nil {
				return err
			}
		}
	} else {
		r, err = os.Open(args[0])
		if err != nil {
			return err
		}
	}

	image, err := upload(r)
	if err != nil {
		return err
	}

	fmt.Println(image.URL)
	return nil
}

func upload(r io.Reader) (*gyazo.Image, error) {
	gyazo, err := gyazo.NewClient(token)
	if err != nil {
		return nil, err
	}
	image, err := gyazo.Upload(r)
	if err != nil {
		return nil, err
	}
	return image, nil
}

var (
	useClipboard bool
)

func main() {
	name := "gyazo"
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.BoolVar(&useClipboard, "c", false, "uplaod from clipboard")
	fs.Usage = func() {
		fs.SetOutput(os.Stdout)
		fmt.Printf(`%[1]s - Gyazo CLI

VERSION: %s

USAGE:
  $ %[1]s [-c] [<] [file]

EXAMPLE:
  $ %[1]s < image.png
  $ cat image.png | %[1]s
  $ %[1]s image.png
  $ %[1]s -c
`, name, version)
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			return
		}
		os.Exit(1)
	}

	if err := run(fs.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
