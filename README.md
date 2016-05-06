# Camerata #

Simple and easy to use server orchestration set of tools made with love and golang.

# camerata command line options: #

```
#!bash

  -args string
    	Module arguments
  -ask-bastion-pass
    	Asks for password on the command line for bastion jump
  -ask-pass
    	Asks for password on the command line (default true)
  -bastion string
    	Bastion or jumpbox server
  -bastion-pass string
    	Bastion or jumpbox server password (default: same as --pass)
  -bastion-user string
    	Bastion or jumpbox server login user (default: same as --user)
  -hosts string
    	Comma separated hosts list
  -inventory string
    	Inventory file
  -module string
    	Module to run (default "test")
  -modules
    	List available modules
  -pass string
    	Use this password
  -pem string
    	Path to pemfile for auth
  -quiet
    	No camerata output
  -sudo
    	Run as sudo
  -sudo-nopass
    	Run as sudo without pass
  -test
    	Runs whoami on remote host
  -user string
    	Login user

```

# camerata-inventory #

Camerata inventory file generator.

## Available backends ##

 * vmware
 * about
 * amazon


## camerata-inventory **about** command line options: ##

No options, just an about message.


## camerata-inventory **vmware** command line options: ##


```
#!bash

 -bastion string
    	Bastion or jumpbox server (name or ip address)
  -bastion-nets string
    	Comma separated list of segments that uses --bastion (e.g.: 10.54.165.,10.54.170.)
  -format string
    	Output format: toml, csv (default "toml")
  -host string
    	vCenter host[:port] or ip[:port]
  -insecure
    	Don't check server certificate (default true)
  -pass string
    	vCenter password
  -user string
    	vCenter username
```

## camerata-inventory **amazon** command line options: ##

Backend into development stage

