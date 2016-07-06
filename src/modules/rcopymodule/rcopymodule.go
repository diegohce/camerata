package rcopymodule

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
	"path"
	"strings"
	//	"time"
)

type Rcopy struct {
	modules.TCamerataModule
	MyArgs   map[string]string
	commands []string
}

func init() {

	rcopymodule_description := `Copy remote file(s) specified in 'source' to 'target' local directory.
		source=/file/to/grab1[|/file/to/grab2...|/file/to/grab#
		target=/local/destination/dir
		
		`

	modules.Register("rcopy", &Rcopy{}, rcopymodule_description)
}

//799241

func (me *Rcopy) Setup(args *cliargs.Arguments, stdout *output.StdoutManager, stderr *output.StderrManager) {
	me.Args = args
	me.Stdout = stdout
	me.Stderr = stderr

	me.MyArgs = modules.ModuleArgsMap(args.MArguments)

	//fmt.Printf("%+v\n", me.MyArgs)

}

func (me *Rcopy) Prepare(host string, sshconn *camssh.SshConnection) error {
	me.Host = host
	me.Sshconn = sshconn

	me.commands = []string{}

	source, ok := me.MyArgs["source"]
	if !ok {
		return errors.New("Missing 'source' argument.")
	}

	if _, ok := me.MyArgs["target"]; !ok {
		return errors.New("Missing 'target' argument.")
	}

	me.commands = strings.Split(source, "|")

	return nil
}

func (me *Rcopy) Run() error {

	me.Stdout.Println(">>> RcopyModule >>> ", me.Args.MArguments)

	for _, command := range me.commands {

		me.Stdout.Print(">>> RcopyModule >>> ", command)

		//		if err := me.runMapOutput(command); err != nil {
		//			return err
		//		}

		fname := path.Base(command)

		fullname := path.Join(me.MyArgs["target"], fname+"_"+me.Host)
		fmt.Println(fullname)

		cmd := "cat " + command
		if err := me.runMapOutputToFile(cmd, fullname); err != nil {
			//if err := me.runMapOutput(cmd); err != nil {
			return err
		}

		me.Stdout.Println("...Done")

	}

	return nil

}

func (me *Rcopy) getFileWriter(filename string) (*os.File, error) {

	fw, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return fw, nil

}

func (me *Rcopy) runMapOutputToFile(command string, filename string) error {

	file_writer, err := me.getFileWriter(filename)
	if err != nil {
		return err
	}
	defer file_writer.Close()

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
		//io.Copy(file_writer, r_stdout)
		//io.Copy(os.Stdout, r_stdout)
		io.Copy(os.Stderr, r_stderr)
	}()

	var gof_error error

	go func() {
		session.Stdin = os.Stdin
		if gof_error = session.Run(commandline); gof_error != nil {
			return
		}
	}()

	//Sync copy
	io.Copy(file_writer, r_stdout)

	return gof_error

}

func (me *Rcopy) runOutput(command string) (string, error) {
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

func (me *Rcopy) runMapOutput(command string) error {
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
