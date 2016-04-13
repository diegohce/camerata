package main

import (
	"fmt"
)

type SomeModule2 TCamerataModule

func init() {
	RegisterModule2("some", &SomeModule2{})
	fmt.Println("somemodule2::init")
}

func (me *SomeModule2) Prepare(host string, sshconn *SshConnection) error {
	me.host = host
	me.sshconn = sshconn

	return nil
}

func (me *SomeModule2) Run() error {
	fmt.Println("Hello from SomeMod2")
	return nil
}
