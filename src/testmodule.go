package main

import (
	"fmt"
)

type TestModule TCamerataModule

func (me *TestModule) Run(sshconn *SshConnection) error {

	if !me.args.Sudo {
		result, err := sshconn.WhoAmI()
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println(">>> WhoAmI @", me.host, result)
		}
	} else {
		sudo_result, err := sshconn.SudoWhoAmI(me.args)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println(">>> SudoWhoAmI @", me.host, sudo_result)
		}
	}

	return nil

}
