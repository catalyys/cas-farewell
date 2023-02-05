package main

import (
	"casf/handler"
	"casf/timer"
	"errors"
	"os"
)

func main() {
	if _, err := os.Stat(os.Getenv("HOME") + "/.config/casf"); errors.Is(err, os.ErrNotExist) {
		handler.FirstBoot()
	}

	handler.FirstBoot()

	timer.StartTimer()
}
