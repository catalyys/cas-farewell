package handler

import "os"

func FirstBoot() {
	os.Mkdir(os.Getenv("HOME")+"/.config/casf", 0755)
}
