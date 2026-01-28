/*
Package i18n implements a simple approach to deal with translations and date/time localization.

Example:

	package main

	import (
		"fmt"
		"time"

		"github.com/whity/go-i18n"
		"github.com/whity/go-i18n/storage"
	)

	func main() {
		// Reading from yaml embedded files
		// yamlStore := storage.NewYaml(*embed.FS)

		// Map store
		store := storage.NewMap(storage.Data{
			"en": storage.Data{
				"translations": storage.Data{
					"hello": "world",
					"new_message": storage.Data{
						"1":     "you have a new message",
						"2..":   "you have {{.count}} new messages",
						"other": "no message",
					},
				},
				"formats": storage.Data{
					"time": storage.Data{
						"formats": storage.Data{
							"default": "%I:%M:%S %p",
						},
					},
					"date": storage.Data{
						"abbr_month_names": []any{
							"jan",
							"feb",
							"mar",
							"apr",
							"may",
							"jun",
							"jul",
							"aug",
							"sep",
							"oct",
							"nov",
							"dec",
						},
						"month_names": []any{
							"january",
							"february",
							"march",
							"april",
							"may",
							"june",
							"july",
							"august",
							"september",
							"october",
							"november",
							"december",
						},
						"abbr_day_names": []any{
							"sun",
							"mon",
							"tue",
							"wed",
							"thu",
							"fri",
							"sat",
						},
						"day_names": []any{
							"sunday",
							"monday",
							"tuesday",
							"wednesday",
							"thursday",
							"friday",
							"saturday",
						},
						"formats": storage.Data{
							"default": "%Y-%m-%d",
							"long":    "%d of %b %Y",
						},
					},
				},
			},
			"pt": storage.Data{
				"translations": storage.Data{
					"hello": "mundo",
					"new_message": storage.Data{
						"1":     "tem uma nova mensagem",
						"2..":   "tem {{.count}} novas mensagens",
						"other": "sem mensagens",
					},
				},
				"formats": storage.Data{
					"time": storage.Data{
						"formats": storage.Data{
							"default": "%H:%M:%S",
						},
					},
				},
			},
		})

		i, err := i18n.New(store, "en")
		if err != nil {
			panic(err)
		}

		fmt.Println(i.Translate("hello"))
		fmt.Println(i.Translate("new_message", i18n.TranslateWithCount(1)))
		fmt.Println(i.Translate("new_message", i18n.TranslateWithCount(0)))
		fmt.Println(i.Translate("new_message", i18n.TranslateWithCount(4)))
		fmt.Println(i.Date(time.Now(), i18n.WithFormat("long")))
		fmt.Println(i.Time(time.Now()))
	}
*/
package i18n

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/lestrrat-go/strftime"
	"github.com/whity/go-i18n/storage"
)

var (
	isNumberRe        *regexp.Regexp = regexp.MustCompile(`^\d+$`)
	isNumberRangeRe   *regexp.Regexp = regexp.MustCompile(`^(\d+)..(\d+)?$`)
	datePlaceholderRe *regexp.Regexp = regexp.MustCompile(`%[aAbB]`)
)

type i18n struct {
	st   storage.Storage
	data storage.Data
}

func New(st storage.Storage, locale string) (*i18n, error) {
	data, err := st.Load(locale)
	if err != nil {
		return nil, err
	}

	return &i18n{st: st, data: data}, nil
}

func (i *i18n) Translate(key string, opts ...translateOption) (string, error) {
	translations, ok := i.data["translations"].(storage.Data)
	if !ok {
		return key, nil
	}

	// If a translation is not present, return the key
	result, ok := translations[key]
	if !ok {
		return key, nil
	}

	cfg := translateConfig{}
	cfg.applyOptions(opts...)

	switch result := result.(type) {
	case string:
		return result, nil
	case storage.Data:
		count := cfg.count

		if count <= 0 {
			if other, ok := result["other"]; ok {
				return convertTranslation2String(other), nil
			}

			return key, nil
		}

		var msg any

		for k, v := range result {
			if k == "other" {
				continue
			}

			// Check if key is a number.
			match := isNumberRe.MatchString(k)
			if match && k == fmt.Sprintf("%d", count) {
				msg = v
				break
			}

			// Check if key is a range, otherwise continue to the next key
			matches := isNumberRangeRe.FindStringSubmatch(k)
			if len(matches) == 0 {
				continue
			}

			start, _ := strconv.Atoi(matches[1])
			endStr := strings.TrimSpace(matches[2])

			// If end isn't defined, just check if count is bigger
			//	or equal than start
			if endStr == "" {
				if count >= start {
					msg = v
					break
				}

				continue
			}

			// Check if count is between the range start and end
			end, _ := strconv.Atoi(endStr)

			if count >= start && count <= end {
				msg = v
				break
			}
		}

		if msg != nil {
			return processTranslation(
				convertTranslation2String(msg),
				count,
			)
		}
	}

	return key, nil
}

func (i *i18n) Date(t time.Time, opts ...localizeOption) (string, error) {
	dateFormats, err := i.getFormats("date")
	if err != nil {
		return "", fmt.Errorf("failed to format date: %v", err)
	}

	formats, ok := dateFormats["formats"].(storage.Data)
	if !ok {
		return "", fmt.Errorf("failed to format date: could not read 'date' formats")
	}

	cfg := newLocalizeConfig()
	cfg.applyOptions(opts...)

	format, ok := formats[cfg.format].(string)
	if !ok {
		return "", fmt.Errorf(
			"failed to format date: failed to read format '%s'",
			cfg.format,
		)
	}

	namesMap := map[string]string{
		"%a": "abbr_day_names",
		"%A": "day_names",
		"%b": "abbr_month_names",
		"%B": "month_names",
	}

	result := datePlaceholderRe.ReplaceAllStringFunc(
		format,
		func(match string) string {
			names, ok := dateFormats[namesMap[match]].([]any)
			if !ok {
				return "%" + match
			}

			var index int

			if match == "%a" || match == "%A" {
				index = int(t.Weekday())
			} else {
				index = int(t.Month()) - 1
			}

			if index >= len(names) {
				return "%" + match
			}

			return fmt.Sprintf("%v", names[index])
		},
	)

	result, err = strftime.Format(result, t)
	if err != nil {
		return "", fmt.Errorf("failed to format date: %v", err)
	}

	return result, nil
}

func (i *i18n) Time(t time.Time, opts ...localizeOption) (string, error) {
	timeFormats, err := i.getFormats("time")
	if err != nil {
		return "", fmt.Errorf("failed to format time: %v", err)
	}

	formats, ok := timeFormats["formats"].(storage.Data)
	if !ok {
		return "", fmt.Errorf("failed to format time: could not read 'time' formats")
	}

	cfg := newLocalizeConfig()
	cfg.applyOptions(opts...)

	format, ok := formats[cfg.format].(string)
	if !ok {
		return "", fmt.Errorf(
			"failed to format time: failed to read format '%s'",
			cfg.format,
		)
	}

	result, err := strftime.Format(format, t)
	if err != nil {
		return "", fmt.Errorf("failed to format time: %v", err)
	}

	return result, nil
}

func (i *i18n) getFormats(tp string) (storage.Data, error) {
	formats, ok := i.data["formats"].(storage.Data)
	if !ok {
		return nil, fmt.Errorf("could not read formats definition")
	}

	typeFormats, ok := formats[tp].(storage.Data)
	if !ok {
		return nil, fmt.Errorf("could not read '%s' formats", tp)
	}

	return typeFormats, nil
}

func convertTranslation2String(msg any) string {
	return fmt.Sprintf("%v", msg)
}

func processTranslation(msg string, count int) (string, error) {
	tmpl, err := template.New("translation").Parse(msg)
	if err != nil {
		return "", nil
	}

	str := strings.Builder{}
	tmpl.Execute(&str, storage.Data{"count": count})

	return str.String(), nil
}

type translateConfig struct {
	count int
}

type translateOption func(c *translateConfig)

func (c *translateConfig) applyOptions(opts ...translateOption) {
	for _, opt := range opts {
		opt(c)
	}
}

func TranslateWithCount(count int) translateOption {
	return func(c *translateConfig) {
		c.count = count
	}
}

type localizeConfig struct {
	format string
}

func newLocalizeConfig() *localizeConfig {
	return &localizeConfig{format: "default"}
}

type localizeOption func(c *localizeConfig)

func WithFormat(name string) localizeOption {
	return func(c *localizeConfig) {
		c.format = name
	}
}

func (c *localizeConfig) applyOptions(opts ...localizeOption) {
	for _, opt := range opts {
		opt(c)
	}
}
