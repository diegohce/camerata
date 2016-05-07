# Camerata #

Simple and easy to use server orchestration set of tools made with love and golang.

# camerata command line options: #

```

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

## Inventory file format ##

Inventory files are [toml v0.2.0](https://github.com/toml-lang/toml/blob/master/versions/en/toml-v0.2.0.md) files with three main sections:

### bastion (optional) ###

Specifies connection data regarding the bastion or jumpbox server. It resembles the command line arguments:

```

-bastion string
        Bastion or jumpbox server
  -bastion-pass string
        Bastion or jumpbox server password (default: same as --pass)
  -bastion-user string
        Bastion or jumpbox server login user (default: same as --user)
```

In inventory file format

```
[bastion]
host="host_name_or_ip.to_bastion[:port]"
user="username" # to jump with
password="the_username_secret"
```

**password** and **user** are optional and, if not specified will be taken from the command line arguments.

### Servers ###

This is you servers inventory list. 

```
[servers]
  [servers.some_server_name]
  host="ip_or_hostname[:port]" #default port is 22
  
  user="username" #if it's not present , command line --user option
  
  password="user_secret_word" #if it's not present , command line --pass option or the one prompted in the console.
  
  sudo=false # or true if you want to run modules as root.
  
  sudo_nopass=false # or true if sudo does not ask for password.
	
	use_bastion=false # or true if we need to jump through the [bastion] host.

  [servers.some_other_server]
  host="other_ip_or_hostname"
  user="username"

# and so on...
```

### Modules ###

Modules are the operation unit. If there's no modules on the inventory file, **camerata** will loop through the servers list and will connect and disconnect with no further action.

Butt (double t intended ;) if you have a **[[modules]]** section on your inventory file, **camerata** will execute each one in order on every server in the servers list.

At the time of this writing, the following modules are available:

- test
- copy
- command
- apt (unstable version)

Execute camerata with *--modules* to get a list of available modules and how to use them.


Inventory file format for modules:

```
[[modules]]
name="test"
args="" #no args, it just execs "whoami" on the server.

[[modules]]
name="copy"
args="source=/path/to/my/file target=/destination/dir/on/server"

[[modules]]
name="command"
args="cat /path/to/my/file" # bash command line

```


## Available backends (inventory generators) ##

 * vmware
 * about
 * amazon (working on it...)


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

