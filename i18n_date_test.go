package i18n

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateNoFormats(t *testing.T) {
	yaml := []string{
		"---",
	}

	i, _ := newI18nFromYAMLString(yaml)
	str, err := i.Date(time.Now())

	assert.EqualValues(t, "", str)
	assert.EqualError(
		t,
		err,
		"failed to format date: could not read formats definition",
	)
}

func TestDateFormatsWrongType(t *testing.T) {
	yamls := [][]string{
		[]string{
			"---",
			"formats:",
			"  date: 1",
		},
		[]string{
			"---",
			"formats:",
			"  date:",
			"    formats: 3",
		},
	}

	for _, yaml := range yamls {
		i, _ := newI18nFromYAMLString(yaml)
		str, err := i.Date(time.Now())

		assert.EqualValues(t, "", str)
		assert.EqualError(
			t,
			err,
			"failed to format date: could not read 'date' formats",
		)
	}
}

func TestDateInvalidDisplayFormat(t *testing.T) {
	yaml := []string{
		"---",
		"formats:",
		"  date:",
		"    formats: ",
		"      def: ",
	}

	i, _ := newI18nFromYAMLString(yaml)
	str, err := i.Date(time.Now())

	assert.EqualValues(t, "", str)
	assert.EqualError(
		t,
		err,
		"failed to format date: failed to read format 'default'",
	)
}

func TestDateNamesPlaceholders(t *testing.T) {
	yaml := []string{
		"---",
		"formats:",
		"  date:",
		"    abbr_month_names:",
		"      - jan",
		"    month_names:",
		"      - january",
		"    abbr_day_names: ['sun', 'mon', 'tue', 'wed', 'thu', 'fri', 'sat']",
		"    day_names:",
		"      - sunday",
		"      - monday",
		"      - tuesday",
		"      - wednesday",
		"      - thursday",
		"      - friday",
		"      - saturday",
		"    formats: ",
		"      default: '%a %A %b %B'",
	}

	i, err := newI18nFromYAMLString(yaml)
	fmt.Println(err)
	dt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	str, _ := i.Date(dt)

	assert.EqualValues(t, "thu thursday jan january", str)
}

func TestDateNamesPlaceholdersMissingMappers(t *testing.T) {
	yaml := []string{
		"---",
		"formats:",
		"  date:",
		"    abbr_month_names:",
		"    month_names:",
		"    abbr_day_names: ['sun', 'mon', 'tue', 'wed']",
		"    day_names:",
		"      - sunday",
		"      - monday",
		"      - tuesday",
		"      - wednesday",
		"    formats: ",
		"      default: '%a %A %b %B'",
	}

	i, err := newI18nFromYAMLString(yaml)
	fmt.Println(err)
	dt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	str, _ := i.Date(dt)

	assert.EqualValues(t, "%a %A %b %B", str)
}

func TestDateISOFormat(t *testing.T) {
	yaml := []string{
		"---",
		"formats:",
		"  date:",
		"    formats: ",
		"      default: '%Y-%m-%d'",
	}

	i, err := newI18nFromYAMLString(yaml)
	fmt.Println(err)
	dt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	str, _ := i.Date(dt)

	assert.EqualValues(t, "2026-01-01", str)
}

func TestDateSpecifyFormat(t *testing.T) {
	yaml := []string{
		"---",
		"formats:",
		"  date:",
		"    formats: ",
		"      default: '%Y-%m-%d'",
		"      no-day: '%Y-%m'",
	}

	i, err := newI18nFromYAMLString(yaml)
	fmt.Println(err)
	dt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	str, _ := i.Date(dt, WithFormat("no-day"))

	assert.EqualValues(t, "2026-01", str)
}
