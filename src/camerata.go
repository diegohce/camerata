package main

import (
	"camerata/errors"
	"fmt"
	"os"
	"strings"
)

const (
	VERSION      = "0.1.0"
	VERSION_NAME = "Jake"
)

func main() {

	args := &Arguments{}
	args.Parse()

	fmt.Printf(">>> Running Camerata v%s %s\n", VERSION, VERSION_NAME)
	fmt.Println(">>>")

	err := args.Validate()
	if err != nil {
		switch err := err.(type) {
		case errors.CamerataArgumentsError:
			fmt.Println("ArgumentsError", err)
		case errors.CamerataError:
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

		result, err := sshconn.WhoAmI()
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println("WhoAmI @", host, result)
		}

		sudo_result, err := sshconn.SudoWhoAmI(args)
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println("SudoWhoAmI @", host, sudo_result)
		}

	}

}
