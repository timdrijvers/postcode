package bagparser

/*
 * City offers parsing capabilities for WPL (Woonplaats) BAG files
 */
import (
	"encoding/xml"
	"fmt"
	"io"
	"maps"
)

type Cities struct {
	cities map[string]string
}

func NewCitiesParser() *Cities {
	return &Cities{
		cities: make(map[string]string),
	}
}

type bagCitiesRoot struct {
	StandBestand bagCitiesBestand `xml:"standBestand"`
}

type bagCitiesBestand struct {
	List    []citiesStand `xml:"stand"`
	DataSet string        `xml:"dataset"`
}

type citiesStand struct {
	Identifier string `xml:"bagObject>Woonplaats>identificatie"`
	City       string `xml:"bagObject>Woonplaats>naam"`
}

func (c *Cities) Parse(r io.Reader) error {
	bagObjects := bagCitiesRoot{}
	err := xml.NewDecoder(r).Decode(&bagObjects)
	if err != nil {
		return err
	}
	for _, addr := range bagObjects.StandBestand.List {
		c.cities[addr.Identifier] = addr.City
	}
	return nil
}

func (c *Cities) Merge(cities []*Cities) {
	for _, source := range cities {
		maps.Copy(c.cities, source.cities)
	}
}

func (c *Cities) Debug() {
	for key, city := range c.cities {
		fmt.Println(key, "=", city)
	}
}

func (c *Cities) Find(id string) *string {
	val, ok := c.cities[id]
	if ok {
		return &val
	} else {
		return nil
	}
}
