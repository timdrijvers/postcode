package bagparser

import (
	"archive/zip"
	"fmt"
	"hash/fnv"
	"io"
	"path/filepath"
	"strings"
)

type ZipParser interface {
	Parse(r io.Reader) error
}

func patternFromFilename(filename string) string {
	return strings.Replace(filepath.Base(filename), ".zip", "", 1)
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func ParseZip(filename string, parser ZipParser) error {
	return ParseZipSharded(filename, parser, 0, 1)
}

func ParseZipSharded(filename string, parser ZipParser, shard uint32, totalShards uint32) error {
	zf, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer zf.Close()
	pattern := patternFromFilename(filename)

	for _, file := range zf.File {
		if strings.HasPrefix(file.Name, pattern) && strings.HasSuffix(file.Name, ".xml") {
			if hash(file.Name)%totalShards != shard {
				continue
			}
			fmt.Printf("ParseZip[%s][%d/%d] - %s\n", pattern, shard, totalShards, file.Name)
			dfp, derr := file.Open()
			if derr != nil {
				return derr
			}
			parseError := parser.Parse(dfp)
			dfp.Close()
			if parseError != nil {
				return parseError
			}
		}

	}
	return nil
}
