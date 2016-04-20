package testmodule

import (
	"camssh"
	"cliargs"
	"fmt"
	"modules"
	"output"
)

type TestModule modules.TCamerataModule

func init() {
	modules.Register("test", &TestModule{},
		"Runs \"whoami\" command on target hosts.")
}

func (me *TestModule) Setup(args *cliargs.Arguments, stdout *output.StdoutManager, stderr *output.StderrManager) {
	me.Args = args
	me.Stdout = stdout
	me.Stderr = stderr
}

func (me *TestModule) Prepare(host string, sshconn *camssh.SshConnection) error {
	me.Host = host
	me.Sshconn = sshconn

	return nil
}

func (me *TestModule) Run() error {

	if !me.Args.Sudo {
		result, err := me.Sshconn.WhoAmI()
		if err != nil {
			panic(err.Error())
		} else {
			me.Stdout.Print(">>> WhoAmI @ ")
			fmt.Println(me.Host, result)
		}
	} else {
		sudo_result, err := me.Sshconn.SudoWhoAmI(me.Args)
		if err != nil {
			panic(err.Error())
		} else {
			me.Stdout.Print(">>> SudoWhoAmI @ ")
			fmt.Println(me.Host, sudo_result)
		}
	}

	return nil

}
