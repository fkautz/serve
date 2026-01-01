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
	"github.com/urfave/cli/v3"
)

// version is set at build time via ldflags
var version = "dev"

func main() {
	app := &cli.Command{
		Name:    "serve",
		Usage:   "Simple HTTP Server",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Value:   ".",
				Usage:   "Directory to serve",
			},
			&cli.StringFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Value:   ":8080",
				Usage:   "Address to listen on",
			},
			&cli.BoolFlag{
				Name:    "log",
				Aliases: []string{"l"},
				Usage:   "Log to stderr",
			},
			&cli.StringFlag{
				Name:    "cert",
				Aliases: []string{"c"},
				Value:   "",
				Usage:   "Certificate for TLS",
			},
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Value:   "",
				Usage:   "Key for TLS",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			dir := cmd.String("dir")
			address := cmd.String("address")
			handler := handlers.CompressHandler(http.FileServer(http.Dir(dir)))
			if cmd.Bool("log") {
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

			cert := cmd.String("cert")
			key := cmd.String("key")

			go func() {
				if cert != "" || key != "" {
					if cert == "" || key == "" {
						serverErr <- cli.Exit("Both a certificate and key must be provided for TLS", 1)
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
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := server.Shutdown(shutdownCtx); err != nil {
					return err
				}
				log.Println("Server stopped")
			}

			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
