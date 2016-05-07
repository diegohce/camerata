package aptmodule

import (
	"bytes"
	"camssh"
	"cliargs"
	"errors"
	"fmt"
	"io"
	"modules"
	"os"
	"output"
	"strconv"
	"strings"
	"time"
)

type Apt struct {
	modules.TCamerataModule
	MyArgs       map[string]string
	commands     []string
	upd_duration time.Duration
}

func init() {

	aptmodule_description := `apt operations module.
		name=pkg_name
		update_cache=yes (force) || update_cache=24h59m59s (if older than 24h59m59s, update)
		deb=/path/to/package.deb (on server)
		deb_dependencies=yes || no (try to install 'deb' dependencies if 'deb' fails	)`

	modules.Register("apt", &Apt{}, aptmodule_description)
}

func (me *Apt) Setup(args *cliargs.Arguments, stdout *output.StdoutManager, stderr *output.StderrManager) {
	me.Args = args
	me.Stdout = stdout
	me.Stderr = stderr

	me.MyArgs = modules.ModuleArgsMap(args.MArguments)

	//fmt.Printf("%+v\n", me.MyArgs)

}

func (me *Apt) Prepare(host string, sshconn *camssh.SshConnection) error {
	me.Host = host
	me.Sshconn = sshconn

	if value, ok := me.MyArgs["update_cache"]; ok {
		if value == "yes" {
			me.commands = append(me.commands, "DEBIAN_FRONTEND=noninteractive apt-get -y update")
		}

		{
			d, err := time.ParseDuration(value)
			fmt.Println("Duration", d)
			if err != nil {
				me.Stderr.Println(err, "apt operations will not update cache")
			} else {
				me.upd_duration = d
				me.commands = append(me.commands, "DEBIAN_FRONTEND=noninteractive apt-get -y update")
			}
		}
	}

	if value, ok := me.MyArgs["name"]; ok {
		me.commands = append(me.commands, fmt.Sprintf("DEBIAN_FRONTEND=noninteractive apt-get -y install %s", value))
	}

	return nil
}

func (me *Apt) Run() error {

	me.Stdout.Println(">>> AptModule >>> ", me.Args.MArguments)

	for _, command := range me.commands {

		//		fmt.Println(command)

		session, err := me.Sshconn.Client.NewSession()
		if err != nil {
			panic("Failed to create session: " + err.Error())
		}
		defer session.Close()

		var commandline string

		if me.Args.User != "root" {
			commandline = fmt.Sprintf("echo %s | sudo -S %s", me.Args.Pass, command)
		} else {
			commandline = command
		}

		r_stdin, _ := session.StdinPipe()
		defer r_stdin.Close()

		r_stdout, _ := session.StdoutPipe()
		r_stderr, _ := session.StderrPipe()

		go func() {
			io.Copy(os.Stdout, r_stdout)
			io.Copy(os.Stderr, r_stderr)
		}()

		session.Stdin = os.Stdin

		if strings.Index(commandline, "update") > -1 {
			last_update := me.getLastUpdate()
			update_time := last_update.Add(me.upd_duration)

			me.Stdout.Println(">>> AptModule >>> Last update", last_update)
			me.Stdout.Println(">>> AptModule >>> Update time", update_time)
			me.Stdout.Println(">>> AptModule >>> Now        ", time.Now())

			if time.Now().Before(update_time) {
				continue
			}
		}

		if err := session.Run(commandline); err != nil {
			return errors.New("Error " + err.Error() + " Failed to run " + command + " aborting apt operation on " + me.Host)
		}
		session.Close()

	}

	return nil

}

func (me *Apt) getLastUpdate() time.Time {
	session, err := me.Sshconn.Client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("stat -c \"%Y\" /var/lib/apt/periodic/update-success-stamp"); err != nil {
		return time.Unix(0, 0)
	}

	i, _ := strconv.Atoi(strings.TrimSuffix(b.String(), "\n"))

	//	fmt.Println("b.string", b.String())
	//	fmt.Println("i", i)

	return time.Unix(int64(i), 0)

}
