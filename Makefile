GOOS=linux

all: camerata camerata_inventory camerata_writer

install: camerata camerata_inventory camerata_writer
	mkdir -p bin
	mv src/camerata \
		src/camerata.exe \
		camerata-inventory/src/camerata-inventory \
		camerata-inventory/src/camerata-inventory.exe
		camerata-writer/src/camerata-writer \
		camerata-writer/src/camerata-writer.exe bin 2>/dev/null; true


camerata: 
	. ./goenv.sh; cd src ; make 

camerata_inventory: 
	cd camerata-inventory ; . ../goenv.sh; cd src; make

camerata_writer: 
	cd camerata-writer ; . ../goenv.sh; cd src; make

.PHONY clean:
	rm -f src/camerata \
		src/camerata.exe \
		camerata-inventory/src/camerata-inventory \
		camerata-inventory/src/camerata-inventory.exe \
		camerata-writer-src/camerata-writer \
		camerata-writer-src/camerata-writer.exe

