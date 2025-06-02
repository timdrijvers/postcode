package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/timdrijvers/postcode/bagparser"
	"github.com/timdrijvers/postcode/storage"
)

func main() {
	shards := flag.Uint("shards", uint(runtime.NumCPU()), "amount of shards, defaults to amount of cores")
	flag.Usage = func() {
		fmt.Println("usage: builder WPL.zip OPR.zip NUM.zip output.db")
		fmt.Println("Builds a processed database for the server, takes a lvbag extract from Kadaster.")
		fmt.Println("")
		fmt.Println("  WPL.zip   9999WPL*.zip file from the extracted lvlbag file")
		fmt.Println("  OPR.zip   9999OPR*.zip file from the extracted lvlbag file")
		fmt.Println("  NUM.zip   9999NUM*.zip file from the extracted lvlbag file")
		fmt.Println("  output.db generated output database, input for server")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  --shards N  set number of parallel processing")
		fmt.Println("")
	}
	flag.Parse()

	if flag.NArg() != 4 {
		flag.Usage()
		return
	}

	cities := bagparser.NewCitiesParser()
	if err := bagparser.ParseZip(flag.Arg(0), cities); err != nil {
		panic(err)
	}

	streets := bagparser.NewStreetsParser()
	if err := bagparser.ParseSharded(flag.Arg(1), *shards, bagparser.NewStreetsParser, streets.Merge); err != nil {
		panic(err)
	}

	numbers := bagparser.NewNumbersParser()
	if err := bagparser.ParseSharded(flag.Arg(2), *shards, bagparser.NewNumbersParser, numbers.Merge); err != nil {
		panic(err)
	}

	output, err := os.Create(flag.Arg(3))
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	db := make(storage.Addresses)
	if err = numbers.Iterate(streets, cities, db.Add); err != nil {
		panic(err)
	}
	if err = db.Write(output); err != nil {
		panic(err)
	}
	if err = output.Close(); err != nil {
		panic(err)
	}
}
