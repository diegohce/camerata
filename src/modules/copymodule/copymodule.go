package copymodule

import (
	"fmt"
	//"bytes"
	"camssh"
	"cliargs"
	"errors"
	"modules"
	"os"
	"output"
	"path/filepath"
	//"strings"
)

type CopyModule struct {
	modules.TCamerataModule
	CopyArgs map[string]string
}

func init() {
	modules.Register("copy", &CopyModule{},
		"Sends --args source=file to target=directory to target hosts.")
}

func (me *CopyModule) Setup(args *cliargs.Arguments, stdout *output.StdoutManager, stderr *output.StderrManager) {
	me.Args = args
	me.Stdout = stdout
	me.Stderr = stderr

	me.CopyArgs = modules.ModuleArgsMap(args.MArguments)

}

func (me *CopyModule) Prepare(host string, sshconn *camssh.SshConnection) error {
	me.Host = host
	me.Sshconn = sshconn

	if len(me.Args.MArguments) == 0 {
		return errors.New("CopyModule: Arguments cannot be empty")
	}

	if len(me.CopyArgs) != 2 {
		return errors.New("CopyModule: Arguments must be \"source=file target=directory\"")
	}

	return nil
}

func (me *CopyModule) Run() error {

	me.Stdout.Print(">>> CopyModule >>> Copying ", me.Args.MArguments, "@", me.Host)
	if me.Args.Sudo {
		me.Stdout.Print(" as sudo")
	}
	me.Stdout.Println("")

	session, err := me.Sshconn.Client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	filename := filepath.Base(me.CopyArgs["source"])
	me.Stdout.Println(">>> CopyModule >>> Filename is", filename)

	fp, err := os.Open(me.CopyArgs["source"])
	if err != nil {
		panic("Error opening file " + err.Error())
	}
	byte_buffer := make([]byte, 4096)
	fileinfo, err := fp.Stat()
	me.Stdout.Println(">>> CopyModule >>> Filesize is", fileinfo.Size())

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		if me.Args.Sudo && !me.Args.SudoNoPass {
			fmt.Fprintln(w, me.Args.Pass)
		}

		fmt.Fprintln(w, "C0644", fileinfo.Size(), filename)

		count, _ := fp.Read(byte_buffer)
		for count > 0 {
			me.Stdout.Println(">>> CopyModule >>> sending", count, "bytes")
			fmt.Fprintf(w, "%s", byte_buffer[:count])
			count, _ = fp.Read(byte_buffer)
		}

		fmt.Fprint(w, "\x00")
	}()

	scp_command := fmt.Sprintf("scp -qrt %s/%s", me.CopyArgs["target"], filename)
	if me.Args.Sudo {
		scp_command = fmt.Sprintf("sudo -S %s", scp_command)
	}
	//fmt.Println(">>> scp command is", scp_command)

	//var b bytes.Buffer
	//session.Stdout = &b
	if err := session.Run(scp_command); err != nil {
		panic("Failed to run: " + err.Error())
	}
	//fmt.Println(">>> CopyModule >>>", b.String())

	return nil

}
