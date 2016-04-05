package main

import (
	"fmt"
	//"bytes"
	"os"
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
	fmt.Println(">>> CopyModule >>> Filename is", filename)

	fp, err := os.Open(commandargs[0])
	if err != nil {
		panic("Error opening file " + err.Error())
	}
	byte_buffer := make([]byte, 4096)
	fileinfo, err := fp.Stat()
	fmt.Println(">>> CopyModule >>> Filesize is", fileinfo.Size())

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		fmt.Fprintln(w, "C0644", fileinfo.Size(), filename)

		count, _ := fp.Read(byte_buffer)
		for count > 0 {
			fmt.Println(">>> CopyModule >>> sending", count, "bytes")
			fmt.Fprintf(w, "%s", byte_buffer[:count])
			count, _ = fp.Read(byte_buffer)
		}
		fmt.Fprint(w, "\x00")

	}()

	scp_command := fmt.Sprintf("scp -qrt %s/%s", commandargs[1], filename)
	//fmt.Println(">>> scp command is", scp_command)

	//var b bytes.Buffer
	//session.Stdout = &b
	if err := session.Run(scp_command); err != nil {
		panic("Failed to run: " + err.Error())
	}
	//fmt.Println(">>> CopyModule >>>", b.String())

	return nil

}
