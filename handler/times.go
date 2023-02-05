package handler

import (
	"encoding/json"
	"io/ioutil"
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

		db = File{buleTimes, pb, LoadFile().Settings}

		file, _ := json.Marshal(db)
		_ = ioutil.WriteFile(path, file, 0644)

		return
	}

	run := Run{m, LoadFile().Pb["any"].Levelnames}
	pb := make(map[string]Run)
	pb["any"] = run

	buleTimes := MergeBule(m, LoadBule())

	db = File{buleTimes, pb, LoadFile().Settings}

	file, _ := json.Marshal(db)
	_ = ioutil.WriteFile(path, file, 0644)

}

func LoadFile() File {
	path := os.Getenv("HOME") + "/.config/casf/casf.json"
	// size, err : := os.Stat("/.config/casf/casf.json").Size()
	// if err != nil {
	//     log.Fatal(err)
	// }

	// if size == 0 {
	// 	fd
	// }

	var file File

	f, err := os.Open(path)
	if err != nil {
		//file = loadEmptyTimes(route)
		// log.Fatalln(err)
		return loadEmptyFile()
	}

	r := json.NewDecoder(f)
	err = r.Decode(&file)
	if err != nil {
		//m = loadEmptyTimes(route)
		// log.Fatalln(err)
		return loadEmptyFile()
	}

	return file
}

func loadEmptyTimes(route string) map[Level]time.Duration {
	var m = make(map[Level]time.Duration)

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
