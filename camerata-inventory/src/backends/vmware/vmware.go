package vmware

import (
	"backends"
	"errors"
	"flag"
)

type vmwArguments struct {
	User         string
	Host         string
	Password     string
	Insecure     bool
	Bastion      string
	BastionNets  string
	OutputFormat string
}

type cloudVmWare struct {
	args *vmwArguments
}

func init() {

	args := &vmwArguments{}

	cloud := &cloudVmWare{
		args: args,
	}
	backends.Register("vmware", cloud)

}

func (be *cloudVmWare) Prepare(args []string) error {

	fs := flag.NewFlagSet("vmware", flag.ExitOnError)

	fs.StringVar(&be.args.User, "user", "", "vCenter username")
	fs.StringVar(&be.args.Password, "pass", "", "vCenter password")
	fs.StringVar(&be.args.Host, "host", "", "vCenter host[:port] or ip[:port]")
	fs.BoolVar(&be.args.Insecure, "insecure", true, "Don't check server certificate")
	fs.StringVar(&be.args.Bastion, "bastion", "", "Bastion or jumpbox server (name or ip address)")
	fs.StringVar(&be.args.BastionNets, "bastion-nets", "", "Comma separated list of segments that uses --bastion (e.g.: 10.54.165.,10.54.170.)")
	fs.StringVar(&be.args.OutputFormat, "format", "toml", "Output format: toml, csv")

	if len(args[1:]) == 0 {
		args = append(args, "--help")
	}

	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	if be.args.User == "" {
		return errors.New("--user cannot be empty")
	}
	if be.args.Host == "" {
		return errors.New("--host cannot be empty")
	}

	return nil
}

func (be *cloudVmWare) Run() error {
	return be.vmwareInventory()
}
