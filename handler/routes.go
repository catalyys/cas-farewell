package handler

import (
	"strings"
)

func CreateRoute(name string, route string) {
	split := strings.Split(route, ",")

	levels := []Level{}

	for _, l := range split {
		levels = append(levels, Level{Chapter(l[0]), Side(l[1])})
	}

	cr := LoadFile().CustomRuns

	cr[name] = levels

	db := File{LoadFile().Settings, LoadFile().DefaultCustomsNames, LoadBule(), LoadFile().Pb, cr}

	saveConfig(db)
}
