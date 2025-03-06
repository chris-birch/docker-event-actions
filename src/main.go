package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/chris-birch/docker-event-actions/src/technitium"
	"github.com/docker/docker/api/types/events"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// hold config options and settings globally
var config Config

// should we only print version information and exit
var showVersion bool

// config file path
var configFilePath string

var PrintUsage = func() {
	fmt.Printf("Usage of %s", os.Args[0])
	fmt.Print(
		`
-v, --version		prints version information
-c, --config [path]	config file path (default "config.yml")
-h, --help		prints help information
`)
}

func init() {
	flag.BoolVar(&showVersion, "v", false, "print version information")
	flag.BoolVar(&showVersion, "version", false, "print version information")
	flag.StringVar(&configFilePath, "c", "config.yml", "config file path")
	flag.StringVar(&configFilePath, "config", "config.yml", "config file path")
	flag.Usage = PrintUsage
	flag.Parse()

	configureLogger()
	loadConfig()

	// after loading the config, check log level
	if config.Options.LogLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	parseArgs()
}

func loadConfig() {
	configFile, err := filepath.Abs(configFilePath)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to set config file path")
	}

	buf, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config file")
	}

	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config file")
	}
}

func parseArgs() {

	// Parse (include) filters
	config.Filter = make(map[string][]string)

	for _, filter := range config.Options.FilterStrings {
		pos := strings.Index(filter, "=")
		if pos == -1 {
			log.Fatal().Msg("each filter should be of the form key=value")
		}
		key := filter[:pos]
		val := filter[pos+1:]
		config.Filter[key] = append(config.Filter[key], val)
	}

	// Parse exclude filters
	config.Exclude = make(map[string][]string)

	for _, exclude := range config.Options.ExcludeStrings {
		pos := strings.Index(exclude, "=")
		if pos == -1 {
			log.Fatal().Msg("each filter should be of the form key=value")
		}
		//trim whitespaces
		key := strings.TrimSpace(exclude[:pos])
		val := exclude[pos+1:]
		config.Exclude[key] = append(config.Exclude[key], val)
	}
}

func configureLogger() {

	// Configure time/timestamp format
	zerolog.TimeFieldFormat = time.RFC1123Z

	// Default level is info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Add timestamp and pretty output
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}

func main() {
	// if the -v flag was set, print version information and exit
	if showVersion {
		printVersion()
	}

	// log all supplied arguments
	logArguments()

	filterArgs := filters.NewArgs()
	for key, values := range config.Filter {
		for _, value := range values {
			filterArgs.Add(key, value)
		}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create new docker client")
	}

	defer func(cli *client.Client) {
		err := cli.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to close docker client")
		}
		fmt.Println("here...")
	}(cli)

	// receives events from the channel
	eventChan, errs := cli.Events(context.Background(), events.ListOptions{Filters: filterArgs})

	// Setup Techchnitium package
	if config.DockHost == "" {
		log.Fatal().Msg("Docker host not set in config file")
	} else {
		log.Info().Msgf("Using docker host: %s", config.DockHost)
	}
	tech := new(technitium.Technitium)
	tech.Init()
	defer tech.Close()

	for {
		select {
		case err := <-errs:
			log.Fatal().Err(err).Msg("")
		case event := <-eventChan:
			// if logging level is debug, log the event
			log.Debug().
				Interface("event", event).Msg("")

			// Check if event should be excluded from reporting
			if len(config.Exclude) > 0 {
				log.Debug().Msg("Performing check for event exclusion")
				if excludeEvent(event) {
					break //breaks out of the select and waits for the next event to arrive
				}
			}

			// Prepare DNS record
			rec, err := technitium.NewRecord(event, cli, config.DockHost)
			if err != nil {
				log.Err(err).Msg("Failed to create technitium record")
			} else {
				// Send record to gRPC server if we don't have any record errors
				if rec != nil {
					tech.SendMsg(rec)
				}
			}
		}
	}
}
