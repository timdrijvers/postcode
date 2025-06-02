package bagparser

/*
 * Streets offers parsing capabilities for OPR (OpenbareRuimte)BAG files
 */
import (
	"encoding/xml"
	"fmt"
	"io"
	"maps"
)

type bagOprRoot struct {
	StandBestand bagOprBestand `xml:"standBestand"`
}

type bagOprBestand struct {
	List    []Street `xml:"stand"`
	DataSet string   `xml:"dataset"`
}

type Street struct {
	Identifier string `xml:"bagObject>OpenbareRuimte>identificatie"`
	Name       string `xml:"bagObject>OpenbareRuimte>naam"`
	CityRef    string `xml:"bagObject>OpenbareRuimte>ligtIn>WoonplaatsRef"`
}

type Streets struct {
	streets map[string]Street
}

func NewStreetsParser() *Streets {
	return &Streets{
		streets: make(map[string]Street),
	}
}

func (c *Streets) Parse(r io.Reader) error {
	bagObjects := bagOprRoot{}
	err := xml.NewDecoder(r).Decode(&bagObjects)
	if err != nil {
		return err
	}
	for _, addr := range bagObjects.StandBestand.List {
		c.streets[addr.Identifier] = addr
	}
	return nil
}
func (c *Streets) Merge(streets []*Streets) {
	for _, source := range streets {
		maps.Copy(c.streets, source.streets)
	}
}

func (c *Streets) Debug() {
	for key, addr := range c.streets {
		fmt.Printf("%d = %s [ref=%s]\n", key, addr.Name, addr.CityRef)
	}
}

func (c *Streets) Find(id string) *Street {
	val, ok := c.streets[id]
	if ok {
		return &val
	} else {
		return nil
	}
}
