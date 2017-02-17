package model

import (
	"fmt"
	"sync"
)

var extractors = make(map[string]Extractor)
var rwMutex sync.RWMutex

// Extractor extracts tests from a domain-specific test results file
type Extractor interface {
	Extract(filepath string) ([]Test, error)
}

// Extract extracts test results from a results file of the given format
func Extract(format string, filepath string) ([]Test, error) {
	rwMutex.RLock()
	extractor := extractors[format]
	rwMutex.RUnlock()

	if extractor == nil {
		return nil, fmt.Errorf("no extractor for format '%s'", format)
	}

	return extractor.Extract(filepath)
}

// RegisterExtractor registers a new extractor
func RegisterExtractor(name string, extractor Extractor) error {
	rwMutex.RLock()
	currentExtractor := extractors[name]
	rwMutex.RUnlock()
    
	if currentExtractor != nil {
		return fmt.Errorf("extract with name '%s' already registered", name)
	}

	rwMutex.Lock()
	defer rwMutex.Unlock()
	extractors[name] = extractor

	return nil
}
