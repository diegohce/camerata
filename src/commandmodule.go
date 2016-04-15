package main

import (
	//"bytes"
	"fmt"
	"io"
	"os"
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

	//var b bytes.Buffer
	//session.Stdout = &b
	go func() {
		var br int64
		r, _ := session.StdoutPipe()
		br, _ = io.Copy(os.Stdout, r)
		for br > 0 {
			br, _ = io.Copy(os.Stdout, r)
		}
	}()
	go func() {
		var br int64
		r, _ := session.StderrPipe()
		br, _ = io.Copy(os.Stderr, r)
		for br > 0 {
			br, _ = io.Copy(os.Stderr, r)
		}
	}()

	if err := session.Run(commandline); err != nil {
		me.stderr.Println("Failed to run: ", err)
	}
	//me.stdout.Print(">>> CommandModule >>>", b.String())
	//fmt.Println(b.String())

	return nil

}
