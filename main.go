package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli"
)

var (
	pbTimes   map[Level]time.Duration
	buleTimes map[Level]time.Duration
	old_times map[Level]time.Duration
)

func main() {
	color.NoColor = false

	app := cli.NewApp()
	app.Name = "Celeste Auto Splitter Farewell"
	app.Usage = "Farewell"
	app.Version = "0.9"
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

func runOverlay(file string, info bool, splits bool, routeP string, number bool, side bool) {
	var saveFile = os.Getenv("HOME") + "/.local/share/Celeste/Saves/" + file + ".celeste"
	buleTimes = loadTimes("bule.json", routeP)
	pbTimes = loadTimes(getFile(routeP), routeP)
	var route = getRun(routeP)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = w.Add(saveFile)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Printf("running %s\n", routeP)

	times := parseSaveFile(saveFile)

	printTimes(times, info, splits, routeP, number, side)
	for {
		select {
		case ev := <-w.Events:
			switch ev.Op {
			case fsnotify.Remove:
				buleTimes = mergeBule(times, buleTimes)
				times = make(map[Level]time.Duration)

				f, err := os.OpenFile(saveFile, os.O_CREATE, 0644)
				if err != nil {
					log.Fatal(err)
				}
				f.Close()
				err = w.Add(saveFile)
				if err != nil {
					log.Fatal(err)
				}
			case fsnotify.Chmod:
				fallthrough
			case fsnotify.Write:
				times = parseSaveFile(saveFile)
			}

			printTimes(times, info, splits, routeP, number, side)
			_, isDone := times[route[len(route)-1]]
			if isDone {
				var d, pbD time.Duration

				for _, k := range route {
					d += times[k]
					pbD += pbTimes[k]

				}
				if d < pbD {
					log.Printf("new pb, congratulations!")
					pbTimes = times
					saveTimes(pbTimes, getFile(routeP))
				}
			}

		case <-c:
			buleTimes = mergeBule(times, buleTimes)
			saveTimes(buleTimes, "bule.json")
			return
		}
	}
}

func showBest(info bool, splits bool, route string, number bool, side bool) {
	pbTimes = loadTimes(getFile(route), route)
	buleTimes = loadTimes("bule.json", route)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Printf("PB Times in %s\n", route)
	printTimes(pbTimes, info, splits, route, number, side)
	fmt.Printf("-----------------------------------------------\n")
}

func showSplits(info bool, splits bool, route string, number bool, side bool) {
	pbTimes = loadTimes(getFile(route), route)
	buleTimes = loadTimes("bule.json", route)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Printf("best Splits in %s\n", route)
	printTimes(buleTimes, info, splits, route, number, side)
	fmt.Printf("-----------------------------------------------\n")
}

func parseSaveFile(path string) map[Level]time.Duration {
	times := make(map[Level]time.Duration)

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := xml.NewDecoder(f)

	var s SaveData
	err = d.Decode(&s)

	if err != nil {
		fmt.Fprintf(os.Stderr, "corrupted or missing savefile!\n")
		log.Fatal(err)
	}

	for _, area := range s.Areas {
		for side, ams := range area.AreaModeStats {
			if ams.BestTime == 0 {
				continue
			}
			times[Level{area.ID, Side(side)}] = time.Duration(ams.BestTime) * 100
		}
	}

	return times
}

func saveTimes(m map[Level]time.Duration, path string) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("failed to open\n")
	}

	w := json.NewEncoder(f)
	err = w.Encode(m)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("failed to save\n")
	}
}

func loadTimes(path string, route string) map[Level]time.Duration {
	var m map[Level]time.Duration

	f, err := os.Open(path)
	if err != nil {
		m = loadEmptyTimes(route)
	}

	r := json.NewDecoder(f)
	err = r.Decode(&m)
	if err != nil {
		m = loadEmptyTimes(route)
	}

	return m
}

func loadEmptyTimes(route string) map[Level]time.Duration {
	var m = make(map[Level]time.Duration)

	for _, value := range getRun(route) {
		m[value] = time.Duration(24 * time.Hour)
	}

	return m
}

func printTimes(times map[Level]time.Duration, info bool, splits bool, routeP string, number bool, side bool) {
	oTotal := time.Duration(0)
	nTotal := time.Duration(0)

	var route = getRun(routeP)

	for _, chapter := range route {
		oTime := old_times[chapter]
		nTime := times[chapter]
		oTotal += oTime
		nTotal += nTime
	}
	if oTotal == nTotal {
		return
	}
	old_times = times

	total := time.Duration(0)
	pbTotal := time.Duration(0)
	besttotal := time.Duration(0)
	pbSplit := time.Duration(0)
	buleSplit := time.Duration(0)

	if splits {
		fmt.Printf("%20s  %7s  %7s  %7s\n", "Chapter", "Time", "Diff", "Split")
	} else {
		fmt.Printf("%20s  %7s  %7s\n", "Chapter", "Time", "Diff")
	}

	for _, chapter := range route {
		d := times[chapter]
		pbD := pbTimes[chapter]
		bD := buleTimes[chapter]

		total += d
		pbTotal += pbD

		if d == 0 {
			if splits {
				fmt.Printf("%20s     -      -       -\n", chapter.String(number, side))
			} else {
				fmt.Printf("%20s     -      -\n", chapter.String(number, side))
			}

			besttotal += bD
			if pbSplit == time.Duration(0) {
				pbSplit += pbD
				buleSplit += bD
			}
		} else {
			if splits {
				fmt.Printf("%20s  %s  %16s  %s\n", chapter.String(number, side), formatWithMinutes(total), formatDiff(total-pbTotal, d < bD), formatWithMinutes(d))
			} else {
				fmt.Printf("%20s  %s  %16s\n", chapter.String(number, side), formatWithMinutes(total), formatDiff(total-pbTotal, d < bD))
			}

			besttotal += d
		}
	}
	if splits && info {
		fmt.Printf("-----------------------------------------------\n")
		fmt.Printf("%20s  %10s  %10s\n", "best possible Time", "PB Split", "best Split")
		fmt.Printf("%20s  %10s  %10s\n", formatWithMinutes(besttotal), formatWithMinutes(pbSplit), formatWithMinutes(buleSplit))
	} else if info {
		fmt.Printf("---------------------------------------\n")
		fmt.Printf("%20s  %10s\n", "best possible Time", "PB Split")
		fmt.Printf("%20s  %10s\n", formatWithMinutes(besttotal), formatWithMinutes(pbSplit))
	}
}

func formatWithMinutes(d time.Duration) string {
	minutes := d / time.Minute

	tenths := d / (100 * time.Millisecond)
	seconds := d / time.Second

	tenths %= 10
	seconds %= 60

	return fmt.Sprintf("%02d:%02d.%01d", minutes, seconds, tenths)
}

func formatDiff(d time.Duration, isBule bool) string {
	var sign byte
	var sprintf func(string, ...interface{}) string
	if d < 0 {
		sign = '-'
		d = -d
		sprintf = color.New(color.FgGreen).SprintfFunc()
	} else if d < 100*time.Millisecond {
		sign = '±'
		sprintf = color.New(color.FgGreen).SprintfFunc()
	} else { // at least 100ms difference
		sign = '+'
		sprintf = color.New(color.FgRed).SprintfFunc()
	}

	if isBule {
		sprintf = color.New(color.FgYellow).SprintfFunc()
	}

	tenths := d / (100 * time.Millisecond)
	seconds := d / time.Second
	minutes := d / time.Minute

	tenths %= 10
	seconds %= 60

	if d >= 1*time.Minute {
		return sprintf("%c%d:%02d.%01d", sign, minutes, seconds, tenths)
	} else {
		return sprintf("%c%02d.%01d", sign, seconds, tenths)
	}

}

func mergeBule(old, new map[Level]time.Duration) map[Level]time.Duration {
	m := make(map[Level]time.Duration)

	for k, v := range old {
		m[k] = v
		w, ok := new[k]
		if !ok {
			continue
		} else if w < v {
			m[k] = w
		}
	}

	for k, v := range new {
		m[k] = v
		w, ok := old[k]
		if !ok {
			continue
		} else if w < v {
			m[k] = w
		}
	}

	return m
}

func listRoutes() {
	var m = getAllRoutes()

	fmt.Printf("%9s  %25s\n", "Route", "Chapters")
	fmt.Printf("------------------------------------------------\n")

	for key, value := range m {
		fmt.Printf("%9s  %25s\n", key, listChapters(value))
	}
}

func getRun(route string) []Level {
	switch route {
	case "any%":
		return anyPercent
	case "any%B":
		return anyPercentB
	case "ForCity":
		return City
	}

	log.Fatal("not a valid route\n")
	return nil
}

func getFile(route string) string {
	switch route {
	case "any%":
		return "pb.json"
	case "any%B":
		return "any%B.json"
	case "ForCity":
		return "city.json"
	}

	log.Fatal("not a valid route\n")
	return ""
}

type SaveData struct {
	xml.Name
	Areas []Area `xml:"Areas>AreaStats"`
}

type Area struct {
	ID            Chapter         `xml:",attr"`
	AreaModeStats []AreaModeStats `xml:"Modes>AreaModeStats"`
}

type AreaModeStats struct {
	TimePlayed uint64 `xml:",attr"` // in 10 millionths of a second
	BestTime   uint64 `xml:",attr"` // in 10 millionths of a second
}
