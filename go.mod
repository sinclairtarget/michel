module github.com/sinclairtarget/michel

go 1.25.4

require (
	github.com/fsnotify/fsnotify v1.9.0
	github.com/sinclairtarget/libatrus-go v0.0.0-20250929114858-c6b44bf459de
	gopkg.in/yaml.v3 v3.0.1
)

require golang.org/x/sys v0.13.0 // indirect

replace github.com/sinclairtarget/libatrus-go => ../libatrus-go
