package main

import (
	"bytes"
	"fmt"
)

type CommandModule TCamerataModule

func (me *CommandModule) Run(sshconn *SshConnection) error {

	commandargs := me.args.MArguments

	if len(commandargs) == 0 {
		return CamerataModuleError{"CommandModule: Arguments cannot be empty"}
	}

	fmt.Print(">>> CommandModule >>> Executing", commandargs, "@", me.host)
	if me.args.Sudo {
		fmt.Print(" as sudo")
	}
	fmt.Println("")

	commandline := commandargs

	session, err := sshconn.client.NewSession()
	if err != nil {
		panic("Failed to create sudoSession: " + err.Error())
	}
	defer session.Close()

	if me.args.Sudo {
		commandline = fmt.Sprintf("sudo -S \"%s\"", commandargs)

		go func() {
			w, _ := session.StdinPipe()
			defer w.Close()
			fmt.Fprintln(w, me.args.Pass)
		}()
	}
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(commandline); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(">>> CommandModule >>>", b.String())

	return nil

}
