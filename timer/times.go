package timer

import (
	"casf/handler"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func SaveTimes(m map[handler.Level]time.Duration, typ string) {
	path := os.Getenv("HOME") + "/.config/casf/casf.json"
	//run := "any"
	fmt.Println(typ)

	run := handler.Run{m, nil}
	pb := make(map[string]handler.Run)
	pb["any"] = run

	buleTimes = mergeBule(m, Loadbule())

	db := handler.File{buleTimes, pb, "any"}

	file, _ := json.Marshal(db)
	_ = ioutil.WriteFile(path, file, 0644)
}

func LoadTimes(route string) handler.File {
	path := os.Getenv("HOME") + "/.config/casf/casf.json"
	var file handler.File

	f, err := os.Open(path)
	if err != nil {
		//file = loadEmptyTimes(route)
		fmt.Println("error read")
	}

	r := json.NewDecoder(f)
	err = r.Decode(&file)
	if err != nil {
		//m = loadEmptyTimes(route)
		fmt.Println(err)
	}

	return file
}

func loadEmptyTimes(route string) map[handler.Level]time.Duration {
	var m = make(map[handler.Level]time.Duration)

	return m
}

func Loadbule() map[handler.Level]time.Duration {
	db := LoadTimes("any")

	fmt.Println(db.Bule)

	return db.Bule
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
