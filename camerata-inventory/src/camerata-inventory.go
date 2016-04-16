package main

import (
	"backends"
	_ "backends/vmware"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Choose a cloud backend:")
		for _, name := range backends.GetBackends() {
			fmt.Printf("\t%s\n", name)
		}
		os.Exit(1)
	}

	be := backends.GetBackend(os.Args[1])
	if be == nil {
		fmt.Println("Unimplemented cloud backend", os.Args[1])
		os.Exit(1)
	}

	if err := be.Prepare(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := be.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
