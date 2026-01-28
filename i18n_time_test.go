package i18n

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeNoFormats(t *testing.T) {
	yaml := []string{
		"---",
	}

	i, _ := newI18nFromYAMLString(yaml)
	str, err := i.Time(time.Now())

	assert.EqualValues(t, "", str)
	assert.EqualError(
		t,
		err,
		"failed to format time: could not read formats definition",
	)
}

func TestTimeFormatsWrongType(t *testing.T) {
	yamls := [][]string{
		[]string{
			"---",
			"formats:",
			"  time: 1",
		},
		[]string{
			"---",
			"formats:",
			"  time:",
			"    formats: 3",
		},
	}

	for _, yaml := range yamls {
		i, _ := newI18nFromYAMLString(yaml)
		str, err := i.Time(time.Now())

		assert.EqualValues(t, "", str)
		assert.EqualError(
			t,
			err,
			"failed to format time: could not read 'time' formats",
		)
	}
}

func TestTimeInvalidDisplayFormat(t *testing.T) {
	yaml := []string{
		"---",
		"formats:",
		"  time:",
		"    formats: ",
		"      def: ",
	}

	i, _ := newI18nFromYAMLString(yaml)
	str, err := i.Time(time.Now())

	assert.EqualValues(t, "", str)
	assert.EqualError(
		t,
		err,
		"failed to format time: failed to read format 'default'",
	)
}

func TestTimeISOFormat(t *testing.T) {
	yaml := []string{
		"---",
		"formats:",
		"  time:",
		"    formats: ",
		"      default: '%H:%M:%S'",
	}

	i, _ := newI18nFromYAMLString(yaml)
	dt := time.Date(2026, 1, 1, 1, 0, 0, 0, time.UTC)
	str, _ := i.Time(dt)

	assert.EqualValues(t, "01:00:00", str)
}

func TestTimeSpecifyFormat(t *testing.T) {
	yaml := []string{
		"---",
		"formats:",
		"  time:",
		"    formats: ",
		"      default: '%H:%M:%S'",
		"      no-seconds: '%H:%M'",
	}

	i, _ := newI18nFromYAMLString(yaml)
	dt := time.Date(2026, 1, 1, 1, 0, 0, 0, time.UTC)
	str, _ := i.Time(dt, WithFormat("no-seconds"))

	assert.EqualValues(t, "01:00", str)
}
