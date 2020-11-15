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

	"github.com/AlecAivazis/survey/v2"
	"github.com/skanehira/clipboard-image/v2"
)

var version = "0.0.1"

var token string

var (
	useClipboard         bool
	generateMarkdownLink bool
	isInteractive        bool
	title                string
	desc                 string
	fileName             string
)

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

	if isInteractive {
		var (
			qs      = []*survey.Question{}
			answers = struct {
				Filename string
				Title    string
				Desc     string
			}{}
		)

		prompt := &survey.Select{
			Message: "Which upload",
			Options: []string{"specify file", "clipboard"},
			Default: "red",
		}

		var selected string
		if err := survey.AskOne(prompt, &selected); err != nil {
			return err
		}

		if selected == "clipboard" {
			r, err = clipboard.Read()
			if err != nil {
				return err
			}
			qs = append(qs, &survey.Question{
				Name: "fileName",
				Prompt: &survey.Input{
					Message: "file name:",
				},
			})
		} else {
			qs = append(qs, &survey.Question{
				Name: "fileName",
				Prompt: &survey.Input{
					Message: "Select file:",
					Suggest: func(path string) []string {
						files, _ := filepath.Glob(path + "*")
						return files
					},
				},
			})
		}

		prompts := []*survey.Question{
			{
				Name: "title",
				Prompt: &survey.Input{
					Message: "Title:",
				},
			},
			{
				Name: "desc",
				Prompt: &survey.Multiline{
					Message: "Description:",
				},
			},
		}
		qs = append(qs, prompts...)
		if err := survey.Ask(qs, &answers); err != nil {
			return err
		}

		if r == nil && answers.Filename != "" {
			r, err = os.Open(answers.Filename)
			if err != nil {
				return err
			}
		}

		fileName = answers.Filename
		title = answers.Title
		desc = answers.Desc
	} else {
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
	}

	meta := Meta{
		Title:    title,
		Desc:     desc,
		Filename: fileName,
	}

	image, err := upload(meta, r)
	if err != nil {
		return err
	}

	if generateMarkdownLink {
		fmt.Println(fmt.Sprintf("![](%s)", image.URL))
		return nil
	}

	fmt.Println(image.URL)
	return nil
}

func upload(meta Meta, r io.Reader) (*Image, error) {
	gyazo, err := NewClient(token)
	if err != nil {
		return nil, err
	}
	image, err := gyazo.Upload(meta, r)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func main() {
	name := "gyazo"
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.BoolVar(&useClipboard, "c", false, "")
	fs.BoolVar(&generateMarkdownLink, "m", false, "")
	fs.BoolVar(&isInteractive, "i", false, "")
	fs.StringVar(&title, "t", "", "")
	fs.StringVar(&desc, "d", "", "")
	fs.Usage = func() {
		fs.SetOutput(os.Stdout)
		fmt.Printf(`%[1]s - Gyazo CLI

VERSION: %s

USAGE:
  $ %[1]s [-cmtdi] [<] [file]

DESCRIPTION:
  -c	upload image from clipboard
  -m	generate markdown link
  -t	title
  -d	description
  -i	interactive mode

EXAMPLE:
  $ %[1]s < image.png
  $ %[1]s -c -t "gorilla image"
  $ %[1]s -m image.png
  $ %[1]s -i
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
