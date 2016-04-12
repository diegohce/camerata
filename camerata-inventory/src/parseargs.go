package main

import (
	"errors"
	"flag"
)

type Arguments struct {
	User     string
	Pass     string
	Host     string
	Insecure bool
	AskPass  bool
}

func (me *Arguments) Parse() {

	flag.StringVar(&me.User, "user", "", "Login user")
	flag.StringVar(&me.Pass, "pass", "", "Use this password")

	flag.BoolVar(&me.Insecure, "insecure", false, "No credentials check")

	flag.StringVar(&me.Host, "host", "", "vSphere host or IP")

	flag.Parse()
}

func (me *Arguments) Validate() error {

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
