package iplist

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type yamlEntry struct {
	IP     string `yaml:"ip"`
	After  string `yaml:"after"`
	Before string `yaml:"before"`
}

func LoadYAMLFile(filename string, logger logrus.FieldLogger) (Entries, error) {
	entries := make(Entries, 0)

	f, err := os.Open(filename)
	if err != nil {
		return entries, fmt.Errorf("failed to open: %v", err)
	}
	defer f.Close()

	rawEntries := make([]yamlEntry, 0)
	if err := yaml.NewDecoder(f).Decode(&rawEntries); err != nil {
		return entries, fmt.Errorf("failed to parse YAML: %v", err)
	}

	for i, rawEntry := range rawEntries {
		entry := Entry{}
		logger = logger.WithField("line", i+1)

		ip, err := parseIP(rawEntry.IP)
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
