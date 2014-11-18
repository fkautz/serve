package main

import (
	"github.com/codegangsta/cli"
	"log"
	"net/http"
	"os"
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
		log.Fatal(http.ListenAndServe(address, http.FileServer(http.Dir(dir))))
	}

	app.Run(os.Args)
}
