package bagparser

/*
 * Numbers offers parsing capabilities for NUM (NummerAanduiding) BAG files
 */
import (
	"encoding/xml"
	"fmt"
	"io"
)

type Numbers struct {
	numbers []Property
}

func NewNumbersParser() *Numbers {
	return &Numbers{
		numbers: make([]Property, 0),
	}
}

type bagNumRoot struct {
	StandBestand bagNumBestand `xml:"standBestand"`
}

type bagNumBestand struct {
	List    []Property `xml:"stand"`
	DataSet string     `xml:"dataset"`
}

type Property struct {
	HouseNumber string `xml:"bagObject>Nummeraanduiding>huisnummer"`
	PostalCode  string `xml:"bagObject>Nummeraanduiding>postcode"`
	StreetRef   string `xml:"bagObject>Nummeraanduiding>ligtAan>OpenbareRuimteRef"`
}

func (c *Numbers) Parse(r io.Reader) error {
	bagObjects := bagNumRoot{}
	err := xml.NewDecoder(r).Decode(&bagObjects)
	if err != nil {
		return err
	}

	c.numbers = append(c.numbers, bagObjects.StandBestand.List...)
	return nil
}

func (c *Numbers) Merge(numbers []*Numbers) {
	var total int = len(c.numbers)

	for _, n := range numbers {
		total += len(n.numbers)
	}
	output := make([]Property, total)

	// Copy self first
	offset := copy(output, c.numbers)
	for _, s := range numbers {
		offset += copy(output[offset:], s.numbers)
	}

	c.numbers = output
}

func (c *Numbers) Debug() {
	for _, num := range c.numbers {
		fmt.Printf("housnumber: %s postalcode: %s street-ref: %d\n", num.HouseNumber, num.PostalCode, num.StreetRef)
	}
}

func (c *Numbers) Iterate(streets *Streets, cities *Cities, callback func(postalcode, housenumber, street, city string)) error {
	for _, property := range c.numbers {
		street := streets.Find(property.StreetRef)
		if street == nil {
			return fmt.Errorf(
				"can not find street-ref: %d for property: %s %s",
				property.StreetRef, property.PostalCode, property.HouseNumber,
			)
		}
		city := cities.Find(street.CityRef)
		if city == nil {
			return fmt.Errorf(
				"can not find city-ref: %d for property: %s %s",
				street.CityRef, property.PostalCode, property.HouseNumber,
			)
		}
		callback(property.PostalCode, property.HouseNumber, street.Name, *city)
	}
	return nil
}
