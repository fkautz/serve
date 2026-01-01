// Copyright (c) 2014-2025 Frederick F. Kautz IV
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/urfave/cli"
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
			Name:  "log, l",
			Usage: "Log to stderr",
		},
		cli.StringFlag{
			Name:  "cert, c",
			Value: "",
			Usage: "Certificate for TLS",
		},
		cli.StringFlag{
			Name:  "key, k",
			Value: "",
			Usage: "Key for TLS",
		},
	}
	app.Action = func(c *cli.Context) error {
		dir := c.String("dir")
		address := c.String("address")
		handler := handlers.CompressHandler(http.FileServer(http.Dir(dir)))
		if c.Bool("log") {
			handler = handlers.LoggingHandler(os.Stderr, handler)
		}

		server := &http.Server{
			Addr:    address,
			Handler: handler,
		}

		// Channel to listen for shutdown signals
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		// Channel to capture server errors
		serverErr := make(chan error, 1)

		cert := c.String("cert")
		key := c.String("key")

		go func() {
			if cert != "" || key != "" {
				if cert == "" || key == "" {
					serverErr <- cli.NewExitError("Both a certificate and key must be provided for TLS", 1)
					return
				}
				serverErr <- server.ListenAndServeTLS(cert, key)
			} else {
				serverErr <- server.ListenAndServe()
			}
		}()

		select {
		case err := <-serverErr:
			if err != nil && err != http.ErrServerClosed {
				return err
			}
		case <-stop:
			log.Println("Shutting down server...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				return err
			}
			log.Println("Server stopped")
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
