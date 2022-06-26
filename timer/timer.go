package timer

import (
	"casf/handler"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

var (
	pbTimes   map[handler.Level]time.Duration
	buleTimes map[handler.Level]time.Duration
	old_times map[handler.Level]time.Duration
)

func RunOverlay(file string, info bool, splits bool, routeP string, number bool, side bool) {
	var saveFile = os.Getenv("HOME") + "/.local/share/Celeste/Saves/" + file + ".celeste"
	buleTimes = LoadTimes("bule.json")
	pbTimes = LoadTimes(getFile(routeP))
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
					//log.Printf("new pb, congratulations!")
					pbTimes = times
					SaveTimes(pbTimes)
				}
			}

		case <-c:
			buleTimes = mergeBule(times, buleTimes)
			SaveTimes(buleTimes)
			return
		}
	}
}

func ShowBest(info bool, splits bool, route string, number bool, side bool) {
	pbTimes = LoadTimes(getFile(route))
	buleTimes = LoadTimes("bule.json")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Printf("PB Times in %s\n", route)
	printTimes(pbTimes, info, splits, route, number, side)
	fmt.Printf("-----------------------------------------------\n")
}

func ShowSplits(info bool, splits bool, route string, number bool, side bool) {
	pbTimes = LoadTimes(getFile(route))
	buleTimes = LoadTimes("bule.json")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Printf("best Splits in %s\n", route)
	printTimes(buleTimes, info, splits, route, number, side)
	fmt.Printf("-----------------------------------------------\n")
}

func parseSaveFile(path string) map[handler.Level]time.Duration {
	times := make(map[handler.Level]time.Duration)

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
			if ams.TimePlayed == 0 {
				continue
			}
			times[handler.Level{area.ID, handler.Side(side)}] = time.Duration(ams.TimePlayed) * 100
		}
	}

	return times
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

func ListRoutes() {
	var m = handler.GetAllRoutes()

	fmt.Printf("%9s | %25s\n", "Route", "Chapters")
	fmt.Printf("----------|--------------------------------------\n")

	for key, value := range m {
		fmt.Printf("%9s | %25s\n", key, handler.ListChapters(value))
	}
}

func getRun(route string) []handler.Level {
	switch route {
	case "any%":
		return handler.AnyPercent
	case "any%B":
		return handler.AnyPercentB
	case "ForCity":
		return handler.City
	}

	//log.Fatal("not a valid route\n")
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

	//log.Fatal("not a valid route\n")
	return ""
}

type SaveData struct {
	xml.Name
	Areas []Area `xml:"Areas>AreaStats"`
}

type Area struct {
	ID            handler.Chapter `xml:",attr"`
	AreaModeStats []AreaModeStats `xml:"Modes>AreaModeStats"`
}

type AreaModeStats struct {
	TimePlayed uint64 `xml:",attr"` // in 10 millionths of a second
	BestTime   uint64 `xml:",attr"` // in 10 millionths of a second
}
