package formatter

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func FormatDiff(t time.Duration, pb time.Duration, isBule bool) string {
	d := t - pb

	var sign byte
	var sprintf func(string, ...interface{}) string
	if d < 0 {
		sign = '-'
		d = -d
		sprintf = color.New(color.FgGreen).SprintfFunc()
	} else if d < 100*time.Millisecond {
		sign = '±'
		sprintf = color.New(color.FgGreen).SprintfFunc()
	} else { // at least 100ms difference
		sign = '+'
		sprintf = color.New(color.FgRed).SprintfFunc()
	}

	if pb == 0 {
		sprintf = color.New(color.FgGreen).SprintfFunc()
	}

	if isBule {
		sprintf = color.New(color.FgYellow).SprintfFunc()
	}

	tenths := d / (100 * time.Millisecond)
	seconds := d / time.Second
	minutes := d / time.Minute

	tenths %= 10
	seconds %= 60

	if d >= 1*time.Minute {
		return sprintf("%c%d:%02d.%01d", sign, minutes, seconds, tenths)
	} else {
		return sprintf("%c%02d.%01d", sign, seconds, tenths)
	}

}

func FormatWithMinutes(d time.Duration) string {
	hours := d / time.Hour
	minutes := (d / time.Minute) - (60 * hours)

	tenths := d / (100 * time.Millisecond)
	seconds := d / time.Second

	tenths %= 10
	seconds %= 60

	return fmt.Sprintf("%02d:%02d:%02d.%01d", hours, minutes, seconds, tenths)
}
