version := `git describe --tags --always --dirty`
pkg := "github.com/sinclairtarget/michel"

build: (_build "dev")

_build build_tag:
    go build -o michel \
        -ldflags '-X main.Version={{version}} -X main.BuildTag={{build_tag}}'

# Run unit tests
[group("test")]
test-unit:

# Run functional tests for CLI
[group("test")]
test-cli: 
    mkdir -p test/bin
    go build -o test/bin/michel \
        -ldflags '-X main.Version={{version}} -X main.BuildTag=test'
    -go test -ldflags \
        '-X {{pkg}}/test/michel.ExecutablePath={{absolute_path("./test/bin/michel")}}' \
        ./test/cli
    rm -rf test/bin

# Run all tests
[group("test")]
test: test-cli

fmt:
    goimports -w $(git ls-files *.go **/*.go)

clobber:
    rm -f michel
