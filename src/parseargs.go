package main

import (
	"camerata/errors"
	"flag"
)

type Arguments struct {
	User    string
	Pass    string
	Bastion string
	AskPass bool
	Sudo    bool
	/* SudoPass    string
	AskSudoPass bool*/
	Hosts      string
	Inventory  string
	Module     string
	MArguments string
}

func (me *Arguments) Parse() {

	flag.StringVar(&me.User, "user", "", "Login user")
	flag.BoolVar(&me.AskPass, "ask-pass", true, "Asks for password on the command line")
	flag.StringVar(&me.Pass, "pass", "", "Asks for password")

	flag.StringVar(&me.Bastion, "bastion", "", "Bastion or jumpbox server")

	flag.BoolVar(&me.Sudo, "sudo", false, "Run as sudo")

	flag.StringVar(&me.Hosts, "hosts", "", "Comma separated hosts list")
	flag.StringVar(&me.Inventory, "inventory", "", "Inventory file")

	flag.StringVar(&me.Module, "module", "command", "Module to run")
	flag.StringVar(&me.MArguments, "args", "", "Module arguments")

	flag.Parse()
}

func (me *Arguments) Validate() error {

	if me.User == "" {
		return errors.CamerataArgumentsError{"User cannot be empty"}
	}
	if len(me.Pass) > 0 && me.AskPass {
		return errors.CamerataArgumentsError{"--pass and --ask-pass cannot be combined"}
	}

	if len(me.Pass) == 0 && !me.AskPass {
		return errors.CamerataArgumentsError{"Must define --pass or --ask-pass"}
	}

	if len(me.Inventory) > 0 && len(me.Hosts) > 0 {
		return errors.CamerataArgumentsError{"--inventory and --hosts cannot be combined"}
	}

	if len(me.Inventory) == 0 && len(me.Hosts) == 0 {
		return errors.CamerataArgumentsError{"Must define --inventory or --hosts"}
	}

	return nil
}
