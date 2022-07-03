package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func SaveTimes(m map[Level]time.Duration, typ string) {
	path := os.Getenv("HOME") + "/.config/casf/casf.json"
	//run := "any"

	var db File

	if typ == "bule" {
		pb := LoadFile().Pb

		buleTimes := MergeBule(m, LoadBule())

		db = File{buleTimes, pb, "any"}

		file, _ := json.Marshal(db)
		_ = ioutil.WriteFile(path, file, 0644)

		return
	}

	run := Run{m, LoadFile().Pb["any"].Levelnames}
	pb := make(map[string]Run)
	pb["any"] = run

	buleTimes := MergeBule(m, LoadBule())

	db = File{buleTimes, pb, "any"}

	file, _ := json.Marshal(db)
	_ = ioutil.WriteFile(path, file, 0644)

}

func LoadFile() File {
	path := os.Getenv("HOME") + "/.config/casf/casf.json"
	var file File

	f, err := os.Open(path)
	if err != nil {
		//file = loadEmptyTimes(route)
		log.Fatalln(err)
	}

	r := json.NewDecoder(f)
	err = r.Decode(&file)
	if err != nil {
		//m = loadEmptyTimes(route)
		log.Fatalln(err)
	}

	return file
}

func loadEmptyTimes(route string) map[Level]time.Duration {
	var m = make(map[Level]time.Duration)

	return m
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
