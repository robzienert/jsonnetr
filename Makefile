BINARY=jsonnetr

LDFLAGS="-X main.version=$(version)"

build: clean
	mkdir build
	env GOOS=darwin GOARCH=amd64 go build -ldflags ${LDFLAGS} -o build/${BINARY}-darwin-amd64 ./*.go

clean:
	rm -rf build
