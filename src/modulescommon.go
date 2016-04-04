package main

import (
	"fmt"
)

type TCamerataModule struct {
	host    string
	args    *Arguments
	sshconn *SshConnection
}

type CamerataModule interface {
	Prepare(string, *SshConnection) error
	Run() error
}

//func NewModule(host string, args *Arguments) (interface{}, error) {
func NewModule(args *Arguments) (interface{}, error) {

	switch args.Module {
	case "test":
		{
			m := &TestModule{args: args}
			return m, nil
		}
	case "command":
		{
			m := &CommandModule{args: args}
			return m, nil
		}

	case "copy":
		{
			m := &CopyModule{args: args}
			return m, nil
		}

	default:
		return nil, CamerataModuleError{fmt.Sprintf("Module %s not implemented", args.Module)}
	}

}
