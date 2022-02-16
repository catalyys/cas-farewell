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
)

var (
	saveFile  = "/home/olli/.local/share/Celeste/Saves/0.celeste"
	pbTimes   map[Level]time.Duration
	buleTimes map[Level]time.Duration
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Fprintf(os.Stderr, "Usage: %s [pbfile]\n", os.Args[0])
			fmt.Fprintln(os.Stderr, "  if no file is provided then it will autosplit the bottom")
			fmt.Fprintln(os.Stderr, "  save for the any% run configured in types.go")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "  if pbFile is provided this program will parse your pb and")
			fmt.Fprintln(os.Stderr, "  initialize your PB and gold splits in a format that this program understands")
			return
		}
		fmt.Fprintf(os.Stderr, "parsing pb save file")
		pbTimes = parseSaveFile(os.Args[1])
		saveTimes(pbTimes, "pb.json")
		_, err := os.Stat("bule.json")
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "bule.json missing, creating...")
			saveTimes(pbTimes, "pb.json")
		}
	}

	color.NoColor = false
	pbTimes = loadTimes("pb.json")
	buleTimes = loadTimes("bule.json")

	fmt.Fprintf(os.Stderr, "read pb and best times")

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = w.Add(saveFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stderr, "added save file to watched files\n")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	times := parseSaveFile(saveFile)
	fmt.Fprintf(os.Stderr, "parsed current save file\n")

	fmt.Printf("PB Times\n")
	printTimes(pbTimes, false)
	fmt.Printf("-----------------------------------------------\n")
	fmt.Printf("best Splits\n")
	printTimes(buleTimes, false)
	fmt.Printf("-----------------------------------------------\n")

	printTimes(times, true)
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

			printTimes(times, true)
			_, isDone := times[anyPercent[len(anyPercent)-1]]
			if isDone {
				var d, pbD time.Duration

				for _, k := range anyPercent {
					d += times[k]
					pbD += pbTimes[k]

				}

				if d < pbD {
					log.Printf("new pb, congratulations!")
					pbTimes = times
					saveTimes(pbTimes, "pb.json")
				}
			}

		case <-c:
			fmt.Fprintf(os.Stderr, "writing bule times\n")
			buleTimes = mergeBule(times, buleTimes)
			saveTimes(buleTimes, "bule.json")
			return

		//default:
		//	times = parseSaveFile(saveFile)
		//	printTimes(times)
		//	time.Sleep(1 * time.Second)
		}
	}
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

func printTimes(times map[Level]time.Duration, inf bool) {
	total := time.Duration(0)
	pbTotal := time.Duration(0)
	besttotal := time.Duration(0)
	pbSplit := time.Duration(0)
	buleSplit := time.Duration(0)

	fmt.Printf("%20s  %7s  %7s  %7s\n", "Chapter", "Time", "Diff", "Split")

	for _, chapter := range anyPercent {
		d := times[chapter]
		pbD := pbTimes[chapter]
		bD := buleTimes[chapter]

		total += d
		pbTotal += pbD

		if d == 0 {
			fmt.Printf("%20s     -      -       -\n", chapter)

			besttotal += bD
			if pbSplit == time.Duration(0){
				pbSplit += pbD
				buleSplit += bD
			}
		} else {
			fmt.Printf("%20s  %s  %16s  %s\n", chapter, formatWithMinutes(total), formatDiff(total-pbTotal, d < bD), formatWithMinutes(d))

			besttotal += d
		}
	}
	if inf {
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

	if d >= 1*time.Minute{
		return sprintf("%c%d:%02d.%01d", sign, minutes, seconds, tenths)
	}else{
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
