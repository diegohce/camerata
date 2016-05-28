package pipmodule

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

type Pip struct {
	modules.TCamerataModule
	MyArgs   map[string]string
	commands []string
}

func init() {

	pipmodule_description := `pyhthon modules dependencies.
	
		*name=pip_package_name
		*requirements=/path/to/requirements.txt
		
		index=http://url/to/pip/index
		
		virtualenv=/path/to/virtualenv_to_activate
		
		virtualenv_create=no || yes (default no)
			If set to yes, create virtualenv and activate before installing
			
		virtualenv_site_packages=no || yes (default no)
			If virtualenv_create == yes, use --site-packages to create virtualenv
			
		* = One of 'name' or 'requirements' is Required`

	modules.Register("pip", &Pip{}, pipmodule_description)
}

//799241

func (me *Pip) Setup(args *cliargs.Arguments, stdout *output.StdoutManager, stderr *output.StderrManager) {
	me.Args = args
	me.Stdout = stdout
	me.Stderr = stderr

	me.MyArgs = modules.ModuleArgsMap(args.MArguments)

	//fmt.Printf("%+v\n", me.MyArgs)

}

func (me *Pip) Prepare(host string, sshconn *camssh.SshConnection) error {
	me.Host = host
	me.Sshconn = sshconn

	me.commands = []string{}

	_, reqs := me.MyArgs["requirements"]
	_, pack := me.MyArgs["name"]

	if !reqs && !pack {
		return errors.New("Missing 'name' and 'requirements' argument. Must specify one.")
	}

	if value, ok := me.MyArgs["virtualenv_create"]; ok && value == "yes" {

		if _, o := me.MyArgs["virtualenv"]; !o {
			return errors.New("Specified 'virtualenv_create' but missing 'virtualenv' path.")
		}
	}

	if value, ok := me.MyArgs["virtualenv_create"]; ok && value == "yes" {

		if value, ok := me.MyArgs["virtualenv_site_packages"]; ok && value == "yes" {

			me.commands = append(me.commands, fmt.Sprintf("virtualenv --site-packages %s", me.MyArgs["virtualenv"]))

		} else {
			me.commands = append(me.commands, fmt.Sprintf("virtualenv  %s", me.MyArgs["virtualenv"]))

		}
	}

	pipcommand := "pip install"

	if value, ok := me.MyArgs["index"]; ok {
		pipcommand = fmt.Sprintf("%s -i %s", pipcommand, value)
	}

	if value, ok := me.MyArgs["virtualenv"]; ok {
		pipcommand = fmt.Sprintf(". %s/bin/activate && %s", value, pipcommand)
	}

	if reqs {
		pipcommand = fmt.Sprintf("%s -r %s", pipcommand, me.MyArgs["requirements"])
	} else {
		pipcommand = fmt.Sprintf("%s %s", pipcommand, me.MyArgs["name"])
	}

	me.commands = append(me.commands, pipcommand)

	return nil
}

func (me *Pip) Run() error {

	me.Stdout.Println(">>> PipModule >>> ", me.Args.MArguments)

	for _, command := range me.commands {

		me.Stdout.Print(">>> PipModule >>> ", command)

		if err := me.runMapOutput(command); err != nil {
			return err
		}

		me.Stdout.Println("...Done")

	}

	return nil

}

func (me *Pip) runOutput(command string) (string, error) {
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

func (me *Pip) runMapOutput(command string) error {
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
