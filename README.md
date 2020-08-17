# Go CMS
CMS in Go

## Testing
```bash
go test
```
## Building
```bash
GOOS=linux GOARCH=amd64 go build
docker build -t go-cms .
```
## Running
Run by providing the database connection and path to the documentation
```bash
go-cms \
    -conn="<username>:<password>@tcp(<hostname>:3306)/<schema>" \
    -port=7335 \
    -path=<path to documentation>
```
or with docker
```bash
docker run --rm -p 7335:7335 go-cms
```
## Usage
Request the page by supplying the page id in the uri
```bash
curl localhost:7335/index
```
