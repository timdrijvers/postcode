package bagparser

import (
	"sync"
)

func ParseSharded[K ZipParser](filename string, totalShards uint, constructor func() K, merger func([]K)) error {
	parsers := make([]K, totalShards)
	errors := make(chan error, totalShards)

	waitGroup := &sync.WaitGroup{}
	for shardIndex := 0; shardIndex < len(parsers); shardIndex++ {
		parsers[shardIndex] = constructor()
		waitGroup.Add(1)

		go func(err chan<- error, parser K, shard uint32) {
			defer waitGroup.Done()
			err <- ParseZipSharded(filename, parser, shard, uint32(totalShards))
		}(errors, parsers[shardIndex], uint32(shardIndex))
	}
	waitGroup.Wait()
	close(errors)

	// Return first error, if any
	for err := range errors {
		if err != nil {
			return err
		}
	}

	merger(parsers)
	return nil
}
