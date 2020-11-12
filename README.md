# gyazo
This is CLI for [Gyazo](https://gyazo.com)

![](https://i.gyazo.com/13bdd42f09cc63e17703c09f7c19514e.gif)

## Supported OS
- Mac
- Windows
- Linux

## Features
- upload image file from clipboard
- generate markdown link

## Requirements
- xclip (only linux)
- file (only linux)

## Installtion
1. Install this CLI
   ```bash
   $ git clone https://github.com/skanehira/gyazo
   $ cd gyazo && go install
   ```

2. Register application and get your token from [here](https://gyazo.com/api) and set it
   ```bash
   # set environment
   export GYAZO_TOKEN = XXXXXXXXXXXXXXX

   # or set $HOME/.gyazo_token
   echo "XXXXXXXXXX" > $HOME/.gyazo_token
   ```

## Usage
```bash
$ gyazo -h
gyazo - Gyazo CLI

VERSION: 0.0.1

USAGE:
  $ gyazo [-c] [-m] [<] [file]

DESCRIPTION:
  -c    upload image from clipboard
  -m    generate markdown link

EXAMPLE:
  $ gyazo < image.png
  $ cat image.png | gyazo
  $ gyazo -m image.png
  $ gyazo -c -m
```

## Author
skanehira

