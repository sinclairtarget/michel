executable_name := "michel"

build:
    go build -o {{executable_name}}

clobber:
    rm -f {{executable_name}}
