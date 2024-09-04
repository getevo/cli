package main

import (
	"fmt"
	"github.com/getevo/cli"
)

func main() {

	cli.Register(
		cli.Command{
			Switch: "client",
			Usage:  "Rum application in client mode",
			Action: func(c *cli.Context) {
				fmt.Println("This is an example command in client mode")
			},
		},
		cli.Command{
			Switch: "server",
			Usage:  "Rum application in server mode",
			Params: cli.Params{
				{Switch: "port", Usage: "Server port", Required: true},
				{Switch: "host", Usage: "Server host", DefaultValue: 5060},
			},
			Action: func(c *cli.Context) {
				c.Print("This is an example command in server mode")
				c.Print("port: %s", c.Params.Get("port").String())
				c.Print("host: %d", c.Params.Get("host").Int())
			},
			Commands: []cli.Command{
				{
					Switch: "start",
					Usage:  "Start the server",
					Params: cli.Params{
						{Switch: "name", Usage: "Sever name"},
					},
					Action: func(c *cli.Context) {
						c.Print("Starting the server...")
						c.Print("server name: %s", c.Param("name").String())
					},
				},
			},
		},
	)

	cli.Run()
}
