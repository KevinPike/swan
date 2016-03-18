package main

import (
	"os"
	"time"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "swan"
	app.Usage = "terminal migrations"
	app.Commands = []cli.Command{
		{
			Name:        "run",
			Aliases:     []string{"r"},
			Usage:       "run migrations",
			Description: "run all migrations since the last migration",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir",
					Value: "./migrations",
					Usage: "migration directory",
				},
				cli.StringFlag{
					Name:  "last",
					Value: ".swan",
					Usage: "last migration file",
				},
			},
			Action: func(c *cli.Context) {
				dir := c.String("dir")
				last := c.String("last")
				if err := Run(last, dir); err != nil {
					os.Exit(1)
				}
			},
		},
		{
			Name:        "create",
			Aliases:     []string{"c"},
			Usage:       "create a new migration: create [name] [ext]",
			Description: "create a new migration in the provided migrations directory",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir",
					Value: "./migrations",
					Usage: "migration directory",
				},
			},
			Action: func(c *cli.Context) {
				name := c.Args().First()
				ext := c.Args().Get(1)
				if ext == "" {
					ext = "sh"
				}
				dir := c.String("dir")

				if err := Create(name, dir, ext, time.Now()); err != nil {
					os.Exit(1)
				}
			},
		},
	}

	app.Run(os.Args)
}
