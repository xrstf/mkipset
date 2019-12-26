package blacklist

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xrstf/mkipset/pkg/ip"
)

var (
	includeLineRegex = regexp.MustCompile(`^include (.+)$`)
	textLineRegex    = regexp.MustCompile(`^([^ ]+)(\s+(.+))?$`)
	afterRegex       = regexp.MustCompile(`after|begins?|since|from|starts?`)
	beforeRegex      = regexp.MustCompile(`before|until|to|ends?`)
	specRegex        = regexp.MustCompile(`([a-z]+)\s+([^ ]+)`)
)

func LoadTextFile(filename string, logger logrus.FieldLogger, allowIncludes bool) (Entries, error) {
	maxIncludes := 0
	if allowIncludes {
		maxIncludes = 100 // meticulously chosen by Skandinavian virgins
	}

	return loadTextFileInternal(make(Entries, 0), filename, logger, &maxIncludes)
}

func loadTextFileInternal(entries Entries, filename string, logger logrus.FieldLogger, remainingIncludes *int) (Entries, error) {
	if *remainingIncludes < 0 {
		return entries, errors.New("include directives were nested too deep or ended up in a loop, giving up")
	}

	absFilename, err := filepath.Abs(filename)
	if err != nil {
		return entries, fmt.Errorf("cannot determine absolute file path for %s: %v", absFilename, err)
	}

	f, err := os.Open(filename)
	if err != nil {
		return entries, fmt.Errorf("failed to open: %v", err)
	}
	defer f.Close()

	lineNo := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNo++
		logger = logger.WithField("line", lineNo)

		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		match := includeLineRegex.FindStringSubmatch(line)
		if match != nil {
			*remainingIncludes--

			rel, err := filepath.Abs(filepath.Join(filepath.Dir(absFilename), match[1]))
			if err != nil {
				return entries, fmt.Errorf("failed to construct file path for include directive '%s': %v", match[1], err)
			}

			entries, err = loadFileInternal(entries, rel, logger, remainingIncludes)
			if err != nil {
				return entries, fmt.Errorf("failed to include %s: %v", rel, err)
			}

			continue
		}

		match = textLineRegex.FindStringSubmatch(line)
		if match == nil {
			logger.Warnln("Line is invalid, no IP/CIDR found.")
			continue
		}

		entry, err := parseTextEntry(match)
		if err != nil {
			logger.Warnf("Line is invalid: %v.", err)
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return entries, fmt.Errorf("failed to read file: %v", err)
	}

	return entries, nil
}

func parseTextEntry(match []string) (Entry, error) {
	e := Entry{}

	ip, err := ip.Parse(match[1])
	if err != nil {
		return e, err
	}
	e.IP = ip

	spec := specRegex.FindAllStringSubmatch(match[3], -1)

	for _, s := range spec {
		verb := strings.ToLower(s[1])
		date := s[2]

		if afterRegex.MatchString(verb) {
			if e.After != nil {
				return e, errors.New("duplicate `after` verb")
			}

			t := parseTime(date)
			if t == nil {
				return e, fmt.Errorf("invalid `after` date '%s'", date)
			}

			e.After = t
		} else if beforeRegex.MatchString(verb) {
			if e.Before != nil {
				return e, errors.New("duplicate `before` verb")
			}

			t := parseTime(date)
			if t == nil {
				return e, fmt.Errorf("invalid `before` date '%s'", date)
			}

			e.Before = t
		} else {
			return e, fmt.Errorf("unknown verb '%s'", verb)
		}
	}

	return e, nil
}
