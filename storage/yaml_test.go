package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type memDirectory struct {
	data map[string]string
}

func (d *memDirectory) ReadFile(name string) ([]byte, error) {
	contents, ok := d.data[name]
	if !ok {
		return []byte{}, fmt.Errorf("file '%s' not found", name)
	}

	return []byte(contents), nil
}

func TestStorageNewYaml(t *testing.T) {
	mDir := &memDirectory{}
	expected := yamlStorage{mDir}
	store := NewYaml(mDir)

	assert.EqualValues(t, &expected, store)
}

func TestStorageYamlLoadUnknownLocale(t *testing.T) {
	mDir := &memDirectory{}
	store := NewYaml(mDir)
	locale, err := store.Load("en")

	assert.Nil(t, locale)
	assert.EqualError(t, err, "failed to read locale file: file 'en.yml' not found")
}

func TestStorageYamlLoadFailedToUnmarshall(t *testing.T) {
	mDir := &memDirectory{
		map[string]string{
			"en.yml": `---
		translations:
			hello: world`,
		},
	}
	store := NewYaml(mDir)
	locale, err := store.Load("en")

	assert.Nil(t, locale)
	assert.EqualError(
		t,
		err,
		"failed to read locale file: yaml: line 2: found character that cannot start any token",
	)
}

func TestStorageYamlLoadSuccess(t *testing.T) {
	mDir := &memDirectory{
		map[string]string{
			"en.yml": `---
translations:
  hello: world`,
		},
	}
	store := NewYaml(mDir)
	locale, err := store.Load("en")

	assert.Nil(t, err)
	assert.EqualValues(
		t,
		Data{
			"translations": Data{
				"hello": "world",
			},
		},
		locale,
	)
}
