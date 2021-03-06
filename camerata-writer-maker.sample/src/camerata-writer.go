package main

import (
	"fmt"
	"os"
	"writer"
	_ "writer/filewriter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Choose a writer:")
		for name, description := range writer.GetWriters() {
			fmt.Printf("\t%s : %s\n", name, description)
		}
		os.Exit(1)
	}

	wr := writer.GetWriter(os.Args[1])
	if wr == nil {
		fmt.Println("Unimplemented writer", os.Args[1])
		os.Exit(1)
	}

	for i := 0; i < 10; i++ {
		wri := wr.New()

		fmt.Printf("wr = %p\nwri = %p\n", wr, wri)

		if err := wri.Prepare(os.Args[1:]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := wri.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}
}
