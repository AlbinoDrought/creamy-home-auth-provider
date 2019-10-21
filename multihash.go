package main

import (
	"errors"

	"github.com/multiformats/go-multihash"
)

type mapper struct {
	multihashCode uint64
}

func (m *mapper) Map(username string) (string, error) {
	hash, err := multihash.Sum([]byte(username), m.multihashCode, -1)
	if err != nil {
		return "", err
	}

	return hash.B58String(), nil
}

// Make a multihash LinkMapper using the given multihash mode
func makeHashMapper(multihashName string) (*mapper, error) {
	code, ok := multihash.Names[multihashName]
	if !ok {
		return nil, errors.New("bad multihash name")
	}

	return &mapper{code}, nil
}

// Must make a multihash LinkMapper with the given mode, or panic on failure
func mustMakeHashMapper(multihashName string) *mapper {
	mapper, err := makeHashMapper(multihashName)
	if err != nil {
		panic(err)
	}
	return mapper
}
