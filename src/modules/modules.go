package modules

import (
	"camssh"
	"cliargs"
	"errors"
	"fmt"
	"output"
)

type TCamerataModule struct {
	Host    string
	Args    *cliargs.Arguments
	Sshconn *camssh.SshConnection
	Stdout  *output.StdoutManager
	Stderr  *output.StderrManager
}

type CamerataModule interface {
	Prepare(string, *camssh.SshConnection) error
	Setup(args *cliargs.Arguments, stdout *output.StdoutManager, stderr *output.StderrManager)
	Run() error
}

var ModulesList = []string{
	"test   : Runs \"whoami\" command on target hosts.",
	"command: Executes --args command line on target hosts.",
	"copy   : Sends --args source_file|dest_directory to target hosts.",
}

var AvailableModules = map[string]CamerataModule{}

func Register(name string, cm CamerataModule) {
	AvailableModules[name] = cm
}

func NewModule(args *cliargs.Arguments, stdout *output.StdoutManager, stderr *output.StderrManager) (CamerataModule, error) {

	m := AvailableModules[args.Module]
	if m == nil {
		return nil, errors.New(fmt.Sprintf("Module %s not implemented", args.Module))
	}

	m.Setup(args, stdout, stderr)

	return m, nil
}
