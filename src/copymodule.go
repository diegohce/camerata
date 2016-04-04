package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type CopyModule TCamerataModule

func (me *CopyModule) Prepare(host string, sshconn *SshConnection) error {
	me.host = host
	me.sshconn = sshconn

	if len(me.args.MArguments) == 0 {
		return CamerataModuleError{"CopyModule: Arguments cannot be empty"}
	}

	commandargs := strings.Split(me.args.MArguments, "|")
	if len(commandargs) != 2 {
		return CamerataModuleError{"CopyModule: Arguments must be src_file|dest_dir"}
	}

	return nil
}

func (me *CopyModule) Run() error {

	commandargs := strings.Split(me.args.MArguments, "|")

	fmt.Print(">>> CopyModule >>> Copying ", me.args.MArguments, "@", me.host)
	if me.args.Sudo {
		fmt.Print(" as sudo")
	}
	fmt.Println("")

	session, err := me.sshconn.client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	filename := filepath.Base(commandargs[0])

	fmt.Println(">>> Filename is", filename)

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		content, err := ioutil.ReadFile(commandargs[0])
		if err != nil {
			panic("Failed to read file " + commandargs[0])
		}

		fmt.Fprintln(w, "C0644", len(content), filename)
		fmt.Fprint(w, content)
		fmt.Fprint(w, "\x00")
	}()

	scp_command := fmt.Sprintf("scp -qrt %s/%s", commandargs[1], filename)

	if err := session.Run(scp_command); err != nil {
		panic("Failed to run: " + err.Error())
	}

	return nil

}
