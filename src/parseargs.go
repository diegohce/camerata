package main

import (
	"flag"
	"strings"
)

type Arguments struct {
	User           string
	Pass           string
	Bastion        string
	BastionUser    string
	BastionPass    string
	AskPass        bool
	AskBastionPass bool
	Sudo           bool
	/* SudoPass    string
	AskSudoPass bool*/
	Hosts      string
	Inventory  string
	Module     string
	MArguments string
	Test       bool
}

func (me *Arguments) Parse() {

	flag.StringVar(&me.User, "user", "", "Login user")
	flag.BoolVar(&me.AskPass, "ask-pass", true, "Asks for password on the command line")
	flag.StringVar(&me.Pass, "pass", "", "Use this password")

	flag.StringVar(&me.Bastion, "bastion", "", "Bastion or jumpbox server")
	flag.StringVar(&me.BastionUser, "bastion-user", "", "Bastion or jumpbox server login user (default: same as --user)")
	flag.StringVar(&me.BastionPass, "bastion-pass", "", "Bastion or jumpbox server password (default: same as --pass)")
	flag.BoolVar(&me.AskBastionPass, "ask-bastion-pass", false, "Asks for password on the command line for bastion jump")

	flag.BoolVar(&me.Sudo, "sudo", false, "Run as sudo")

	flag.StringVar(&me.Hosts, "hosts", "", "Comma separated hosts list")
	flag.StringVar(&me.Inventory, "inventory", "", "Inventory file")

	flag.StringVar(&me.Module, "module", "test", "Module to run")
	flag.StringVar(&me.MArguments, "args", "", "Module arguments")

	flag.BoolVar(&me.Test, "test", false, "Runs whoami on remote host")

	flag.Parse()
}

func (me *Arguments) Validate() error {

	if me.User == "" {
		return CamerataArgumentsError{"User cannot be empty"}
	}
	if len(me.Pass) > 0 {
		me.AskPass = false
	}

	if len(me.Pass) == 0 && !me.AskPass {
		return CamerataArgumentsError{"Must define --pass or --ask-pass"}
	}

	if len(me.Inventory) > 0 && len(me.Hosts) > 0 {
		return CamerataArgumentsError{"--inventory and --hosts cannot be combined"}
	}

	if len(me.Inventory) == 0 && len(me.Hosts) == 0 {
		return CamerataArgumentsError{"Must define --inventory or --hosts"}
	}

	if len(me.Bastion) > 0 {
		if strings.Index(me.Bastion, ":") < 0 {
			me.Bastion = me.Bastion + ":22"
		}
		if len(me.BastionUser) == 0 {
			me.BastionUser = me.User
		}
		if len(me.BastionPass) == 0 {
			me.BastionPass = me.Pass
		}
	}

	return nil
}
