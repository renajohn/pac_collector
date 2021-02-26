package collector

import (
	"github.com/renajohn/pac_collector/api"
)

// Collector bind a source to a store
type Collector struct {
	Source api.Source
	Sink   api.Sink
}

// Start initiate the collection process
func (c *Collector) Start() {
	c.Source.Start()
	c.collect()
}

func (c *Collector) collect() {
	for measure := range c.Source.MeasurementsChannel() {
		c.Sink.Put(measure)
	}
}
