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
		fmt.Print(">>> Bastion Password [Enter for user pass]: ")
		bastion_password_b, err := terminal.ReadPassword(0)
		fmt.Println("")
		if err != nil {
			return CamerataError{"Error reading bastion password from terminal"}
		}

		bastion_password := string(bastion_password_b)
		if bastion_password != "" {
			args.BastionPass = bastion_password
		} else {
			args.BastionPass = args.Pass
		}

	} else {
		args.BastionPass = args.Pass

	}

	return nil

}

func main() {

	args := &Arguments{}
	args.Parse()

	stdout := NewStdoutManager(args)
	stderr := NewStderrManager(args)

	stdout.Println(">>> Hi!")
	stdout.Printf(">>> Running Camerata v%s (%s)\n", VERSION, VERSION_NAME)
	stdout.Println(">>>")

	err := args.Validate()
	if err != nil {
		switch err := err.(type) {
		case CamerataArgumentsError:
			stderr.Println(">>> ArgumentsError", err)
		case CamerataError:
			stderr.Println(">>> CamerataError", err)
		}
		os.Exit(1)
	}

	if len(args.Inventory) > 0 {
		inventory, err := ParseInventory(args.Inventory)
		if err != nil {
			panic(err.Error())
		}

		if strings.Index(inventory.Bastion.Host, ":") < 0 {
			inventory.Bastion.Host = inventory.Bastion.Host + ":22"
		}

		passwords_saved := false

		for name, server := range inventory.Servers {
			stdout.Println(">>> Playing", name, "from inventory")

			if server.User == "" && args.User == "" {
				stderr.Println(">>>", CamerataArgumentsError{"No user defined on inventory file nor --user argument for " + name})
				continue
			}

			if !passwords_saved {

				if inventory.Bastion.Password == "" {
					args.AskBastionPass = true
				}
				if server.Password == "" && args.Pass == "" {
					args.AskPass = true
				} else {
					args.AskPass = false
				}

				askpasswords(args)
				inventory.Bastion.Password = args.BastionPass
				passwords_saved = true
			}

			if server.User == "" {
				server.User = args.User
			}
			if server.Password == "" {
				server.Password = args.Pass
			}

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

			sshconn, err := NewSshConnection(host, inv_args, stdout, stderr)
			if err != nil {
				stderr.Printf("%+v\n", inv_args)
				panic(err.Error())
			}
			defer sshconn.Close()

			for _, module := range inventory.Modules {
				inv_args.Module = module.Name
				inv_args.MArguments = module.Args

				mod, err := NewModule(inv_args, stdout, stderr)
				if err != nil {
					panic(err.Error())
				}

				if err := mod.(CamerataModule).Prepare(host, sshconn); err != nil {
					stderr.Println(">>>", err)
					os.Exit(1)
				}

				if err := mod.(CamerataModule).Run(); err != nil {
					stderr.Println(">>>", err)
					os.Exit(1)
				}
			}
			sshconn.Close()
		}

		os.Exit(0)

	}

	askpasswords(args)

	//	var mod2 CamerataModule2
	mod, err := NewModule(args, stdout, stderr)
	if err != nil {
		panic(err.Error())
		//		mod2 = AvailableModules2[args.Module]
		//		mod2.Run()
		//		os.Exit(1)
	}

	for _, host := range strings.Split(args.Hosts, ",") {
		host = strings.TrimSpace(host)
		if strings.Index(host, ":") < 0 {
			host = host + ":22"
		}

		sshconn, err := NewSshConnection(host, args, stdout, stderr)
		if err != nil {
			stderr.Printf("%+v\n", args)
			panic(err.Error())
		}
		defer sshconn.Close()

		if err := mod.(CamerataModule).Prepare(host, sshconn); err != nil {
			stderr.Println(">>>", err)
			os.Exit(1)
		}

		if err := mod.(CamerataModule).Run(); err != nil {
			stderr.Println(">>>", err)
			os.Exit(1)
		}

	}

}
