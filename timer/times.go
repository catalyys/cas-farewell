package timer

import (
	"casf/handler"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func SaveTimes(m map[handler.Level]time.Duration, typ string) {
	path := os.Getenv("HOME") + "/.config/casf/casf.json"
	//run := "any"
	fmt.Println(typ)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("failed to open\n")
	}

	// w := json.NewEncoder(f)
	// err = w.Encode(m)

	run := handler.Run{m, nil}
	pb := make(map[string]handler.Run)
	pb["any"] = run
	db := handler.File{m, pb, "any"}

	w := json.NewEncoder(f)
	err = w.Encode(db)

	if err != nil {
		log.Fatal(err)
		fmt.Printf("failed to save\n")
	}
}

func LoadTimes(route string) map[handler.Level]time.Duration {
	path := os.Getenv("HOME") + "/.config/casf/casf.json"
	var m map[handler.Level]time.Duration

	f, err := os.Open(path)
	if err != nil {
		m = loadEmptyTimes(route)
	}

	r := json.NewDecoder(f)
	err = r.Decode(&m)
	if err != nil {
		m = loadEmptyTimes(route)
	}

	return m
}

func loadEmptyTimes(route string) map[handler.Level]time.Duration {
	var m = make(map[handler.Level]time.Duration)

	return m
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
