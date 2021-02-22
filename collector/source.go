package collector

import "github.com/renajohn/pac_collector/store"

// Source represents a source
type Source interface {
	Start() <-chan store.Measurement
}
