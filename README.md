etcdsh v0.2
===========
A command line application that provides an interactive shell for the etcd server.

*Aug 13th, 2014: This is still a work in progress, and right now only supports the commands cd, ls, get, and set.*

* [Overview](#overview)
* [Installation](#installation)
* [Configuration](#configuration)
* [Custom Prompt](#custom-prompt)
* [TODO](#todo)
* [Bugs](#bugs)


### Overview
Using etcdsh you can browse your etcd "filesystem" like an ordinary filesystem, using familar commands line "cd" and "ls". For example:


```
Connecting to http://127.0.0.1:4001
Interactive etcd shell started.
Type 'help' for a list of commands.
Type 'q' to quit.

joe@etcd:/$ ls
total 4
0 0   0  k: .
0 0   0  k: ..
3 3 300  k: go
1 1   0  k: apps

joe@etcd:/$ cd apps
joe@etcd:/apps$ ls
total 4
 0  0 0  k: .
 0  0 0  k: ..
 6  6 0  o: mobile
15 15 0  o: website

joe@etcd:/apps$ get website
http://example.com

joe@etcd:/apps$ set website http://example.net
http://example.net

joe@etcd:/apps cd ..
joe@etcd:/$ q
```


### Installation
Go version 1.3 is required. See [this page](http://golang.org/doc/install) for Go installation instructions.

The readline development libraries are required.

```
apt-get install libreadline-dev
```

Clone the repo and install using `go install`.

```
git clone git@github.com:headzoo/etcdsh.git
cd etcdsh
go get
go install
```

Use `etcdsh -help` for a list of command line arguments. Once in the shell type "help" for a list of commands.


### Configuration
There are three ways to configure etcdsh:

1. Using command line arguments.
2. Using environment variables.
3. Using a JSON file saved at `$HOME/.etcdsh`.

The following is the list of configuration options.

* "machine" The etcd server to connect to.
* "colors" Whether to use colors in output. Only applicable to Linux. The [LS_COLORS](http://blog.twistedcode.org/2008/04/lscolors-explained.html) environment variable is used to determine which colors to use.
* "ps1" The first custom prompt.
* "ps2" The second custom prompt.

When used at the command line, prefix the option with "-", eg `-machine`. When defined as an environment variable, prefix the option with "ETCDSH_", eg `ETCDSH_MACHINE`.


An example configuration file.

```
{
  "machine": "http://127.0.0.1:4001",
  "colors": true
}
```


### Custom Prompt
The terminal prompt may be customized using the environmental variables ETCDSH_PS1 and ETCDSH_PS2. The prompt format is nearly identital to the format used by bash. See [How to: Change / Setup bash custom prompt (PS1)](http://www.cyberciti.biz/tips/howto-linux-unix-bash-shell-setup-prompt.html) for a complete list of escape codes.

Example:
```
# This is the default etcdsh prompt. The prompt will look like 'sean@etcd:domains$ '.
export ETCDSH_PS1="\u@etcd:\w\$ "
```

You can even embed color codes.

```
export ETCDSH_PS1="\e[36m\u@\h\e[0m:\e[32m\w\e[0m"
```

Escape sequences not currently supported by etcdsh: \\D, \\V, \\!, \\#, \\j, \nnn, \\[, and \\]. Additionally bash commands cannot be embedded in the prompt. For example you can't use `\u@$(hostname):`.


### TODO
* Write the command history to a file, eg `$HOME/.etcdsh_history`.
* Auto complete needs to work recursivly.
* Process command switches.
* Handle pipes and redirection.


### Bugs
This was the first app I've written using [Go](http://golang.org/). The code is probably sloppy, and breaks some Go conventions. Use the [issue system](https://github.com/headzoo/etcdsh/issues) to report bugs and broken conventions.

