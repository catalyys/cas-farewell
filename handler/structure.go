package handler

import (
	"encoding/xml"
	"time"
)

type File struct {
	Settings map[string]string       `json:"settings,omitempty"`
	Bule     map[Level]time.Duration `json:"bule"`
	Pb       map[string]Run          `json:"pb,omitempty"`
}

type Run struct {
	Times      map[Level]time.Duration `json:"times"`
	Levelnames map[Level]string        `json:"level_names,omitempty"`
}

type Level struct {
	Chapter Chapter
	Side    Side
}

type Side int

type Chapter int

type SaveData struct {
	xml.Name
	Areas []Area `xml:"Areas>AreaStats"`
}

type Area struct {
	ID            Chapter         `xml:",attr"`
	AreaModeStats []AreaModeStats `xml:"Modes>AreaModeStats"`
}

type AreaModeStats struct {
	TimePlayed uint64 `xml:",attr"` // in 10 millionths of a second
	BestTime   uint64 `xml:",attr"` // in 10 millionths of a second
}
