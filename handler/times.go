package handler

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"time"
)

func SaveTimes(m map[Level]time.Duration, typ string) {
	var db File

	if typ == "bule" {
		pb := LoadFile().Pb

		buleTimes := MergeBule(m, LoadBule())

		db = File{LoadFile().Settings, LoadFile().DefaultCustomsNames, buleTimes, pb}

		saveConfig(db)

		return
	}

	run := Run{m, LoadFile().Pb["any"].Levelnames}
	pb := make(map[string]Run)
	pb["any"] = run

	buleTimes := MergeBule(m, LoadBule())

	db = File{LoadFile().Settings, LoadFile().DefaultCustomsNames, buleTimes, pb}

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
		fmt.Fprintf(os.Stderr, "corrupted or missing savefile!\n")
		log.Fatal(err)
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
	levels := GetAllRoutes()[route]

	var m = make(map[Level]time.Duration)

	for _, i := range levels {
		m[i] = 1 * time.Second
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
