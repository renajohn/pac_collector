package collector

import (
	"log"

	"github.com/renajohn/pac_collector/api"
)

// Collector bind a source to a store
type Collector struct {
	Source        api.Source
	Sink          api.Sink
	sourceChannel <-chan api.Measurement
}

// Start initiate the collection process
func (c *Collector) Start() {
	var err error
	c.sourceChannel, err = c.Source.Start()

	if err != nil {
		log.Fatal(err)
	}

	c.collect()
}

func (c *Collector) collect() {
	for measure := range c.sourceChannel {
		c.Sink.Put(measure)
	}
}
