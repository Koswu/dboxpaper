# dboxpaper

[![Test Status](https://github.com/mattn/dboxpaper/workflows/test/badge.svg?branch=master)](https://github.com/mattn/dboxpaper/actions?workflow=test)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattn/dboxpaper)](https://goreportcard.com/report/github.com/mattn/dboxpaper)

![dboxpaper](https://raw.githubusercontent.com/mattn/dboxpaper/master/dboxpaper-logo256.png)

client for Dropbox Paper

## Usage

```
NAME:
   dboxpaper - Dropbox Paper client

USAGE:
   dboxpaper [global options] command [command options] [arguments...]
   
VERSION:
   0.0.1
   
AUTHOR(S):
   mattn <mattn.jp@gmail.com> 
   
COMMANDS:
     cat         Cat paper
     delete      Delete paper permanently
     list, ls    Show papers
     upload, up  Upload paper
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

list papers

```
$ dboxpaper list
```

cat paper

```
$ dboxpaper cat XXXXXXXXXXX
```

upload paper

```
$ cat README.md | dboxpaper upload
```

update paper

```
$ cat README.md | dboxpaper upload XXXXXXXXXXX
```

delete paper

```
$ dboxpaper delete XXXXXXXXXXX
```

## Installation

```
$ go install github.com/mattn/dboxpaper@latest
```

## Configuration

### Dropbox App Credentials

dboxpaper requires you to set the following environment variables with your Dropbox App credentials:

- `DROPBOX_CLIENT_ID`: Your Dropbox App's Client ID (required)
- `DROPBOX_CLIENT_SECRET`: Your Dropbox App's Client Secret (required)

Example:

```bash
export DROPBOX_CLIENT_ID="your_client_id"
export DROPBOX_CLIENT_SECRET="your_client_secret"
dboxpaper list
```

To create a Dropbox App and obtain these credentials, visit the [Dropbox App Console](https://www.dropbox.com/developers/apps).

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
