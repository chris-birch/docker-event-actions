package main

import "time"

type pushover struct {
	Enabled  bool
	APIToken string `yaml:"api_token"`
	UserKey  string `yaml:"user_key"`
}
type gotify struct {
	Enabled bool
	URL     string `yaml:"url"`
	Token   string `yaml:"token"`
}
type mail struct {
	Enabled  bool
	From     string `yaml:"from"`
	To       string `yaml:"to"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
}
type mattermost struct {
	Enabled bool
	URL     string `yaml:"url"`
	Channel string `yaml:"channel"`
	User    string `yaml:"user"`
}

type reporter struct {
	Pushover   pushover
	Gotify     gotify
	Mail       mail
	Mattermost mattermost
}

type options struct {
	FilterStrings  []string      `yaml:"filter_strings,flow"`
	ExcludeStrings []string      `yaml:"exclude_strings,flow"`
	LogLevel       string        `yaml:"log_level"`
	ServerTag      string        `yaml:"server_tag"`
	Delay          time.Duration `yaml:"delay"`
}

type config struct {
	Reporter  reporter
	Options   options
	Reporters []string            `yaml:"-"`
	Filter    map[string][]string `yaml:"-"`
	Exclude   map[string][]string `yaml:"-"`
}
