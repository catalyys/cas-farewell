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

func (l Level) String() string {
	switch l.Chapter {
	case Prologue:
		fallthrough
	case Epilogue:
		return l.Chapter.String()
	default:
		return l.Chapter.String() + l.Side.String()
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

func LongChapterName(c Chapter) string {
	return longChapterName[c]
}

func (c Chapter) String() string {
	return longChapterName[c]
}

type Side int

const (
	SideA Side = iota
	SideB
	SideC
)

var sideName = []string{"A", "B", "C"}

func (s Side) String() string {
	return ""
}
