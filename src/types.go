package main

type options struct {
	FilterStrings  []string `yaml:"filter_strings,flow"`
	ExcludeStrings []string `yaml:"exclude_strings,flow"`
	LogLevel       string   `yaml:"log_level"`
	ServerTag      string   `yaml:"server_tag"`
}

type Config struct {
	Options         options
	EnabledReporter []string            `yaml:"-"`
	Filter          map[string][]string `yaml:"-"`
	Exclude         map[string][]string `yaml:"-"`
}
