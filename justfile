executable_name := "michel"
version := `git describe --tags --always --dirty`

build:
    go build -o {{executable_name}} \
        -ldflags '-X main.Version={{version}} -X main.BuildTag=dev'

clobber:
    rm -f {{executable_name}}
