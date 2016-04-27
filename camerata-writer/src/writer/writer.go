package writer

type CamerataWriter interface {
	Prepare(args []string) error
	Run() error
}

var writersList = map[string]CamerataWriter{}
var writersDesc = map[string]string{}

func Register(name string, cw CamerataWriter, description string) {
	writersList[name] = cw
	writersDesc[name] = description
}

func GetWriter(name string) CamerataWriter {
	return writersList[name]
}

func GetWriters() map[string]string {
	return writersDesc
}
