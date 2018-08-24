git describe --always --long --dirty
go build -i -v -ldflags "-X 'main.version=2.0-0-g08e574e-dirty'"
