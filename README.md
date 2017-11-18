[![Build Status](https://travis-ci.org/thomaslorentsen/go-cms.svg?branch=master)](https://travis-ci.org/thomaslorentsen/go-cms)

# Go CMS
CMS in Go

## Running
Run by providing the database connection and path to the documentation
```bash
go-cms \
    -conn="<username>:<password>@tcp(<hostname>:3306)/<schema>" \
    -port=7335 \
    -path=<path to documentation>
```

## Usage
Request the page by supplying the page id in the uri
```bash
curl localhost:7335/index
```
