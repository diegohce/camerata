package writer

type CamerataWriter interface {
	Prepare(args []string) error
	Run() error
}

type CamerataWriterMaker interface {
	New() CamerataWriter
}

var writersList = map[string]CamerataWriterMaker{}
var writersDesc = map[string]string{}

func Register(name string, cw CamerataWriterMaker, description string) {
	writersList[name] = cw
	writersDesc[name] = description
}

func GetWriter(name string) CamerataWriterMaker {
	return writersList[name]
}

func GetWriters() map[string]string {
	return writersDesc
}
