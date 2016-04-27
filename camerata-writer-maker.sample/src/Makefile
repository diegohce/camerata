GOOS=linux
GOFILES_main=camerata-writer.go 

all: build_camerata-writer

run: run_camerata-writer	

build_camerata-writer:
	GOOS=$(GOOS) go build $(GOFILES_main)

run_camerata-writer:
	GOOS=$(GOOS) go run $(GOFILES_main)

.PHONY clean:
	rm -f camerata-writer camerata-writer.exe

