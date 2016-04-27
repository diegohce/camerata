package filewriter

import (
	"errors"
	"flag"
	"fmt"
	"writer"
)

type arguments struct {
	File string
}

type FileWriter struct {
	Args *arguments
}

func init() {
	writer.Register("file", &FileWriter{}, "Writes stdin to --file")
}

func (me *FileWriter) Prepare(args []string) error {

	me.Args = &arguments{}

	fs := flag.NewFlagSet("FileWriter", flag.ExitOnError)
	fs.StringVar(&me.Args.File, "file", "", "Filename to write to")

	if len(args[1:]) == 0 {
		args = append(args, "--help")
	}

	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	if me.Args.File == "" {
		return errors.New("--file cannot be empty")
	}

	return nil
}

func (me *FileWriter) Run() error {
	fmt.Println("Hello from FileWriter instance")
	return nil
}
