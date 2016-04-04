package main

import (
	//	"bytes"
	"fmt"
)

type CopyModule TCamerataModule

func (me *CopyModule) Prepare(host string, sshconn *SshConnection) error {
	me.host = host
	me.sshconn = sshconn

	if len(me.args.MArguments) == 0 {
		return CamerataModuleError{"CopyModule: Arguments cannot be empty"}
	}

	return nil
}

func (me *CopyModule) Run() error {

	commandargs := me.args.MArguments

	fmt.Print(">>> CopyModule >>> Copying ", commandargs, "@", me.host)
	if me.args.Sudo {
		fmt.Print(" as sudo")
	}
	fmt.Println("")

	//	commandline := commandargs

	//	session, err := me.sshconn.client.NewSession()
	//	if err != nil {
	//		panic("Failed to create session: " + err.Error())
	//	}
	//	defer session.Close()

	//	if me.args.Sudo {
	//		commandline = fmt.Sprintf("sudo -S \"%s\"", commandargs)

	//		go func() {
	//			w, _ := session.StdinPipe()
	//			defer w.Close()
	//			fmt.Fprintln(w, me.args.Pass)
	//		}()
	//	}
	//	var b bytes.Buffer
	//	session.Stdout = &b
	//	if err := session.Run(commandline); err != nil {
	//		panic("Failed to run: " + err.Error())
	//	}
	//	fmt.Println(">>> CopyModule >>>", b.String())

	return nil

}
