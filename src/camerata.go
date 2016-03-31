package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	VERSION      = "0.1.0"
	VERSION_NAME = "Jake"
)

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
			fmt.Println("ArgumentsError", err)
		case CamerataError:
			fmt.Println("CamerataError", err)
		}
		os.Exit(1)
	}

	if len(args.Inventory) > 0 {
		fmt.Println("Inventory file not implemented")
		os.Exit(1)
	}

	for _, host := range strings.Split(args.Hosts, ",") {
		host = strings.TrimSpace(host)
		if strings.Index(host, ":") < 0 {
			host = host + ":22"
		}

		sshconn, err := NewSshConnection(host, args)
		if err != nil {
			panic(err.Error())
		}
		defer sshconn.Close()

		if !args.Sudo {
			result, err := sshconn.WhoAmI()
			if err != nil {
				panic(err.Error())
			} else {
				fmt.Println("WhoAmI @", host, result)
			}
		} else {
			sudo_result, err := sshconn.SudoWhoAmI(args)
			if err != nil {
				panic(err.Error())
			} else {
				fmt.Println("SudoWhoAmI @", host, sudo_result)
			}
		}

	}

}
