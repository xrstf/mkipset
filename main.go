package main

import (
	"flag"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xrstf/mkipset/pkg/blacklist"
	"github.com/xrstf/mkipset/pkg/config"
	"github.com/xrstf/mkipset/pkg/ipset"
)

func main() {
	var (
		configFile         string
		verbose            bool
		pretty             bool
		ignoreMissingFiles bool
	)

	flag.StringVar(&configFile, "config", "", "configuration file to use")
	flag.BoolVar(&ignoreMissingFiles, "ignore-missing", false, "when given, do not abort on missing blacklist files")
	flag.BoolVar(&verbose, "verbose", false, "enable more verbose logging")
	flag.BoolVar(&pretty, "pretty", false, "enable pretty logging output with colors")
	flag.Parse()

	logger := logrus.New()
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	if pretty {
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.Stamp,
		})
	}

	logger.Debugf("Loading configuration file %s…", configFile)
	config, err := config.LoadFromFile(configFile)
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v.", err)
	}

	if flag.NArg() == 0 {
		logger.Fatalln("No blacklist files given.")
	}

	var entries blacklist.Entries

	files := 0

	for _, file := range flag.Args() {
		flogger := logger.WithField("file", file)
		flogger.Debugln("Loading file…")

		fentries, err := blacklist.LoadFile(file, flogger)
		if err != nil {
			if ignoreMissingFiles {
				flogger.Warnf("Failed to load IP list: %v.", err)
				continue
			}

			flogger.Fatalf("Failed to load IP list: %v.", err)
		}

		entries = entries.Merge(fentries)
		files++
	}

	if files == 0 {
		logger.Fatalln("Could not load any of the given files, refusing to continue.")
	}

	active := entries.Active(time.Now())
	logger.Debugf("Found a total of %d entries, %d of which are active.", len(entries), len(active))

	filtered := active.RemoveCollisions(config.WhitelistIPs())
	logger.Debugf("Of those, %d entries remained after removing whitelisted elements.", len(filtered))

	logger.Debugln("Building ipset client…")
	client, err := ipset.NewExec()
	if err != nil {
		logger.Fatalf("Failed to build ipset client: %v.", err)
	}

	logger.Debugf("Ensuring ipset set %s…", config.SetName)
	err = client.Create(config.SetName, ipset.SetTypeHashNet)
	if err != nil {
		logger.Fatalf("Failed to create ipset set: %v.", err)
	}

	logger.Debugln("Finding existing set entries…")
	set, err := client.Show(config.SetName)
	if err != nil {
		logger.Fatalf("Failed to show ipset set: %v.", err)
	}
	logger.Debugf("Set contains %d members.", len(set.Members))

	ipsToBlacklist := filtered.IPs()

	if set.MembersEquals(ipsToBlacklist) {
		logger.Debugln("Blacklist equals current set members, nothing to do.")
		os.Exit(0)
	}

	logger.Debugln("Synchronizing ipset…")
	err = client.Synchronize(*set, ipsToBlacklist)
	if err != nil {
		logger.Fatalf("Failed to synchronize ipset set: %v.", err)
	}
	logger.Debugln("Set has been synchronized.")
}
