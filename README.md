etcdsh
======
A command line application that provides an interactive shell for the etcd server.

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

With Go installed clone the repo and install using `go install`.

```
git clone git@github.com:headzoo/etcdsh.git
cd etcdsh
go get
go install
```

### Usage
Use `etcdsh -help` for a list of command line arguments. Once in the shell type "help" for a list of commands.


### Configuration
Default configuration values are stored in `$HOME/.etcdsh`, although command line values override the values from `.etcdsh`.

An example configuration.

```
{
  "machine": "http://127.0.0.1:4001",
  "colors": true
}
```

* "machine" The etcd server to connect to.
* "colors" Whether to use colors in output. Only applicable to Linux. The [LS_COLORS](http://blog.twistedcode.org/2008/04/lscolors-explained.html) environment variable is used to determine which colors to use.


### Custom Prompt
The terminal prompt may be customized using the environmental variables ETCDSH_PS1 and ETCDSH_PS2. The prompt format is nearly identital to the format used by bash. See [How to: Change / Setup bash custom prompt (PS1)](http://www.cyberciti.biz/tips/howto-linux-unix-bash-shell-setup-prompt.html) for a complete list of escape codes.

Example:
```
# This is the default etcdsh prompt. The prompt will look like 'sean@etcd:domains$ '.
export ETCDSH_PS1="\u@etcd:\w\$ "
```

Escape sequences not currently supported by etcdsh: \\D, \\V, \\!, \\#, \\j, \nnn, \\[, and \\]. Additionally bash commands cannot be embedded in the prompt. For example you can't use `\u@$(hostname):`.

### Bugs
This was the first app I've written using [Go](http://golang.org/). The code is probably sloppy, and breaks some Go conventions. Use the [issue system](https://github.com/headzoo/etcdsh/issues) to report bugs and broken conventions.


