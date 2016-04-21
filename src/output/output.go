package output

import (
	"cliargs"
	"fmt"
	"os"
)

type StdoutManager struct {
	write bool
}

type StderrManager struct {
	write bool
}

func NewStdoutManager(args *cliargs.Arguments) *StdoutManager {
	return &StdoutManager{write: !args.Quiet}
}

func (me *StdoutManager) Print(v ...interface{}) (int, error) {
	if me.write {
		return fmt.Print(v...)
	}
	return 0, nil
}

func (me *StdoutManager) Println(v ...interface{}) (int, error) {
	if me.write {
		return fmt.Println(v...)
	}
	return 0, nil
}

func (me *StdoutManager) Printf(format string, v ...interface{}) (int, error) {
	if me.write {
		return fmt.Printf(format, v...)
	}
	return 0, nil
}

func NewStderrManager(args *cliargs.Arguments) *StderrManager {
	return &StderrManager{write: !args.Quiet}
}

func (me *StderrManager) Print(v ...interface{}) (int, error) {
	if me.write {
		return fmt.Fprint(os.Stderr, v...)
	}
	return 0, nil
}

func (me *StderrManager) Println(v ...interface{}) (int, error) {
	if me.write {
		return fmt.Fprintln(os.Stderr, v...)
	}
	return 0, nil
}

func (me *StderrManager) Printf(format string, v ...interface{}) (int, error) {
	if me.write {
		return fmt.Fprintf(os.Stderr, format, v...)
	}
	return 0, nil
}
