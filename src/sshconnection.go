package main

import (
	"bytes"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

type SshConnection struct {
	config         *ssh.ClientConfig
	client         *ssh.Client
	bastion_config *ssh.ClientConfig
	bastion        *ssh.Client
}

func NewSshConnection(host string, args *Arguments) (*SshConnection, error) {

	config := &ssh.ClientConfig{
		User: args.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(args.Pass),
		},
	}

	sshconn := &SshConnection{config: config}

	if len(args.Bastion) == 0 {
		var err error

		fmt.Println(">>> Dialing", host)
		sshconn.client, err = ssh.Dial("tcp", host, config)
		if err != nil {
			return nil, CamerataConnectionError{"Failed to dial: " + err.Error()}
		}

	} else {
		var err error

		sshconn.bastion_config = &ssh.ClientConfig{
			User: args.BastionUser,
			Auth: []ssh.AuthMethod{
				ssh.Password(args.BastionPass),
			},
		}

		fmt.Println(">>> Dialing bastion", args.Bastion, "with user", args.BastionUser)
		sshconn.bastion, err = ssh.Dial("tcp", args.Bastion, sshconn.bastion_config)
		if err != nil {
			return nil, CamerataConnectionError{"Failed to dial: " + err.Error()}
		}

		fmt.Println(">>> Creating connection between", args.Bastion, "and", host)
		var client_tcp_conn net.Conn
		client_tcp_conn, err = sshconn.bastion.Dial("tcp", host)
		if err != nil {
			return nil, CamerataConnectionError{"Failed to dial: " + err.Error()}
		}

		client_conn, new_ch, req_ch, err := ssh.NewClientConn(client_tcp_conn, host, config)
		if err != nil {
			return nil, CamerataConnectionError{"Failed jumping from bastion to target: " + err.Error()}
		}

		sshconn.client = ssh.NewClient(client_conn, new_ch, req_ch)
	}

	return sshconn, nil

}

func (me *SshConnection) Close() {
	me.client.Close()
	if me.bastion != nil {
		me.bastion.Close()
	}
}

func (me *SshConnection) WhoAmI() (string, error) {
	session, err := me.client.NewSession()
	if err != nil {
		return "", CamerataRunError{"Failed to create session: " + err.Error()}
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		return "", CamerataRunError{"Failed to run: " + err.Error()}
	}
	return b.String(), nil
}

func (me *SshConnection) SudoWhoAmI(args *Arguments) (string, error) {
	session, err := me.client.NewSession()
	if err != nil {
		return "", CamerataRunError{"Failed to create session: " + err.Error()}
	}
	defer session.Close()

	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintln(w, args.Pass)
	}()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("sudo -S /usr/bin/whoami"); err != nil {
		return "", CamerataRunError{"Failed to run: " + err.Error()}
	}
	return b.String(), nil
}
