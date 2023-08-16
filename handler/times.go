package handler

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func SaveTimes(m map[Level]time.Duration, typ string) {
	var db File

	cr := LoadFile().CustomRuns
	pb := LoadFile().Pb

	if typ == "bule" {

		buleTimes := MergeBule(m, LoadBule())

		db = File{LoadFile().Settings, LoadFile().DefaultCustomsNames, buleTimes, pb, cr}

		saveConfig(db)

		return
	}

        if pb == nil {
                pb = make(map[string]Run)
        }

	run := Run{m, LoadFile().Pb[typ].Levelnames}
	pb[typ] = run

	buleTimes := MergeBule(m, LoadBule())

	db = File{LoadFile().Settings, LoadFile().DefaultCustomsNames, buleTimes, pb, cr}

	saveConfig(db)
}

func ParseSaveFile(path string) map[Level]time.Duration {
	times := make(map[Level]time.Duration)

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "missing savefile!\n")
		log.Fatal(err)
	}
	defer f.Close()

	d := xml.NewDecoder(f)

	var s SaveData
	err = d.Decode(&s)

	if err != nil {
                return nil
	}

	for _, area := range s.Areas {
		for side, ams := range area.AreaModeStats {
			if ams.TimePlayed == 0 {
				continue
			}
			times[Level{Chapter: area.ID, Side: Side(side)}] = time.Duration(ams.TimePlayed) * 100
		}
	}

	return times
}

func LoadEmptyTimes(route string) map[Level]time.Duration {
	levels := GetAllRoutes()[strings.ToLower(route)]

	var m = make(map[Level]time.Duration)

	for _, i := range levels {
		m[i] = time.Duration(0)
	}

	return m
}

func loadEmptyFile() File {
	// var file = make(map[Level]time.Duration)
	var f1 File

	return f1
}

func LoadBule() map[Level]time.Duration {
	db := LoadFile()

	return db.Bule
}

func LoadRun(route string) map[Level]time.Duration {
	db := LoadFile()

	return db.Pb[route].Times
}

func MergeBule(old, new map[Level]time.Duration) map[Level]time.Duration {
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
