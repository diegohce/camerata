GOOS=linux
GOFILES_main=camerata-inventory.go parseargs.go

all: build_camerata-inventory

run: run_camerata-inventory	

build_camerata-inventory:
	GOOS=$(GOOS) go build $(GOFILES_main)

run_camerata-inventory:
	GOOS=$(GOOS) go run $(GOFILES_main)

.PHONY clean:
	rm -f camerata-inventory camerata-inventory.exe

