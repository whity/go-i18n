package storage

import (
	"fmt"

	"go.yaml.in/yaml/v3"
)

type Directory interface {
	ReadFile(name string) ([]byte, error)
}

type yamlStorage struct {
	dir Directory
}

func (s *yamlStorage) Load(locale string) (Data, error) {
	contents, err := s.dir.ReadFile(locale + ".yml")
	if err != nil {
		return nil, fmt.Errorf("failed to read locale file: %v", err)
	}

	var data Data
	err = yaml.Unmarshal(contents, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to read locale file: %v", err)
	}

	return data, nil
}

func NewYaml(dir Directory) Storage {
	return &yamlStorage{dir: dir}
}
