package model

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Package struct {
	XMLName    xml.Name `xml:"package"`
	Name       string   `xml:"name,attr"`
	Id         string   `xml:"id,attr"`
	Version    int      `xml:"version,attr"`
	Difficulty int      `xml:"difficulty,attr"`

	Info     Info    `xml:"info"`
	Tags     Tags    `xml:"tags"`
	Authors  Authors `xml:"authors"`
	Sources  Sources `xml:"sources"`
	Comments string  `xml:"comments"`

	Rounds Rounds `xml:"rounds"`
}

type Tags struct {
	Content []string `xml:"tag"`
}

type Info struct {
	Authors Authors `xml:"authors"`
}

type Authors struct {
	Content []string `xml:"author"`
}

type Sources struct {
	Content []string `xml:"source"`
}

type Rounds struct {
	Rounds []Round `xml:"round"`
}

type Round struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`

	Themes Themes `xml:"themes"`
}

type Themes struct {
	Themes []Theme `xml:"theme"`
}

type Theme struct {
	Name      string    `xml:"name,attr"`
	Questions Questions `xml:"questions"`
}

type Questions struct {
	Questions []Question `xml:"question"`
}

type Question struct {
	Price    int      `xml:"price,attr"`
	Scenario Scenario `xml:"scenario"`
}

type Scenario struct {
	Atom Atom `xml:"atom"`
}

type Atom struct {
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

func main() {
	xmlFile, err := os.Open("./xml/test.xml")
	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := io.ReadAll(xmlFile)

	defer xmlFile.Close()

	var pack Package
	err = xml.Unmarshal(byteValue, &pack)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(pack.Rounds)
}
