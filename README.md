serve
=====

A very simple static HTTP server

install
=======

```sh
go get github.com/fkautz/serve
```

usage
=====

```
NAME:
   serve - Simple HTTP Server

USAGE:
   serve [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --dir, -d '.'		Directory to serve
   --address, -a ':8080'	Address to listen on
   --log, -l			Log to stderr
   --help, -h			show help
   --version, -v		print the version
```

examples
========

```sh
# current directory port 8080
serve

# specified directory port 8080
serve -d '/opt/webpage'

# current directory port 8081
serve -a '8081'

# current directory localhost access only port 8082
serve -a '127.0.0.1:8082'

# specified directory localhost port 8081 with logging
serve -d '/opt/webpage' -a '127.0.0.1:8081' -l
```
