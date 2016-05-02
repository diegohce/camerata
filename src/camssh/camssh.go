package camssh

import (
	"bytes"
	"cliargs"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"output"

	"golang.org/x/crypto/ssh"
)

type SshConnection struct {
	Config         *ssh.ClientConfig
	Client         *ssh.Client
	Bastion_config *ssh.ClientConfig
	Bastion        *ssh.Client
	Stdout         *output.StdoutManager
	Stderr         *output.StderrManager
}

func makeSigner(pemfile string) (ssh.Signer, error) {

	pemBytes, err := ioutil.ReadFile(pemfile)
	if err != nil {
		//log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, err
	}
	return signer, nil
}

func NewSshConnection(host string, args *cliargs.Arguments, stdout *output.StdoutManager, stderr *output.StderrManager) (*SshConnection, error) {

	var authm ssh.AuthMethod

	if args.PemFile != "" {
		signer, err := makeSigner(args.PemFile)
		if err != nil {
			return nil, err
		}
		authm = ssh.PublicKeys(signer)

	} else {
		authm = ssh.Password(args.Pass)
	}

	config := &ssh.ClientConfig{
		User: args.User,
		Auth: []ssh.AuthMethod{
			authm,
			//ssh.Password(args.Pass),
		},
	}

	/*
		config := &ssh.ClientConfig{
			User: args.User,
			Auth: []ssh.AuthMethod{
				ssh.Password(args.Pass),
			},
		}
	*/

	sshconn := &SshConnection{
		Config: config,
		Stdout: stdout,
		Stderr: stderr,
	}

	if len(args.Bastion) == 0 {
		var err error

		stdout.Println(">>> Dialing", host)
		sshconn.Client, err = ssh.Dial("tcp", host, config)
		if err != nil {
			return nil, errors.New("Failed to dial: " + err.Error())
		}

	} else {
		var err error
		var bastion_authm ssh.AuthMethod

		if args.PemFile != "" {
			signer, err := makeSigner(args.PemFile)
			if err != nil {
				return nil, err
			}
			bastion_authm = ssh.PublicKeys(signer)

		} else {
			bastion_authm = ssh.Password(args.BastionPass)
		}

		sshconn.Bastion_config = &ssh.ClientConfig{
			User: args.BastionUser,
			Auth: []ssh.AuthMethod{
				bastion_authm,
				//ssh.Password(args.BastionPass),
			},
		}

		stdout.Println(">>> Dialing bastion", args.Bastion, "with user", args.BastionUser)
		sshconn.Bastion, err = ssh.Dial("tcp", args.Bastion, sshconn.Bastion_config)
		if err != nil {
			return nil, errors.New("Failed to dial: " + err.Error())
		}

		stdout.Println(">>> Creating connection between", args.Bastion, "and", host)
		var client_tcp_conn net.Conn
		client_tcp_conn, err = sshconn.Bastion.Dial("tcp", host)
		if err != nil {
			return nil, errors.New("Failed to dial: " + err.Error())
		}

		client_conn, new_ch, req_ch, err := ssh.NewClientConn(client_tcp_conn, host, config)
		if err != nil {
			return nil, errors.New("Failed jumping from bastion to target: " + err.Error())
		}

		sshconn.Client = ssh.NewClient(client_conn, new_ch, req_ch)
	}

	return sshconn, nil

}

func (me *SshConnection) Close() {
	me.Client.Close()
	if me.Bastion != nil {
		me.Bastion.Close()
	}
}

func (me *SshConnection) WhoAmI() (string, error) {
	session, err := me.Client.NewSession()
	if err != nil {
		return "", errors.New("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		return "", errors.New("Failed to run: " + err.Error())
	}
	return b.String(), nil
}

func (me *SshConnection) SudoWhoAmI(args *cliargs.Arguments) (string, error) {
	session, err := me.Client.NewSession()
	if err != nil {
		return "", errors.New("Failed to create session: " + err.Error())
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
		return "", errors.New("Failed to run: " + err.Error())
	}
	return b.String(), nil
}
