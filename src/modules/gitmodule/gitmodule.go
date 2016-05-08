package gitmodule

import (
	"bytes"
	"camssh"
	"cliargs"
	"errors"
	"fmt"
	"io"
	"modules"
	"os"
	"output"
	//	"strconv"
	//	"strings"
	//	"time"
)

type Git struct {
	modules.TCamerataModule
	MyArgs   map[string]string
	commands []string
}

func init() {

	gitmodule_description := `git deploy module.
		**Clone / update**
		repo=git://github.com/diegohce/camerata.git
		dest=/usr/src/camerata
		version=0.1.2
		
		**Rollback version**
		dest=/usr/src/camerata
		rollback_to=0.1.1
		
		**Get version**
		dest=/usr/src/camerata
		version=?`

	modules.Register("git", &Git{}, gitmodule_description)
}

func (me *Git) Setup(args *cliargs.Arguments, stdout *output.StdoutManager, stderr *output.StderrManager) {
	me.Args = args
	me.Stdout = stdout
	me.Stderr = stderr

	me.MyArgs = modules.ModuleArgsMap(args.MArguments)

	//fmt.Printf("%+v\n", me.MyArgs)

}

func (me *Git) Prepare(host string, sshconn *camssh.SshConnection) error {
	me.Host = host
	me.Sshconn = sshconn

	if _, ok := me.MyArgs["dest"]; !ok {
		return errors.New("Missing 'dest' argument")
	}

	if value, ok := me.MyArgs["version"]; ok {
		if value == "?" {

			me.commands = append(me.commands, fmt.Sprintf("cd %s && git describe", me.MyArgs["dest"]))
		}
	}

	return nil
}

func (me *Git) Run() error {

	me.Stdout.Println(">>> GitModule >>> ", me.Args.MArguments)

	for _, command := range me.commands {

		fmt.Println(command)
		output, err := me.runOutput(command)
		if err != nil {
			return err
		}
		fmt.Print(output)

	}

	return nil

}

func (me *Git) runOutput(command string) (string, error) {
	session, err := me.Sshconn.Client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var commandline string

	if me.Args.Sudo && me.Args.User != "root" {
		commandline = fmt.Sprintf("sudo -S -- bash -s <<CMD\n%s\n%s\nCMD", me.Args.Pass, command)
	} else {
		commandline = command
	}

	var b bytes.Buffer
	session.Stdout = &b
	var c bytes.Buffer
	session.Stderr = &c

	if err := session.Run(commandline); err != nil {
		fmt.Println(c.String())
		return c.String(), err
	}

	return b.String(), nil

}

func (me *Git) runMapOutput(command string) error {
	session, err := me.Sshconn.Client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var commandline string

	if me.Args.Sudo && me.Args.User != "root" {
		commandline = fmt.Sprintf("sudo -S -- bash -s <<CMD\n%s\n%s\nCMD", me.Args.Pass, command)
	} else {
		commandline = command
	}

	r_stdin, _ := session.StdinPipe()
	defer r_stdin.Close()

	r_stdout, _ := session.StdoutPipe()
	r_stderr, _ := session.StderrPipe()

	go func() {
		io.Copy(os.Stdout, r_stdout)
		io.Copy(os.Stderr, r_stderr)
	}()

	session.Stdin = os.Stdin
	if err := session.Run(commandline); err != nil {
		return err
	}

	return nil

}
