package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func FirstBoot() {
	os.Mkdir(os.Getenv("HOME")+"/.config/casf", 0755)

	os.Create(os.Getenv("HOME") + "/.config/casf/casf.json")

	setDefaults()
}

func setDefaults() {
	path := os.Getenv("HOME") + "/.config/casf/casf.json"
	//run := "any"

	var db File

	pb := LoadFile().Pb

	buleTimes := LoadFile().Bule

	defaultSettings := map[string]string{
		"default_run":        "any",
		"celeste_savefolder": os.Getenv("HOME") + "/.local/share/Celeste/Saves/",
	}

	db = File{defaultSettings, buleTimes, pb}

	file, _ := json.Marshal(db)
	_ = ioutil.WriteFile(path, file, 0644)
}

func ParseSaveFile(path string) map[Level]time.Duration {
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

func ListRoutes() {
	var m = GetAllRoutes()

	fmt.Printf("%9s | %25s\n", "Route", "Chapters")
	fmt.Printf("----------|--------------------------------------\n")

	for key, value := range m {
		fmt.Printf("%9s | %25s\n", key, ListChapters(value))
	}
}

func GetSetting(key string) string {
	return LoadFile().Settings[key]
}
