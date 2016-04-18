package about

import (
	"backends"
	"fmt"
)

type aboutT struct {
	message string
}

func init() {
	backends.Register("about", &aboutT{})
}

func (be *aboutT) Prepare(args []string) error {
	return nil
}

func (be *aboutT) Run() error {
	fmt.Println("Camerata inventory tool by @diegohc")
	return nil
}
