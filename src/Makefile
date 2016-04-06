GOOS=linux
GOFILES_main=camerata.go parseargs.go sshconnection.go errors.go modulescommon.go testmodule.go commandmodule.go copymodule.go inventory.go


all: build_camerata

run: run_camerata

build_camerata:
	GOOS=$(GOOS) go build $(GOFILES_main)

run_camerata:
	GOOS=$(GOOS) go run $(GOFILES_main)

.PHONY clean:
	rm -f camerata camerata.exe

