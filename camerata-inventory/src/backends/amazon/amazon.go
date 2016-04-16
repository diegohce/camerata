// https://aws.amazon.com/documentation/sdk-for-go/

package amazon

import (
	"backends"
	"errors"
	"flag"
)

type amazonArguments struct {
	User         string
	Host         string
	Password     string
	Insecure     bool
	OutputFormat string
}

type cloudAmazon struct {
	args *amazonArguments
}

func init() {

	args := &amazonArguments{}

	cloud := &cloudAmazon{
		args: args,
	}
	backends.Register("amazon", cloud)

}

func (be *cloudAmazon) Prepare(args []string) error {

	fs := flag.NewFlagSet("amazon", flag.ExitOnError)

	fs.StringVar(&be.args.User, "user", "", "vCenter username")
	fs.StringVar(&be.args.Password, "pass", "", "vCenter password")
	fs.StringVar(&be.args.Host, "host", "", "vCenter host[:port] or ip[:port]")
	fs.BoolVar(&be.args.Insecure, "insecure", true, "Don't check server certificate")
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

func (be *cloudAmazon) Run() error {
	return errors.New("Not implemented (...yet! ;)")
}
