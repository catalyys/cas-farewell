package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Level struct {
	Chapter Chapter
	Side    Side
}

func (l Level) String(number bool, side bool) string {
	switch l.Chapter {
	case Prologue:
		fallthrough
	case Epilogue:
		return l.Chapter.String(number)
	default:
		if number {
			return l.Chapter.String(number) + l.Side.String(side)
		} else {
			return l.Chapter.String(number) + " " + l.Side.String(side)
		}
	}
}

func (l Level) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%d:%d", l.Chapter, l.Side)), nil
}

func (l *Level) UnmarshalText(b []byte) error {
	s := strings.Split(string(b), ":")
	if len(s) != 2 {
		return fmt.Errorf("malformed level: %s", b)
	}

	c, err := strconv.Atoi(s[0])
	if err != nil {
		return err
	}
	if c < 0 || c >= 10 {
		return fmt.Errorf("invalid chapter: %v", c)
	}
	(*l).Chapter = Chapter(c)

	side, err := strconv.Atoi(s[1])
	if err != nil {
		return err
	}
	if side < 0 || side >= 3 {
		return fmt.Errorf("invalid side: %v", side)
	}
	(*l).Side = Side(side)

	return nil
}

var anyPercent = []Level{
	{Prologue, SideA},
	{Chapter1, SideA},
	{Chapter2, SideA},
	{Chapter3, SideA},
	{Chapter4, SideA},
	{Chapter5, SideA},
	{Chapter6, SideA},
	{Chapter7, SideA},
}

var City = []Level{
	{Chapter1, SideA},
}

var anyPercentB = []Level{
	{Prologue, SideA},
	{Chapter1, SideA},
	{Chapter2, SideA},
	{Chapter3, SideA},
	{Chapter4, SideA},
	{Chapter5, SideA},
	{Chapter5, SideB},
	{Chapter6, SideA},
	{Chapter6, SideB},
	{Chapter7, SideA},
}

func getAllRoutes() map[string][]Level {
	var allRoutes = make(map[string][]Level)

	allRoutes["any%"] = anyPercent
	allRoutes["any%B"] = anyPercentB
	allRoutes["ForCity"] = City

	return allRoutes
}

func listChapters(levels []Level) string {
	var s string
	var i int = 0

	for _, value := range levels {
		s = s + value.String(true, true)
		if i < len(levels)-1 {
			s = s + "->"
		}
		i++
	}
	return fmt.Sprint(s)
}

type Chapter int

const (
	Prologue = iota
	Chapter1
	Chapter2
	Chapter3
	Chapter4
	Chapter5
	Chapter6
	Chapter7
	Epilogue
	Chapter8
	Chapter9
)

var shortChapterName = []string{
	"Prologue",
	"1",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"Epilogue",
	"8",
	"9",
}

var longChapterName = []string{
	"Prologue",
	"Forsaken City",
	"Old Ruins",
	"Celestial Resort",
	"Golden Ridge",
	"Mirror Temple",
	"Reflection",
	"The Summit",
	"Epilogue",
	"Core",
	"Farewell",
}

func (c Chapter) String(number bool) string {
	if number {
		return shortChapterName[c]
	} else {
		return longChapterName[c]
	}
}

type Side int

const (
	SideA Side = iota
	SideB
	SideC
)

var sideName = []string{"A", "B", "C"}

func (s Side) String(side bool) string {
	if side {
		return sideName[s]
	} else {
		return ""
	}
}
