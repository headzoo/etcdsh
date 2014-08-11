etcdsh
======
An interactive shell for the etcd server.

Use `etcdsh -help` for a list of command line arguments. Once in the shell type "help" for a list of commands.

Default configuration values are stored in `$HOME/.etcdsh`, which may contain the following values.

```
{
  "machine": "http://127.0.0.1:4001"
}
```

Command line values override the values from `.etcdsh`.
