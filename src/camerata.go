package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	VERSION      = "0.1.0"
	VERSION_NAME = "Jake"
)

func askpasswords(args *Arguments) error {

	if args.AskPass {
		fmt.Print(">>> User Password: ")
		password_b, err := terminal.ReadPassword(0)
		fmt.Println("")
		if err != nil {
			return CamerataError{"Error reading password from terminal"}
		}

		password := string(password_b)
		args.Pass = password
	}

	if args.AskBastionPass {
		fmt.Print(">>> Bastion Password: ")
		bastion_password_b, err := terminal.ReadPassword(0)
		fmt.Println("")
		if err != nil {
			return CamerataError{"Error reading bastion password from terminal"}
		}

		bastion_password := string(bastion_password_b)
		args.BastionPass = bastion_password

	} else {
		args.BastionPass = args.Pass

	}

	return nil

}

func main() {

	fmt.Println(">>> Hi!")
	fmt.Printf(">>> Running Camerata v%s codename %s\n", VERSION, VERSION_NAME)
	fmt.Println(">>>")

	args := &Arguments{}
	args.Parse()

	err := args.Validate()
	if err != nil {
		switch err := err.(type) {
		case CamerataArgumentsError:
			fmt.Println(">>> ArgumentsError", err)
		case CamerataError:
			fmt.Println(">>> CamerataError", err)
		}
		os.Exit(1)
	}

	if len(args.Inventory) > 0 {
		fmt.Println(">>> Inventory file not implemented.")
		fmt.Println(">>> Use --hosts in the meantime.")
		os.Exit(1)
	}

	askpasswords(args)

	for _, host := range strings.Split(args.Hosts, ",") {
		host = strings.TrimSpace(host)
		if strings.Index(host, ":") < 0 {
			host = host + ":22"
		}

		sshconn, err := NewSshConnection(host, args)
		if err != nil {
			fmt.Printf("%+v\n", args)
			panic(err.Error())
		}
		defer sshconn.Close()

		if args.Test {
			if !args.Sudo {
				result, err := sshconn.WhoAmI()
				if err != nil {
					panic(err.Error())
				} else {
					fmt.Println(">>> WhoAmI @", host, result)
				}
			} else {
				sudo_result, err := sshconn.SudoWhoAmI(args)
				if err != nil {
					panic(err.Error())
				} else {
					fmt.Println(">>> SudoWhoAmI @", host, sudo_result)
				}
			}
		}

	}

}
