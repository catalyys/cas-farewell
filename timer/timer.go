package timer

import (
	"casf/formatter"
	"casf/handler"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	tm "github.com/buger/goterm"
	"github.com/fsnotify/fsnotify"
)

var (
	pbTimes   map[handler.Level]time.Duration
	buleTimes map[handler.Level]time.Duration
	old_times map[handler.Level]time.Duration
)

func RunOverlay(file string, info bool, splits bool, routeP string, number bool, side bool) {
	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Flush()

	if handler.GetSetting("flag_i") == "true" {
		info = true
	}
	if handler.GetSetting("flag_s") == "true" {
		splits = true
	}
	if handler.GetSetting("flag_n") == "true" {
		number = true
	}
	if handler.GetSetting("flag_z") == "true" {
		side = true
	}

	var saveFile string
	switch file {
	case "1":
		saveFile = handler.GetSetting("celeste_savefolder") + "0.celeste"
	case "2":
		saveFile = handler.GetSetting("celeste_savefolder") + "1.celeste"
	default:
		saveFile = handler.GetSetting("celeste_savefolder") + "2.celeste"
	}

	buleTimes = handler.LoadBule()
	pbTimes = handler.LoadRun(routeP)
	var route = getRun(routeP)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// err = w.Add(saveFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	for w.Add(saveFile) != nil {
		printTimes(handler.LoadEmptyTimes(routeP), info, splits, routeP, number, side)
		time.Sleep(time.Second)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Printf("running %s\n", routeP)

	times := handler.ParseSaveFile(saveFile)

	printTimes(times, info, splits, routeP, number, side)
	for {
		select {
		case ev := <-w.Events:
			switch ev.Op {
			case fsnotify.Remove:
				buleTimes = handler.MergeBule(times, buleTimes)
				times = make(map[handler.Level]time.Duration)

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
                                var new_times map[handler.Level]time.Duration = nil
                                for {
				        new_times = handler.ParseSaveFile(saveFile)
                                        if new_times != nil {
                                                break
                                        }
                                }
                                times = new_times
			}

			printTimes(times, info, splits, routeP, number, side)

			_, isDone := times[route[len(route)-1]]
			if isDone {
				var d, pbD time.Duration

				for _, k := range route {
					d += times[k]
					pbD += pbTimes[k]

				}
				if d < pbD || len(pbTimes) == 0 {
					//log.Printf("new pb, congratulations!")
					pbTimes = times
					handler.SaveTimes(pbTimes, routeP)
				}
			}

		case <-c:
			//buleTimes = mergeBule(times, buleTimes)
			handler.SaveTimes(times, "bule")
			return
		}
	}
}

func ShowBest(info bool, splits bool, route string, number bool, side bool) {
	pbTimes = handler.LoadRun(route)
	buleTimes = handler.LoadBule()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Printf("PB Times in %s\n", route)
	printTimes(pbTimes, info, splits, route, number, side)
	fmt.Printf("-----------------------------------------------\n")
}

func ShowBule() {
	buleTimes = handler.LoadBule()
	fmt.Printf("%20s  %7s\n", "Chapter ", "Time")

	for _, level := range handler.AllChapter {
		d := buleTimes[level]

		if d == 0 {
			fmt.Printf("%20s     -\n", level.String(false, true))
		} else {
			fmt.Printf("%20s  %s\n", level.String(false, true), formatter.FormatWithMinutes(d))
		}
	}
}

func printTimes(times map[handler.Level]time.Duration, info bool, splits bool, routeP string, number bool, side bool) {
	oTotal := time.Duration(0)
	nTotal := time.Duration(0)

	var route = getRun(routeP)

	for _, chapter := range route {
		oTime := old_times[chapter]
		nTime := times[chapter]
		oTotal += oTime
		nTotal += nTime
	}
	// if oTotal == nTotal {
	// 	return
	// }

	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Flush()

	old_times = times

	total := time.Duration(0)
	pbTotal := time.Duration(0)
	besttotal := time.Duration(0)
	pbSplit := time.Duration(0)
	buleSplit := time.Duration(0)

	if splits {
		fmt.Printf("%20s  %7s  %7s  %7s\n", "Chapter ", "Time", "Diff", "Split")
	} else {
		fmt.Printf("%20s  %7s  %7s\n", "Chapter ", "Time", "Diff")
	}

	for _, level := range route {

		d := times[level]
		pbD := pbTimes[level]
		bD := buleTimes[level]

		total += d
		pbTotal += pbD

		if d == 0 {
			if splits {
				fmt.Printf("%20s     -        -      -\n", level.String(number, side))
			} else {
				fmt.Printf("%20s     -        -\n", level.String(number, side))
			}

			besttotal += bD
			if pbSplit == time.Duration(0) {
				pbSplit += pbD
				buleSplit += bD
			}
		} else {
			if splits {
				fmt.Printf("%20s  %s  %16s  %s\n", level.String(number, side), formatter.FormatWithMinutes(total), formatter.FormatDiff(total, pbTotal, d < bD), formatter.FormatWithMinutes(d))
			} else {
				fmt.Printf("%20s  %s  %16s\n", level.String(number, side), formatter.FormatWithMinutes(total), formatter.FormatDiff(total, pbTotal, d < bD))
			}

			besttotal += d
		}
	}
	if splits && info {
		fmt.Printf("-----------------------------------------------\n")
		fmt.Printf("%20s  %10s  %10s\n", "best possible Time", "PB Split", "best Split")
		fmt.Printf("%20s  %10s  %10s\n", formatter.FormatWithMinutes(besttotal), formatter.FormatWithMinutes(pbSplit), formatter.FormatWithMinutes(buleSplit))
	} else if info {
		fmt.Printf("---------------------------------------\n")
		fmt.Printf("%20s  %10s\n", "best possible Time", "PB Split")
		fmt.Printf("%20s  %10s\n", formatter.FormatWithMinutes(besttotal), formatter.FormatWithMinutes(pbSplit))
	}
}

func getRun(route string) []handler.Level {
	levels := handler.GetAllRoutes()[route]

	if levels == nil {
		fmt.Printf("route %s not found", route)
		fmt.Println()
		log.Fatal()
	}

	return levels
}
