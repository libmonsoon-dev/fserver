all: build-windows build

build-windows:
	GOOS=windows go build -v -o fserve.exe .

build:
	go build -v .
