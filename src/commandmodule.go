package main

import (
	"bytes"
	"fmt"
)

type CommandModule TCamerataModule

func (me *CommandModule) Prepare(host string, sshconn *SshConnection) error {
	me.host = host
	me.sshconn = sshconn

	if len(me.args.MArguments) == 0 {
		return CamerataModuleError{"CommandModule: Arguments cannot be empty"}
	}

	return nil
}

func (me *CommandModule) Run() error {

	commandargs := me.args.MArguments

	me.stdout.Print(">>> CommandModule >>> Executing ", commandargs, " @", me.host)
	if me.args.Sudo {
		me.stdout.Print(" as sudo")
	}
	me.stdout.Println("")

	commandline := commandargs

	session, err := me.sshconn.client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	if me.args.Sudo {
		commandline = fmt.Sprintf("sudo -S \"%s\"", commandargs)

		go func() {
			w, _ := session.StdinPipe()
			defer w.Close()
			if me.args.Sudo && !me.args.SudoNoPass {
				fmt.Fprintln(w, me.args.Pass)
			}
		}()
	}

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(commandline); err != nil {
		panic("Failed to run: " + err.Error())
	}
	me.stdout.Print(">>> CommandModule >>>", b.String())
	fmt.Println(b.String())

	return nil

}
