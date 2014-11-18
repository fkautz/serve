serve
=====

A very simple static HTTP server

usage
=====

```sh
# current directory port 8080
serve

# specified directory port 8080
serve -d '/opt/webpage'

# current directory port 8081
serve -a '8081'

# current directory localhost access only port 8082
serve -a '127.0.0.1:8082'

# specified directory localhost port 8081
serve -d '/opt/webpage' -a '127.0.0.1:8081'
```
