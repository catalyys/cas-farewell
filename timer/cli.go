package timer

import (
	"casf/handler"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func StartTimer() {
	color.NoColor = false
	// handler.FirstBoot()

	app := cli.NewApp()
	app.Name = "Celeste Auto Splitter"
	app.Usage = "Farewell"
	app.Version = "1.0.0"
	app.UseShortOptionHandling = true

	myFlags := []cli.Flag{
		cli.BoolFlag{Name: "splits, s", Usage: "shows more information on splits"},
		cli.BoolFlag{Name: "info, i", Usage: "shows best possible times"},
		cli.BoolFlag{Name: "number, n", Usage: "displays numbers instead of names"},
		cli.BoolFlag{Name: "sides, z", Usage: "displays chapter side next to the name/number"},
		cli.StringFlag{
			Name: "saveslot, slot, save",
			// Value: "3",
			Usage: "indicates the saveslot `slot` 1, 2 or 3",
		},
		cli.StringFlag{
			Name:  "route, r",
			Value: "any",
			Usage: "indicates the route/run",
		},
	}
	app.Flags = myFlags

	importFlags := []cli.Flag{
		cli.BoolFlag{Name: "pb, p"},
		cli.BoolFlag{Name: "bule, b"},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "filepath of the pb or bule file to import",
		},
		cli.StringFlag{
			Name:  "run",
			Value: "any",
			Usage: "name of the run to import the pb",
		},
	}

	routeFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "name, n",
			Usage: "name of the custom run",
		},
		cli.StringFlag{
			Name:     "route, r",
			Usage:    "route of the custom run",
			Required: true,
		},
	}

	routeDeleteFlags := []cli.Flag{
		cli.StringFlag{
			Name:     "name, n",
			Usage:    "name of the custom `run` to delete",
			Required: true,
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.String("saveslot") == "" {
			RunOverlay(handler.GetSetting("default_saveslot"), c.Bool("info"), c.Bool("splits"), strings.ToLower(c.String("route")), c.Bool("number"), c.Bool("sides"))
			return nil
		}

		if c.String("saveslot") != "1" && c.String("saveslot") != "2" && c.String("saveslot") != "3" {
			fmt.Printf("saveslot needs to be 1, 2 or 3\n")
			return nil
		}

		RunOverlay(c.String("saveslot"), c.Bool("info"), c.Bool("splits"), strings.ToLower(c.String("route")), c.Bool("number"), c.Bool("sides"))

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
						ShowBest(c.Bool("info"), c.Bool("splits"), strings.ToLower(c.String("route")), c.Bool("number"), c.Bool("sides"))
						return nil
					},
				},
				{
					Name:  "bule",
					Usage: "show all best times for each chapter",
					Flags: myFlags,
					Action: func(c *cli.Context) error {
						ShowBule()
						return nil
					},
				},
				{
					Name:  "routes",
					Usage: "show all pre configured routes",
					Flags: myFlags,
					Action: func(c *cli.Context) error {
						handler.ListRoutes()
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
				if c.String("saveslot") == "" {
					RunOverlay(handler.GetSetting("default_saveslot"), c.Bool("info"), c.Bool("splits"), strings.ToLower(c.String("route")), c.Bool("number"), c.Bool("sides"))
					return nil
				}

				if c.String("saveslot") != "1" && c.String("saveslot") != "2" && c.String("saveslot") != "3" {
					fmt.Printf("saveslot needs to be 1, 2 or 3\n")
					return nil
				}

				RunOverlay(c.String("saveslot"), c.Bool("info"), c.Bool("splits"), strings.ToLower(c.String("route")), c.Bool("number"), c.Bool("sides"))

				return nil
			},
		},
		{
			Name:    "import",
			Aliases: []string{"i "},
			Usage:   "import pre v1 pb and bule files",
			Flags:   importFlags,
			Action: func(c *cli.Context) error {
				if c.Bool("pb") {
					handler.ImportOldPb(c.String("file"), strings.ToLower(c.String("run")))
					return nil
				} else if c.Bool("bule") {
					handler.ImportOldBule(c.String("file"))
					return nil
				}
				return nil
			},
		},
		{
			Name:  "route",
			Usage: "configure routes",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "create route",
					Flags: routeFlags,
					Action: func(c *cli.Context) error {
						handler.CreateRoute(c.String("name"), c.String("route"))
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove route",
					Flags: routeDeleteFlags,
					Action: func(c *cli.Context) error {
						handler.DeleteCustomRoute(c.String("name"))
						return nil
					},
				},
				{
					Name:  "show",
					Usage: "show routes",
					// Flags: routeDeleteFlags,
					Action: func(c *cli.Context) error {
						handler.ListRoutes()
						return nil
					},
				},
			},
		},
	}

	// start application
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
