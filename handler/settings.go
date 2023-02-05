package handler

import (
	"encoding/json"
	"io/ioutil"
	"os"
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

	db = File{buleTimes, pb, defaultSettings}

	file, _ := json.Marshal(db)
	_ = ioutil.WriteFile(path, file, 0644)
}
