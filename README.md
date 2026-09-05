# cmgrep

cmgrep (Charlie Max GREP) is a tool for distributed grepping.

## Hosts

fa26-cs425-72NN.cs.illinois.edu

## Installation

Ensure you have go installed on your machine per <https://go.dev/doc/install>

```
git config core.hooksPath .githooks
```

## Usage

```
go build -o cmgrep .

# generate this machine's log (node number read from the hostname)
./scripts/gen_logs.py

# run the server (node number read from the hostname)
./cmgrep serve

# pass the node explicitly if the hostname is not a VM name
./scripts/gen_logs.py 3
./cmgrep serve --node=3

# run the client (from anywhere that can reach the VMs)
./cmgrep "/morerare"
./cmgrep -i "/MORERARE"
./cmgrep -E "/(morerare|alsoquiterare)"
```

## Testing

```
# off the VMs, no argument generates all ten logs
# (fixed seed, so a VM generating its own log produces the same bytes)
./scripts/gen_logs.py
```

## Building

```
# build (the whole package, not just main.go)
go build -o cmgrep .
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


