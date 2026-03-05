version := `git describe --tags --always --dirty`
pkg := "github.com/sinclairtarget/michel"

build: build-static

# Build michel, statically linking libatrus
# Requires built version of the library in libatrus-go
build-static build_tag="dev":
    go build -o michel \
        -ldflags '-X main.Version={{version}} -X main.BuildTag={{build_tag}}' \
        -tags bundled_libatrus

# Build michel, dynamically linking libatrus at a standard path
build-system build_tag="dev":
    go build -o michel \
        -ldflags '-X main.Version={{version}} -X main.BuildTag={{build_tag}}'

# Build michel, dynamically linking libatrus at a non-standard path
# Still requires libatrus pkg-config file in pkg-config search path
# Uses -rpath argument to linker
build-shared build_tag="dev":
    go build -o michel \
        -ldflags "-X main.Version={{version}} -X main.BuildTag={{build_tag}} -r $(pkg-config --variable=libdir libatrus)"

# Run unit tests. Libatrus is statically linked
[group("test")]
test-unit:
    -go test -tags bundled_libatrus ./internal/...

# Run functional tests for CLI. Libatrus is statically linked
[group("test")]
test-cli: 
    mkdir -p test/bin
    go build -o test/bin/michel \
        -ldflags '-X main.Version={{version}} -X main.BuildTag=test' \
        -tags bundled_libatrus
    -go test -ldflags \
        '-X {{pkg}}/test/michel.ExecutablePath={{absolute_path("./test/bin/michel")}}' \
        ./test/cli
    rm -rf test/bin

# Run all tests
[group("test")]
test: test-unit test-cli

fmt:
    goimports -w $(git ls-files *.go **/*.go)

clobber:
    rm -f michel
