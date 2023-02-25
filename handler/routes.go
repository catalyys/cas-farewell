package handler

import (
	"regexp"
	"strconv"
	"strings"
)

func CreateRoute(name string, route string) {
	var re = regexp.MustCompile(`(?m)\d:?[aAbBcC012]`)

	levels := []Level{}

	for _, l := range re.FindAllString(route, -1) {
		l = strings.Trim(l, " ")
		x := strings.Replace(l, ":", "", -1)
		var side Side

		switch strings.ToLower(x[1:]) {
		case "a":
			side = SideA
		case "b":
			side = SideB
		case "c":
			side = SideC
		case "0":
			side = SideA
		case "1":
			side = SideB
		case "2":
			side = SideC
		default:
			side = SideA
		}

		i, _ := strconv.Atoi(x[:1])

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

func DeleteCustomRoute(key string) {
	curr := LoadFile()

	delete(curr.CustomRuns, key)

	saveConfig(curr)
}
