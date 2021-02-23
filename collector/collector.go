package collector

import "github.com/renajohn/pac_collector/api"

// Collector bind a source to a store
type Collector struct {
	Source        api.Source
	Store         api.Store
	sourceChannel <-chan api.Measurement
}

// Start initiate the collection process
func (c *Collector) Start() {
	c.sourceChannel = c.Source.Start()

	c.collect()
}

func (c *Collector) collect() {
	for measure := range c.sourceChannel {
		c.Store.Put(measure)
	}
}
