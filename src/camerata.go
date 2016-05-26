package main

import (
	"camssh"
	"cliargs"
	"errors"
	"fmt"
	"modules"
	_ "modules/aboutmodule"
	_ "modules/aptmodule"
	_ "modules/commandmodule"
	_ "modules/copymodule"
	_ "modules/pipmodule"
	_ "modules/testmodule"
	"os"
	"output"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	VERSION      = "1.2.2"
	VERSION_NAME = "Jake"
)

func askpasswords(args *cliargs.Arguments) error {

	if args.AskPass {
		fmt.Fprint(os.Stderr, ">>> User Password: ")
		password_b, err := terminal.ReadPassword(0)
		fmt.Fprintln(os.Stderr, "")
		if err != nil {
			return errors.New("Error reading password from terminal")
		}

		password := string(password_b)
		args.Pass = password
	}

	if args.Bastion != "" && args.AskBastionPass {
		fmt.Fprint(os.Stderr, ">>> Bastion Password [Enter for user pass]: ")
		bastion_password_b, err := terminal.ReadPassword(0)
		fmt.Fprintln(os.Stderr, "")
		if err != nil {
			return errors.New("Error reading bastion password from terminal")
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

	args := &cliargs.Arguments{}
	args.Parse()

	stdout := output.NewStdoutManager(args)
	stderr := output.NewStderrManager(args)

	stdout.Println(">>> Hi!")
	stdout.Printf(">>> Running Camerata v%s (%s)\n", VERSION, VERSION_NAME)
	stdout.Println(">>>")

	if args.ModulesList {
		fmt.Println(strings.Join(modules.ModulesList, "\n"))
		os.Exit(0)
	}

	err := args.Validate()
	if err != nil {
		stderr.Println(">>>", err)
		os.Exit(1)
	}

	if len(args.Inventory) > 0 {
		inventory, err := ParseInventory(args.Inventory)
		if err != nil {
			stderr.Println(">>>", err)
			os.Exit(1)
		}

		if strings.Index(inventory.Bastion.Host, ":") < 0 {
			inventory.Bastion.Host = inventory.Bastion.Host + ":22"
		}

		passwords_saved := false

		for name, server := range inventory.Servers {
			stdout.Println(">>> Playing", name, "from inventory")

			if server.User == "" && args.User == "" {
				stderr.Println(">>>", errors.New("No user defined on inventory file nor --user argument for "+name))
				continue
			}

			if !passwords_saved && server.PemFile == "" {

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

			var pemfile string

			if args.PemFile != "" {
				pemfile = args.PemFile
			} else {
				pemfile = server.PemFile
			}

			if args.Sudo {
				server.Sudo = true
			}

			inv_args := &cliargs.Arguments{
				User:       server.User,
				Pass:       server.Password,
				Sudo:       server.Sudo,
				SudoNoPass: server.SudoNoPass,
				PemFile:    pemfile,
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

			//fmt.Printf("%+v\n\n", inv_args)

			sshconn, err := camssh.NewSshConnection(host, inv_args, stdout, stderr)
			if err != nil {
				//stderr.Printf("%+v\n", inv_args)
				stderr.Println(">>>", err)
				continue
			}
			defer sshconn.Close()

			for _, module := range inventory.Modules {
				inv_args.Module = module.Name
				inv_args.MArguments = module.Args

				mod, err := modules.NewModule(inv_args, stdout, stderr)
				if err != nil {
					stderr.Println(">>> Error", err, "with module", module, "args", inv_args)
				}

				if err := mod.Prepare(host, sshconn); err != nil {
					stderr.Println(">>>", err)
					continue
				}

				if err := mod.Run(); err != nil {
					stderr.Println(">>>", err)
					continue
				}
			}
			sshconn.Close()
		}

		os.Exit(0)

	}

	askpasswords(args)

	mod, err := modules.NewModule(args, stdout, stderr)
	if err != nil {
		panic(err.Error())
	}

	for _, host := range strings.Split(args.Hosts, ",") {
		host = strings.TrimSpace(host)
		if strings.Index(host, ":") < 0 {
			host = host + ":22"
		}

		sshconn, err := camssh.NewSshConnection(host, args, stdout, stderr)
		if err != nil {
			stderr.Println(">>>", err)
			//stderr.Printf(">>> %+v\n", args)
			continue
		}
		defer sshconn.Close()

		if err := mod.Prepare(host, sshconn); err != nil {
			stderr.Println(">>>", err)
			os.Exit(1)
		}

		if err := mod.Run(); err != nil {
			stderr.Println(">>>", err)
			os.Exit(1)
		}

	}

}
