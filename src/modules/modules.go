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

var ModulesList = []string{}

var AvailableModules = map[string]CamerataModule{}

func Register(name string, cm CamerataModule, description string) {
	AvailableModules[name] = cm
	ModulesList = append(ModulesList, fmt.Sprintf("%s\n\t%s", name, description))
}

func NewModule(args *cliargs.Arguments, stdout *output.StdoutManager, stderr *output.StderrManager) (CamerataModule, error) {

	m := AvailableModules[args.Module]
	if m == nil {
		return nil, errors.New(fmt.Sprintf("Module %s not implemented", args.Module))
	}

	m.Setup(args, stdout, stderr)

	return m, nil
}
