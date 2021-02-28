package main

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/renajohn/pac_collector/collector"
	"github.com/renajohn/pac_collector/internal/mocksink"
	"github.com/renajohn/pac_collector/internal/swcsource"
)

type pacMonConfig struct {
	sourceURL       string
	sinkURL         string
	pollingInterval time.Duration
}

func parseCmdParams(args []string) (*pacMonConfig, error) {

	var commandLine = flag.NewFlagSet(args[0], flag.ExitOnError)
	sourceURLPtr := commandLine.String("sourceURL", "", "Source end point URL")
	sinkURLPtr := commandLine.String("sinkURL", "", "Sink end point URL")
	intervalPtr := commandLine.Int("pollingInterval", 60, "[Optional] Interval in seconds at which data will be fetched (min 1s)")

	commandLine.Parse(args[1:])

	if len(*sourceURLPtr) == 0 || len(*sinkURLPtr) == 0 || *intervalPtr < 1 {
		commandLine.Usage()
		return nil, errors.New("incorrect parameters")
	}

	config := pacMonConfig{
		sourceURL:       *sourceURLPtr,
		sinkURL:         *sinkURLPtr,
		pollingInterval: time.Duration(*intervalPtr) * time.Second,
	}

	return &config, nil
}

func main() {
	config, err := parseCmdParams(os.Args)

	if err != nil {
		os.Exit(1)
	}

	sink := mocksink.MockSink{}
	//source := swcsource.NewSWCSource("ws://192.168.086.29:8214/", 2*time.Second)
	source := swcsource.NewSWCSource(config.sourceURL, config.pollingInterval)

	collector := collector.Collector{
		Source: source,
		Sink:   &sink,
	}

	collector.Start()
}
