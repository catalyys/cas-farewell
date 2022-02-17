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
	saveFile  = os.Getenv("HOME") + "/.local/share/Celeste/Saves/0.celeste"
	pbTimes   map[Level]time.Duration
	buleTimes map[Level]time.Duration
)

func main() {
	color.NoColor = false

	app := cli.NewApp()
	app.Name = "Celeste Auto Splitter"
	app.Usage = "automatically splits the time for your Celeste run"

	// We'll be using the same flag for all our commands
	// so we'll define it up here
	myFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "view, v",
			Value: "best",
			Usage: "shows you personal best or best splits",
		},
		cli.StringFlag{
			Name:  "savefile, save, s",
			Value: "2",
			Usage: "indicates the savefile slot `0`, `1` or `2`",
		},
	}

	// we create our commands
	app.Commands = []cli.Command{
		{
			Name:    "show",
			Aliases: []string{"s"},
			Usage:   "Show best splits or peronal best time",
			Flags:   myFlags,
			// the action, or code that will be executed when
			// we execute our `show` command
			Subcommands: []cli.Command{
				{
					Name:  "best",
					Usage: "show personal best",
					Action: func(c *cli.Context) error {
						showBest()
						return nil
					},
				},
				{
					Name:  "splits",
					Usage: "show best splits",
					Action: func(c *cli.Context) error {
						showSplits()
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
			// the action, or code that will be executed when
			// we execute our `show` command
			Action: func(c *cli.Context) error {
				if c.String("savefile") != "0" && c.String("savefile") != "1" && c.String("savefile") != "2" {
					fmt.Printf("savefile needs to be 0, 1 or 2\n")
					return nil
				}
				runOverlay(c.String("savefile"))
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

func runOverlay(file string) {
	saveFile = os.Getenv("HOME") + "/.local/share/Celeste/Saves/" + file + ".celeste"
	buleTimes = loadTimes("bule.json")
	pbTimes = loadTimes("pb.json")

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

	times := parseSaveFile(saveFile)

	printTimes(times, 3)
	//fmt.Fprintf(os.Stderr, "starting loop, press ^C to exit\n")
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

			printTimes(times, 2)

		case <-c:
			buleTimes = mergeBule(times, buleTimes)
			saveTimes(buleTimes, "bule.json")
			return
		}
	}
}

func showBest() {
	pbTimes = loadTimes("pb.json")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Printf("PB Times\n")
	printTimes(pbTimes, 1)
	fmt.Printf("-----------------------------------------------\n")
}

func showSplits() {
	buleTimes = loadTimes("bule.json")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Printf("best Splits\n")
	printTimes(buleTimes, 1)
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
			times[Level{area.ID, Side(side)}] = time.Duration(ams.TimePlayed) * 100
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

func loadTimes(path string) map[Level]time.Duration {
	var m map[Level]time.Duration

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	r := json.NewDecoder(f)
	err = r.Decode(&m)
	if err != nil {
		log.Fatal(err)
	}

	return m
}

func printTimes(times map[Level]time.Duration, view int) {
	total := time.Duration(0)
	pbTotal := time.Duration(0)
	besttotal := time.Duration(0)
	pbSplit := time.Duration(0)
	buleSplit := time.Duration(0)

	if view == 1 || view == 3 {
		fmt.Printf("%20s  %7s  %7s  %7s\n", "Chapter", "Time", "Diff", "Split")
	} else {
		fmt.Printf("%20s  %7s  %7s\n", "Chapter", "Time", "Diff")
	}

	for _, chapter := range anyPercent {
		d := times[chapter]
		pbD := pbTimes[chapter]
		bD := buleTimes[chapter]

		total += d
		pbTotal += pbD

		if d == 0 {
			if view == 1 || view == 3 {
				fmt.Printf("%20s     -      -       -\n", chapter)
			} else {
				fmt.Printf("%20s     -      -\n", chapter)
			}

			besttotal += bD
			if pbSplit == time.Duration(0) {
				pbSplit += pbD
				buleSplit += bD
			}
		} else {
			if view == 1 || view == 3 {
				fmt.Printf("%20s  %s  %16s  %s\n", chapter, formatWithMinutes(total), formatDiff(total-pbTotal, d < bD), formatWithMinutes(d))
			} else {
				fmt.Printf("%20s  %s  %16s\n", chapter, formatWithMinutes(total), formatDiff(total-pbTotal, d < bD))
			}

			besttotal += d
		}
	}
	if view == 2 {
		fmt.Printf("---------------------------------------\n")
		fmt.Printf("%20s  %10s\n", "best possible Time", "PB Split")
		fmt.Printf("%20s  %10s\n", formatWithMinutes(besttotal), formatWithMinutes(pbSplit))
	} else if view == 3 {
		fmt.Printf("-----------------------------------------------\n")
		fmt.Printf("%20s  %10s  %10s\n", "best possible Time", "PB Split", "best Split")
		fmt.Printf("%20s  %10s  %10s\n", formatWithMinutes(besttotal), formatWithMinutes(pbSplit), formatWithMinutes(buleSplit))
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
