package handler

import (
	"encoding/json"
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
	var db File

	pb := LoadFile().Pb

	buleTimes := LoadFile().Bule

	cr := LoadFile().CustomRuns

	defaultSettings := map[string]string{
		"default_run":        "any",
		"celeste_savefolder": os.Getenv("HOME") + "/.local/share/Celeste/Saves/",
	}

	defaultChapterNames := map[Level]string{
		{Chapter1, SideA}: "",
		{Chapter1, SideB}: "",
		{Chapter1, SideC}: "",
		{Chapter2, SideA}: "",
		{Chapter2, SideB}: "",
		{Chapter2, SideC}: "",
		{Chapter3, SideA}: "",
		{Chapter3, SideB}: "",
		{Chapter3, SideC}: "",
		{Chapter4, SideA}: "",
		{Chapter4, SideB}: "",
		{Chapter4, SideC}: "",
		{Chapter5, SideA}: "",
		{Chapter5, SideB}: "",
		{Chapter5, SideC}: "",
		{Chapter6, SideA}: "",
		{Chapter6, SideB}: "",
		{Chapter6, SideC}: "",
		{Chapter7, SideA}: "",
		{Chapter7, SideB}: "",
		{Chapter7, SideC}: "",
		{Chapter8, SideA}: "",
		{Chapter8, SideB}: "",
		{Chapter8, SideC}: "",
		{Chapter9, SideA}: "",
	}

	db = File{defaultSettings, defaultChapterNames, buleTimes, pb, cr}

	saveConfig(db)
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
		fmt.Println("unable to decode json")
		return loadEmptyFile()
	}

	f.Close()

	return file
}

func saveConfig(db File) {
	path := os.Getenv("HOME") + "/.config/casf/casf.json"

	file, _ := json.MarshalIndent(db, "", "    ")
	_ = ioutil.WriteFile(path, file, 0644)
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

func ImportOldPb(file string, cat string) {
	var pb map[Level]time.Duration

	f, err := os.Open(file)
	if err != nil {
		//file = loadEmptyTimes(route)
		log.Fatalln(err)
	}

	r := json.NewDecoder(f)
	err = r.Decode(&pb)
	if err != nil {
		//m = loadEmptyTimes(route)
		log.Fatalln(err)
	}

	pbs := LoadFile().Pb

	if pbs == nil {
		pbs = map[string]Run{}
	}

	run := Run{pb, pbs[cat].Levelnames}
	pbs[cat] = run

	cr := LoadFile().CustomRuns

	db := File{LoadFile().Settings, LoadFile().DefaultCustomsNames, LoadBule(), pbs, cr}

	saveConfig(db)
}

func ImportOldBule(file string) {
	var bule map[Level]time.Duration

	f, err := os.Open(file)
	if err != nil {
		//file = loadEmptyTimes(route)
		log.Fatalln(err)
	}

	r := json.NewDecoder(f)
	err = r.Decode(&bule)
	if err != nil {
		//m = loadEmptyTimes(route)
		log.Fatalln(err)
	}

	m := MergeBule(bule, LoadBule())

	db := File{LoadFile().Settings, LoadFile().DefaultCustomsNames, m, LoadFile().Pb, LoadFile().CustomRuns}

	saveConfig(db)
}
