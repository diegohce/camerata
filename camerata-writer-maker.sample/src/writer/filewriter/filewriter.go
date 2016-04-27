package filewriter

import (
	"fmt"
	"writer"
)

type FileWriter struct {
	filename string
}

func init() {
	writer.Register("file", &FileWriter{}, "Writes stdin to --file")
}

func (me *FileWriter) New() writer.CamerataWriter {
	return &FileWriter{}
}

func (me *FileWriter) Prepare(args []string) error {
	return nil
}

func (me *FileWriter) Run() error {
	fmt.Println("Hello from FileWriter instance")
	return nil
}
