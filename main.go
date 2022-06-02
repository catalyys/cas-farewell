package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func main() {
	color.NoColor = false

	app := cli.NewApp()
	app.Name = "Celeste Auto Splitter Farewell"
	app.Usage = "Farewell"
	app.Version = "0.9.6"
	app.UseShortOptionHandling = true

	myFlags := []cli.Flag{
		cli.BoolFlag{Name: "splits, s"},
		cli.BoolFlag{Name: "info, i"},
		cli.BoolFlag{Name: "number, n"},
		cli.BoolFlag{Name: "sides, z"},
		cli.StringFlag{
			Name:  "savefile, save",
			Value: "2",
			Usage: "indicates the savefile slot `0`, 1 or 2",
		},
		cli.StringFlag{
			Name:  "route, r",
			Value: "any%",
			Usage: "indicates the route/run",
		},
	}
	app.Flags = myFlags

	app.Action = func(c *cli.Context) error {
		if c.String("savefile") != "0" && c.String("savefile") != "1" && c.String("savefile") != "2" {
			fmt.Printf("savefile needs to be 0, 1 or 2\n")
			return nil
		}
		runOverlay(c.String("savefile"), c.Bool("info"), c.Bool("splits"), c.String("route"), c.Bool("number"), c.Bool("sides"))
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "Show best splits or peronal best time",
			Subcommands: []cli.Command{
				{
					Name:  "best",
					Usage: "show personal best",
					Flags: myFlags,
					Action: func(c *cli.Context) error {
						showBest(c.Bool("info"), c.Bool("splits"), c.String("route"), c.Bool("number"), c.Bool("sides"))
						return nil
					},
				},
				{
					Name:  "splits",
					Usage: "show best splits",
					Flags: myFlags,
					Action: func(c *cli.Context) error {
						showSplits(c.Bool("info"), c.Bool("splits"), c.String("route"), c.Bool("number"), c.Bool("sides"))
						return nil
					},
				},
				{
					Name:  "routes",
					Usage: "show all pre configured routes",
					Flags: myFlags,
					Action: func(c *cli.Context) error {
						listRoutes()
						return nil
					},
				},
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "start the overlay for the run",
			Flags:   myFlags,
			Action: func(c *cli.Context) error {
				if c.String("savefile") != "0" && c.String("savefile") != "1" && c.String("savefile") != "2" {
					fmt.Printf("savefile needs to be 0, 1 or 2\n")
					return nil
				}
				runOverlay(c.String("savefile"), c.Bool("info"), c.Bool("splits"), c.String("route"), c.Bool("number"), c.Bool("side"))
				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "for dev testing",
			Flags:   myFlags,
			Action: func(c *cli.Context) error {
				//printTimes(loadEmptyTimes("any%"), true, true, "any%", false, true)
				//saveTimes(loadEmptyTimes("any%"), "test.json")
				listRoutes()
				return nil
			},
		},
	}

	// start our application
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
