package main

/*
[bastion]
host="trustedhost.olleros"
user="dcena"
password=""

[servers]

	[servers.avinet]
	host="avinet.olleros"
	user=""
	password=""
	sudo=false
	sudo_nopass=false
	use_bastion=true

	[servers.callbacksd-pci]
	host="callbacksd-pci.olleros"
	user=""
	password=""
	sudo=false
	sudo_nopass=false
	use_bastion=true

[[modules]]
name="command"
args='''for i in {1,2,3}
do
	echo "HELLO WORLD $i"
done'''

[[modules]]
name="command"
args='''for i in {1,2,3}
do
	echo "HELLO WORLD $i"
done'''

*/

import (
	//"fmt"

	"github.com/BurntSushi/toml"
)

type Inventory struct {
	Bastion BastionServer `toml:"bastion"`
	Servers map[string]Server
	Modules []ServerModule `toml:"modules"`
}

type BastionServer struct {
	Host     string
	User     string
	Password string
}

type Server struct {
	Host       string
	User       string
	Password   string
	Sudo       bool
	SudoNoPass bool `toml:"sudo_nopass"`
	UseBastion bool `toml:"use_bastion"`
}

type ServerModule struct {
	Name string
	Args string
}

func ParseInventory(inventoryfile string) (*Inventory, error) {
	inventory := &Inventory{}

	_, err := toml.DecodeFile(inventoryfile, inventory)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}
