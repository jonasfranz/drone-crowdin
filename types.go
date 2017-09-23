package main

import "encoding/json"

// StringMapFlag represents a string-based map as a cli flag
type StringMapFlag struct {
	parts map[string]string
}

// String is implemented from cli
func (s *StringMapFlag) String() string {
	return ""
}

// Get returns the parsed map
func (s *StringMapFlag) Get() map[string]string {
	return s.parts
}

// Set parses the map (via json)
func (s *StringMapFlag) Set(value string) error {
	s.parts = map[string]string{}
	err := json.Unmarshal([]byte(value), &s.parts)
	if err != nil {
		s.parts["*"] = value
	}
	return nil
}
