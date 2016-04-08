package main

import (
	"fmt"
)

type TestModule TCamerataModule

func (me *TestModule) Prepare(host string, sshconn *SshConnection) error {
	me.host = host
	me.sshconn = sshconn

	return nil
}

func (me *TestModule) Run() error {

	if !me.args.Sudo {
		result, err := me.sshconn.WhoAmI()
		if err != nil {
			panic(err.Error())
		} else {
			me.stdout.Print(">>> WhoAmI @ ")
			fmt.Println(me.host, result)
		}
	} else {
		sudo_result, err := me.sshconn.SudoWhoAmI(me.args)
		if err != nil {
			panic(err.Error())
		} else {
			me.stdout.Print(">>> SudoWhoAmI @ ")
			fmt.Println(me.host, sudo_result)
		}
	}

	return nil

}
