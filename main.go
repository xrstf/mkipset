package main

import (
	"flag"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xrstf/mkipset/pkg/config"
	"github.com/xrstf/mkipset/pkg/iplist"
	"github.com/xrstf/mkipset/pkg/ipset"
)

func main() {
	var (
		configFile string
		verbose    bool
		pretty     bool
	)

	flag.StringVar(&configFile, "config", "", "configuration file to use")
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

	listFile := flag.Arg(0)
	if len(listFile) == 0 {
		logger.Fatalln("No LIST_FILE argument given.")
	}

	logger.Debugf("Loading IP list %s…", listFile)
	list, err := iplist.LoadFile(listFile, logger)
	if err != nil {
		logger.Fatalf("Failed to load IP list: %v.", err)
	}

	active := list.Active(time.Now())
	logger.Debugf("List contains %d total entries, %d of which are active.", len(list), len(active))

	filtered := list.RemoveCollisions(config.WhitelistIPs())
	logger.Debugf("List contains %d entries after removing whitelisted entries.", len(filtered))

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

	logger.Debugln("Synchronizing ipset…")
	err = client.Synchronize(*set, filtered.IPs())
	if err != nil {
		logger.Fatalf("Failed to synchronize ipset set: %v.", err)
	}
	logger.Debugln("Set has been synchronized.")
}
