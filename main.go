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
		log.Fatal(http.ListenAndServe(address, server))
	}

	app.Run(os.Args)
}
