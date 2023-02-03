package handler

import "os"

func FirstBoot() {
	os.Mkdir(os.Getenv("HOME")+"/.config/casf", 0755)

	os.Create(os.Getenv("HOME") + "/.config/casf/casf.json")
}

// func setDefaults() {

// }
