package handler

import (
	"strconv"
	"strings"
)

func CreateRoute(name string, route string) {
	split := (strings.Split(route, ","))

	levels := []Level{}

	for _, l := range split {
		l = strings.Trim(l, " ")
		x := strings.Split(l, ":")
		var side Side

		switch strings.ToLower(x[1]) {
		case "a":
			side = SideA
		case "b":
			side = SideB
		case "c":
			side = SideC
		default:
			side = SideA
		}

		i, _ := strconv.Atoi(x[0])

		levels = append(levels, Level{Chapter(i), side})
	}

	cr := LoadFile().CustomRuns

	if cr == nil {
		cr := make(map[string][]Level)

		cr[name] = levels

		db := File{LoadFile().Settings, LoadFile().DefaultCustomsNames, LoadBule(), LoadFile().Pb, cr}

		saveConfig(db)

		return
	}

	cr[name] = levels

	db := File{LoadFile().Settings, LoadFile().DefaultCustomsNames, LoadBule(), LoadFile().Pb, cr}

	saveConfig(db)
}
