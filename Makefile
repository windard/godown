
ifeq ("${GOOS}", "")
	GOOS:=darwin
endif

ifeq ("${GOARCH}", "")
	GOARCH:=amd64
endif

ifeq ("${OUTPUT_FILE}", "")
	OUTPUT_FILE:=godown
endif


build: prepare
	go build -v -o ${OUTPUT_FILE} .

prepare:
	go get -v ./...

test:
	go test -v .

package:
	tar -czvf godown.tar.gz ${OUTPUT_FILE} LICENSE README.md
