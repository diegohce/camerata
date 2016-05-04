package aboutmodule

import (
	"camssh"
	"cliargs"
	"fmt"
	"modules"
	"output"
)

type AboutModule struct {
	modules.TCamerataModule
	MyArgs map[string]string
}

func init() {
	modules.Register("about", &AboutModule{},
		"Demo module.")
}

func (me *AboutModule) Setup(args *cliargs.Arguments, stdout *output.StdoutManager, stderr *output.StderrManager) {
	me.Args = args
	me.Stdout = stdout
	me.Stderr = stderr

	me.MyArgs = modules.ModuleArgsMap(args.MArguments)

	fmt.Printf("%+v\n", me.MyArgs)

}

func (me *AboutModule) Prepare(host string, sshconn *camssh.SshConnection) error {
	me.Host = host
	me.Sshconn = sshconn

	return nil
}

func (me *AboutModule) Run() error {

	fmt.Println("source:", me.MyArgs["source"])
	fmt.Println("Target:", me.MyArgs["target"])

	return nil

}
