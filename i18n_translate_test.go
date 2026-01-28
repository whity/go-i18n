package i18n

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/whity/go-i18n/storage"
)

type tmpDir struct {
	readFile func(string) ([]byte, error)
}

func (d *tmpDir) ReadFile(name string) ([]byte, error) {
	return d.readFile(name)
}

func newI18nFromYAMLString(str []string) (*i18n, error) {
	dir := &tmpDir{
		readFile: func(name string) ([]byte, error) {
			return []byte(strings.Join(str, "\n")), nil
		},
	}

	store := storage.NewYaml(dir)
	return New(store, "en")
}

func TestNewError(t *testing.T) {
	dir := &tmpDir{
		readFile: func(name string) ([]byte, error) {
			return nil, fmt.Errorf("file '%s' not found", name)
		},
	}

	store := storage.NewYaml(dir)
	i, err := New(store, "en")

	assert.Nil(t, i)
	assert.EqualError(
		t,
		err,
		"failed to read locale file: file 'en.yml' not found",
	)
}

func TestTranslateNoTranslationsKey(t *testing.T) {
	i, _ := newI18nFromYAMLString([]string{"---"})
	translation, _ := i.Translate("hello")

	assert.EqualValues(t, "hello", translation)
}

func TestTranslateNoKey(t *testing.T) {
	yaml := []string{
		"---",
		"translations: {}",
	}

	i, _ := newI18nFromYAMLString(yaml)
	translation, _ := i.Translate("hello")

	assert.EqualValues(t, "hello", translation)
}

func TestTranslateKeyExists(t *testing.T) {
	yaml := []string{
		"---",
		"translations:",
		"  hello: world",
	}

	i, _ := newI18nFromYAMLString(yaml)
	translation, _ := i.Translate("hello")

	assert.EqualValues(t, "world", translation)
}

func TestTranslateKeyWithCount(t *testing.T) {
	yaml := []string{
		"---",
		"translations:",
		"  new_message:",
		"    '1': you have a new message",
		"    '2..': you have {{.count}} new messages",
		"    'other': no messages",
	}

	i, _ := newI18nFromYAMLString(yaml)

	tWithoutCount, _ := i.Translate("new_message")
	assert.EqualValues(t, "no messages", tWithoutCount, "no count passed")

	tWithCountOne, _ := i.Translate("new_message", TranslateWithCount(1))
	assert.EqualValues(
		t,
		"you have a new message",
		tWithCountOne,
		"count passed (1)",
	)

	tWithCountThree, _ := i.Translate("new_message", TranslateWithCount(3))
	assert.EqualValues(
		t,
		"you have 3 new messages",
		tWithCountThree,
		"count passed (3)",
	)
}

func TestTranslateKeyWithCountRange(t *testing.T) {
	yaml := []string{
		"---",
		"translations:",
		"  new_message:",
		"    '1..2': you have {{.count}} new messages",
	}

	i, _ := newI18nFromYAMLString(yaml)

	insideRange, _ := i.Translate("new_message", TranslateWithCount(2))
	assert.EqualValues(t, "you have 2 new messages", insideRange, "inside range")

	outOfRange, _ := i.Translate("new_message", TranslateWithCount(3))
	assert.EqualValues(
		t,
		"new_message",
		outOfRange,
		"outside range",
	)
}

func TestTranslateKeyWithCountZeroNoOther(t *testing.T) {
	yaml := []string{
		"---",
		"translations:",
		"  new_message:",
		"    '1..2': you have {{.count}} new messages",
	}

	i, _ := newI18nFromYAMLString(yaml)

	countZero, _ := i.Translate("new_message")
	assert.EqualValues(t, "new_message", countZero)
}
