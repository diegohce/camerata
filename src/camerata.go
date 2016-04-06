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
	fmt.Printf(">>> Running Camerata v%s (%s)\n", VERSION, VERSION_NAME)
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
		inventory, err := ParseInventory(args.Inventory)
		if err != nil {
			panic(err.Error())
		}
		//fmt.Printf("%+v", inventory)

		mod, err := NewModule(args)
		if err != nil {
			panic(err.Error())
		}

		for name, server := range inventory.Servers {
			fmt.Println("Playing", name, "from inventory")

			inv_args := &Arguments{
				User:       server.User,
				Pass:       server.Password,
				Sudo:       server.Sudo,
				SudoNoPass: server.SudoNoPass,
			}
			if server.UseBastion {
				inv_args.Bastion = inventory.Bastion.Host
				inv_args.BastionUser = inventory.Bastion.User
				inv_args.BastionPass = inventory.Bastion.Password
			}
			host := server.Host

			if strings.Index(host, ":") < 0 {
				host = host + ":22"
			}

			sshconn, err := NewSshConnection(host, inv_args)
			if err != nil {
				fmt.Printf("%+v\n", inv_args)
				panic(err.Error())
			}
			defer sshconn.Close()

			if err := mod.(CamerataModule).Prepare(host, sshconn); err != nil {
				fmt.Println(">>>", err)
				os.Exit(1)
			}

			if err := mod.(CamerataModule).Run(); err != nil {
				fmt.Println(">>>", err)
				os.Exit(1)
			}
			sshconn.Close()
		}

		os.Exit(0)

	}

	askpasswords(args)

	mod, err := NewModule(args)
	if err != nil {
		panic(err.Error())
	}

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

		if err := mod.(CamerataModule).Prepare(host, sshconn); err != nil {
			fmt.Println(">>>", err)
			os.Exit(1)
		}

		if err := mod.(CamerataModule).Run(); err != nil {
			fmt.Println(">>>", err)
			os.Exit(1)
		}

	}

}
