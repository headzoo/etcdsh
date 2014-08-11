etcdsh
======
An interactive shell for the etcd server.


### Installation
Go version 1.3 is required. See [this page](http://golang.org/doc/install) for Go installation instructions.

With Go installed clone the repo and install using `go install`.

```
git clone git@github.com:headzoo/etcdsh.git
cd etcdsh
go install
```

### Usage
Use `etcdsh -help` for a list of command line arguments. Once in the shell type "help" for a list of commands.


### Configuration
Default configuration values are stored in `$HOME/.etcdsh`, which may contain the following values.

```
{
  "machine": "http://127.0.0.1:4001"
}
```

Command line values override the values from `.etcdsh`.
