package backends

type CloudBackend interface {
	Prepare(args []string) error
	Run() error
}

var backendsList = map[string]CloudBackend{}

func Register(name string, cb CloudBackend) {
	backendsList[name] = cb
}

func GetBackend(name string) CloudBackend {
	return backendsList[name]
}

func GetBackends() []string {
	var names []string

	for name, _ := range backendsList {
		names = append(names, name)
	}

	return names
}
