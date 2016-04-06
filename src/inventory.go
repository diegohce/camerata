package main

/*
[bastion]
host="IP:port or fqdn:port"
user="user to use on login"
password="*******"

[servers]

	[servers.heymamma]
	host="heymamma.tienda"
	user="diego"
	password=""
	sudo=false
	sudo_nopass=false
	use_bastion=false


	[servers.pbot01]
	host="pbot01.laspornografas.com"
	user="root"
	password=""
	sudo=false
	sudo_nopass=false
	use_bastion=false
*/

import (
	//"fmt"

	"github.com/BurntSushi/toml"
)

type Inventory struct {
	Bastion BastionServer `toml:"bastion"`
	Servers map[string]Server
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

func ParseInventory(inventoryfile string) (*Inventory, error) {
	inventory := &Inventory{}

	_, err := toml.DecodeFile(inventoryfile, inventory)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}
