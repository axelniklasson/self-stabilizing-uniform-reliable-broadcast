# Self-stabilizing Uniform Reliable Broadcast (URB)

## Set up
First, make sure that [Go](https://golang.org/doc/install) is installed. This project was developed for version `1.13`, so that is the recommended version to use. Then, install [dep](https://golang.github.io/dep/docs/installation.html).

When that is done, run the following commands.
```
cd $GOPATH/src/github.com && git clone https://github.com/axelniklasson/self-stabilizing-uniform-reliable-broadcast
cd self-stabilizing-uniform-reliable-broadcast
echo "go vet ./... && go list ./... | xargs -n 1 golint -set_exit_status && go test -v ./..." >> .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit scripts/*.sh
dep ensure
```

## Running locally
```
./scripts/start.sh NUMBER_OF_NODES
```

## Testing
All unit tests can be run through the bash script as `sh scripts/test.sh`.
