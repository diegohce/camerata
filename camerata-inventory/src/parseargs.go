package main

import (
	"errors"
	"flag"
)

type Arguments struct {
	User        string
	Pass        string
	Host        string
	Insecure    bool
	AskPass     bool
	Bastion     string
	BastionNets string
}

func (me *Arguments) Parse() {

	flag.StringVar(&me.User, "user", "", "Login user")
	flag.StringVar(&me.Pass, "pass", "", "Use this password")

	flag.BoolVar(&me.Insecure, "insecure", true, "No credentials check")

	flag.StringVar(&me.Host, "host", "", "vSphere host or IP")
	flag.StringVar(&me.Bastion, "bastion", "", "Bastion or jumpbox server (name or ip address)")
	flag.StringVar(&me.BastionNets, "bastion-nets", "", "Comma separated list of segments that uses --bastion (e.g.: 10.54.165.,10.54.170.)")

	flag.Parse()
}

func (me *Arguments) Validate() error {

	me.AskPass = true

	if me.User == "" {
		return errors.New("--user cannot be empty")
	}

	if len(me.Pass) > 0 {
		me.AskPass = false
	}

	if len(me.Host) == 0 {
		return errors.New("--host cannot be empty")
	}

	return nil
}