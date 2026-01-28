package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageNewMap(t *testing.T) {
	expected := mapStorage{data: Data{}}
	store := NewMap(Data{})

	assert.EqualValues(t, &expected, store)
}

func TestStorageMapLoadUnknownLocale(t *testing.T) {
	store := NewMap(Data{})
	locale, err := store.Load("en")

	assert.Nil(t, err)
	assert.EqualValues(t, Data{}, locale)
}

func TestStorageMapLoadKnownLocale(t *testing.T) {
	store := NewMap(Data{
		"en": Data{
			"translations": Data{
				"hello": "world",
			},
		},
	})
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
