package main

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/renajohn/pac_collector/collector"
	"github.com/renajohn/pac_collector/internal/kafkasink"
	"github.com/renajohn/pac_collector/internal/swcsource"
)

type pacMonConfig struct {
	sourceURL       string
	sinkURL         string
	kafkaTopic      string
	pollingInterval time.Duration
}

func parseCmdParams(args []string) (*pacMonConfig, error) {

	var commandLine = flag.NewFlagSet(args[0], flag.ExitOnError)
	sourceURLPtr := commandLine.String("sourceURL", "", "Source end point URL")
	sinkURLPtr := commandLine.String("sinkURL", "", "Sink end point URL. This URL should point to a kafka broker.")
	topicPtr := commandLine.String("topic", "SWCTemperature", "[Optional] Which kafka topic to use")
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
		kafkaTopic:      *topicPtr,
	}

	return &config, nil
}

func main() {
	config, err := parseCmdParams(os.Args)

	if err != nil {
		os.Exit(1)
	}

	sink := kafkasink.NewKafkaSink(config.sinkURL, config.kafkaTopic)
	source := swcsource.NewSWCSource(config.sourceURL, config.pollingInterval)

	collector := collector.Collector{
		Source: source,
		Sink:   sink,
	}

	collector.Start()
}
