package blacklist

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func LoadFile(filename string, logger logrus.FieldLogger) (Entries, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".yaml":
		fallthrough
	case ".yml":
		return LoadYAMLFile(filename, logger)

	case ".json":
		return LoadJSONFile(filename, logger)

	case ".txt":
		return LoadTextFile(filename, logger)

	default:
		return nil, fmt.Errorf("unknown file extension '%s'", ext)
	}
}

var timeFormats = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04",
	"2006-01-02",
	"2006-01",
	"2006",
}

func parseTime(d string) *time.Time {
	for _, format := range timeFormats {
		parsed, err := time.Parse(format, d)
		if err == nil {
			return &parsed
		}
	}

	return nil
}
