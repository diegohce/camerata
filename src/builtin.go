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
	Run(*SshConnection) error
}

func NewModule(host string, args *Arguments) (interface{}, error) {
	//(*CamerataModule, error) {

	switch args.Module {
	case "test":
		{
			m := &TestModule{host: host, args: args}
			return m, nil
		}

	default:
		return nil, CamerataModuleError{fmt.Sprintf("Module %s not implemented", args.Module)}
	}

}
