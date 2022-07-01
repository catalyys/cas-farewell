package timer

import (
	"casf/handler"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func SaveTimes(m map[handler.Level]time.Duration, typ string) {
	path := os.Getenv("HOME") + "/.config/casf/casf.json"
	//run := "any"

	var db handler.File

	if typ == "bule" {
		pb := LoadFile().Pb

		buleTimes = mergeBule(m, LoadBule())

		db = handler.File{buleTimes, pb, "any"}

		file, _ := json.Marshal(db)
		_ = ioutil.WriteFile(path, file, 0644)

		return
	}

	run := handler.Run{m, nil}
	pb := make(map[string]handler.Run)
	pb["any"] = run

	buleTimes = mergeBule(m, LoadBule())

	db = handler.File{buleTimes, pb, "any"}

	file, _ := json.Marshal(db)
	_ = ioutil.WriteFile(path, file, 0644)

}

func LoadFile() handler.File {
	path := os.Getenv("HOME") + "/.config/casf/casf.json"
	var file handler.File

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

func loadEmptyTimes(route string) map[handler.Level]time.Duration {
	var m = make(map[handler.Level]time.Duration)

	return m
}

func LoadBule() map[handler.Level]time.Duration {
	db := LoadFile()

	return db.Bule
}

func LoadRun(route string) map[handler.Level]time.Duration {
	db := LoadFile()

	return db.Pb[route].Times
}

func mergeBule(old, new map[handler.Level]time.Duration) map[handler.Level]time.Duration {
	m := make(map[handler.Level]time.Duration)

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
