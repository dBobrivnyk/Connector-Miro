package idstorage

import (
	"encoding/json"
	"io"
	"os"
)

type Id struct {
	Owner  string `json:"owner"`
	ApiKey string `json:"api_key"`
}

func Load(from func() (io.ReadCloser, error)) ([]Id, error) {
	source, err := from()
	if err != nil {
		return nil, err
	}

	defer source.Close()
	bytes, err := io.ReadAll(source)
	if err != nil {
		return nil, err
	}

	var allIds []Id
	err = json.Unmarshal(bytes, &allIds)
	if err != nil {
		return nil, err
	}

	return allIds, nil
}

func FromFile(path string) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) {
		return os.Open(path)
	}
}
