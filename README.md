# cmgrep

cmgrep (Charlie Max GREP) is a tool for distributed grepping.

## Installation

Ensure you have go installed on your machine per <https://go.dev/doc/install>

## Usage

```
go run main.go
```

## Building

```
# build
go build -o cmgrep main.go
# run
./cmgrep
```

## Generating protobufs

```
brew install protobuf
protoc --version # check

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
which protoc-gen-go # check

protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/*.proto

```
