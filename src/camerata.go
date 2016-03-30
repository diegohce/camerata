package main

import (
	"camerata/errors"
	"fmt"
	"os"
)

func main() {

	args := &Arguments{}
	args.Parse()

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

	fmt.Printf("%+v\n", args)

}
