package iplist

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/xrstf/mkipset/pkg/ip"
)

type jsonEntry struct {
	IP     string `json:"ip"`
	After  string `json:"after"`
	Before string `json:"before"`
}

func LoadJSONFile(filename string, logger logrus.FieldLogger) (Entries, error) {
	entries := make(Entries, 0)

	f, err := os.Open(filename)
	if err != nil {
		return entries, fmt.Errorf("failed to open: %v", err)
	}
	defer f.Close()

	rawEntries := make([]jsonEntry, 0)
	if err := json.NewDecoder(f).Decode(&rawEntries); err != nil {
		return entries, fmt.Errorf("failed to parse JSON: %v", err)
	}

	for i, rawEntry := range rawEntries {
		entry := Entry{}
		logger = logger.WithField("line", i+1)

		ip, err := ip.Parse(rawEntry.IP)
		if err != nil {
			logger.Warnf("Entry is invalid: %v.", err)
			continue
		}
		entry.IP = ip

		if len(rawEntry.After) > 0 {
			t := parseTime(rawEntry.After)
			if t == nil {
				logger.Warnln("Entry is invalid: invalid `after` date.")
			}

			entry.After = t
		}

		if len(rawEntry.Before) > 0 {
			t := parseTime(rawEntry.Before)
			if t == nil {
				logger.Warnln("Entry is invalid: invalid `before` date.")
			}

			entry.Before = t
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
