GOOS=linux
GOFILES_main=camerata.go inventory.go version.go

all: build_camerata

run: run_camerata

build_camerata:
	GOOS=$(GOOS) go build $(GOFILES_main)

run_camerata:
	GOOS=$(GOOS) go run $(GOFILES_main)

.PHONY clean:
	rm -f camerata camerata.exe

