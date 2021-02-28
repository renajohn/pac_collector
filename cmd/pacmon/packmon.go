package main

import (
	"time"

	"github.com/renajohn/pac_collector/collector"
	"github.com/renajohn/pac_collector/internal/mocksink"
	"github.com/renajohn/pac_collector/internal/swcsource"
)

func main() {
	sink := mocksink.MockSink{}
	source := swcsource.NewSWCSource("ws://192.168.086.29:8214/", 2*time.Second)

	collector := collector.Collector{
		Source: source,
		Sink:   &sink,
	}

	collector.Start()
}
