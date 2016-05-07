package commandmodule

import (
	//"bytes"
	"camssh"
	"cliargs"
	"errors"
	"fmt"
	"io"
	"modules"
	"os"
	"output"
)

type CommandModule modules.TCamerataModule

func init() {
	modules.Register("command", &CommandModule{},
		"Executes --args command line on target hosts.")
}

func (me *CommandModule) Setup(args *cliargs.Arguments, stdout *output.StdoutManager, stderr *output.StderrManager) {
	me.Args = args
	me.Stdout = stdout
	me.Stderr = stderr
}

func (me *CommandModule) Prepare(host string, sshconn *camssh.SshConnection) error {
	me.Host = host
	me.Sshconn = sshconn

	if len(me.Args.MArguments) == 0 {
		return errors.New("CommandModule: Arguments cannot be empty")
	}

	return nil
}

func (me *CommandModule) Run() error {

	commandargs := me.Args.MArguments

	me.Stdout.Print(">>> CommandModule >>> Executing ", commandargs, " @", me.Host)
	if me.Args.Sudo {
		me.Stdout.Print(" as sudo")
	}
	me.Stdout.Println("")

	commandline := commandargs

	session, err := me.Sshconn.Client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

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

	if me.Args.Sudo && me.Args.User != "root" {
		if !me.Args.SudoNoPass {
			commandline = fmt.Sprintf("echo %s | sudo -S %s", me.Args.Pass, commandargs)
			//			commandline = fmt.Sprintf("sudo -S \"%s\"", commandargs)
			//			w, _ := session.StdinPipe()
			//			defer w.Close()
			//			go fmt.Fprintln(w, me.Args.Pass)

		} else {
			commandline = fmt.Sprintf("sudo %s", commandargs)
		}

		//		go func() {
		//			w, err := session.StdinPipe()
		//			if err != nil {
		//				panic("Error on stdinpipe: " + err.Error())
		//			}
		//			defer w.Close()

		//			if !me.Args.SudoNoPass {
		//				fmt.Fprintln(w, me.Args.Pass)
		//			}
		//			session.Stdin = os.Stdin
		//		}()

	}
	session.Stdin = os.Stdin
	if err := session.Run(commandline); err != nil {
		me.Stderr.Println("Failed to run: ", err)
	}

	return nil

}
