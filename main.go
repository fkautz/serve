package main

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gorilla/handlers"
)

func main() {
	app := cli.NewApp()
	app.Name = "serve"
	app.Usage = "Simple HTTP Server"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dir, d",
			Value: ".",
			Usage: "Directory to serve",
		},
		cli.StringFlag{
			Name:  "address, a",
			Value: ":8080",
			Usage: "Address to listen on",
		},
		cli.BoolFlag{
			Name:  "log,l",
			Usage: "Log to stderr",
		},
		cli.StringFlag{
			Name:  "cert,c",
			Value: "",
			Usage: "Certificate for TLS",
		},
		cli.StringFlag{
			Name:  "key,k",
			Value: "",
			Usage: "Key for TLS",
		},
	}
	app.Action = func(c *cli.Context) {
		address := ":8080"
		dir := "."

		if c.String("dir") != "" {
			dir = c.String("dir")
		}
		if c.String("address") != "" {
			address = c.String("address")
		}
		server := handlers.CompressHandler(http.FileServer(http.Dir(dir)))
		if c.Bool("log") {
			server = handlers.LoggingHandler(os.Stderr, server)
		}
		cert := c.String("cert")
		key := c.String("key")
		// if cert or key are set, start TLS
		if cert != "" || key != "" {
			// Both cert and key must be set
			if cert == "" || key == "" {
				log.Fatalln("Both a certificate and key must be provided for TLS")
			}
			log.Fatalln(http.ListenAndServeTLS(address, cert, key, server))
		} else {
			// No cert and no key, just serve unencrypted
			log.Fatalln(http.ListenAndServe(address, server))
		}
	}

	app.Run(os.Args)
}
