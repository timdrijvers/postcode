package storage

import (
	"encoding/gob"
	"os"
	"strings"
)

type Address struct {
	Street string
	City   string
}

type AddressKey struct {
	PostalCode  string
	HouseNumber string
}

type FullAddress struct {
	Street      string
	City        string
	PostalCode  string
	HouseNumber string
}

func addressKey(postalcode string, housenumber string) AddressKey {
	return AddressKey{PostalCode: normalize(postalcode), HouseNumber: normalize(housenumber)}
}

type Addresses map[AddressKey]*Address

func normalize(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "", -1)
}

// Load data from file into memory
func NewAddressesStorageFromFile(f *os.File) (Addresses, error) {
	decoder := gob.NewDecoder(f)
	addresses := make(Addresses)

	err := decoder.Decode(&addresses)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// Add a new entry to the map
func (a Addresses) Add(postalcode string, housenumber string, street string, city string) {
	a[addressKey(postalcode, housenumber)] = &Address{Street: street, City: city}
}

// Look for an entry in the map
func (a Addresses) Find(postalcode string, housenumber string) *FullAddress {
	key := addressKey(postalcode, housenumber)
	val, ok := a[key]
	if ok {
		return &FullAddress{
			HouseNumber: key.HouseNumber,
			PostalCode:  key.PostalCode,
			Street:      val.Street,
			City:        val.City,
		}
	} else {
		return nil
	}
}

// Write map to file
func (a Addresses) Write(f *os.File) error {
	return gob.NewEncoder(f).Encode(a)
}
