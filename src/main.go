package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog"
)

// Create global logger
var logger zerolog.Logger

// hold the supplied run-time arguments globally
var glb_arguments config

// version information, are injected during build process
var (
	version string = "n/a"
	commit  string = "n/a"
	date    string = "0"
	gitdate string = "0"
	branch  string = "n/a"
)

func init() {
	loadConfig()
	parseArgs()
	configureLogger(glb_arguments.Options.LogLevel)

	if glb_arguments.Reporter.Pushover.Enabled {
		if len(glb_arguments.Reporter.Pushover.APIToken) == 0 {
			logger.Fatal().Msg("Pushover Enabled. Pushover API token required!")
		}
		if len(glb_arguments.Reporter.Pushover.UserKey) == 0 {
			logger.Fatal().Msg("Pushover Enabled. Pushover user key required!")
		}
	}
	if glb_arguments.Reporter.Gotify.Enabled {
		if len(glb_arguments.Reporter.Gotify.URL) == 0 {
			logger.Fatal().Msg("Gotify Enabled. Gotify URL required!")
		}
		if len(glb_arguments.Reporter.Gotify.Token) == 0 {
			logger.Fatal().Msg("Gotify Enabled. Gotify APP token required!")
		}
	}
	if glb_arguments.Reporter.Mail.Enabled {
		if len(glb_arguments.Reporter.Mail.User) == 0 {
			logger.Fatal().Msg("E-Mail notification Enabled. SMTP username required!")
		}
		if len(glb_arguments.Reporter.Mail.To) == 0 {
			logger.Fatal().Msg("E-Mail notification Enabled. Recipient address required!")
		}
		if len(glb_arguments.Reporter.Mail.From) == 0 {
			glb_arguments.Reporter.Mail.From = glb_arguments.Reporter.Mail.User
		}
		if len(glb_arguments.Reporter.Mail.Password) == 0 {
			logger.Fatal().Msg("E-Mail notification Enabled. SMTP Password required!")
		}
		if len(glb_arguments.Reporter.Mail.Host) == 0 {
			logger.Fatal().Msg("E-Mail notification Enabled. SMTP host address required!")
		}
	}
	if glb_arguments.Reporter.Mattermost.Enabled {
		if len(glb_arguments.Reporter.Mattermost.URL) == 0 {
			logger.Fatal().Msg("Mattermost Enabled. Mattermost URL required!")
		}
	}
}

func main() {
	// if the -v flag was set, print version information and exit
	if glb_arguments.Version {
		printVersion()
	}

	// log all supplied arguments
	logArguments()

	timestamp := time.Now()
	startup_message := buildStartupMessage(timestamp)
	sendNotifications(timestamp, startup_message, "Starting docker event monitor", glb_arguments.Reporters)

	filterArgs := filters.NewArgs()
	for key, values := range glb_arguments.Filter {
		for _, value := range values {
			filterArgs.Add(key, value)
		}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create new docker client")
	}
	defer cli.Close()

	// receives events from the channel
	event_chan, errs := cli.Events(context.Background(), types.EventsOptions{Filters: filterArgs})

	for {
		select {
		case err := <-errs:
			logger.Fatal().Err(err).Msg("")
		case event := <-event_chan:
			// if logging level is debug, log the event
			logger.Debug().
				Interface("event", event).Msg("")

			// Check if event should be exlcuded from reporting
			if len(glb_arguments.Exclude) > 0 {
				logger.Debug().Msg("Performing check for event exclusion")
				if excludeEvent(event) {
					break //breaks out of the select and waits for the next event to arrive
				}
			}
			processEvent(event)
		}
	}
}

func loadConfig() {
	configFile, err := filepath.Abs("./config.yml")

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to set config file path")
	}

	buf, err := os.ReadFile(configFile)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to read config file")
	}

	err = yaml.Unmarshal(buf, &glb_arguments)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse config file")
	}
}

func parseArgs() {

	// Parse (include) filters
	glb_arguments.Filter = make(map[string][]string)

	for _, filter := range glb_arguments.Options.FilterStrings {
		pos := strings.Index(filter, "=")
		if pos == -1 {
			logger.Fatal().Msg("each filter should be of the form key=value")
		}
		key := filter[:pos]
		val := filter[pos+1:]
		glb_arguments.Filter[key] = append(glb_arguments.Filter[key], val)
	}

	// Parse exclude filters
	glb_arguments.Exclude = make(map[string][]string)

	for _, exclude := range glb_arguments.Options.ExcludeStrings {
		pos := strings.Index(exclude, "=")
		if pos == -1 {
			logger.Fatal().Msg("each filter should be of the form key=value")
		}
		//trim whitespaces
		key := strings.TrimSpace(exclude[:pos])
		val := exclude[pos+1:]
		glb_arguments.Exclude[key] = append(glb_arguments.Exclude[key], val)
	}

	//Parse Enabled reportes

	if glb_arguments.Reporter.Gotify.Enabled {
		glb_arguments.Reporters = append(glb_arguments.Reporters, "Gotify")
	}
	if glb_arguments.Reporter.Mattermost.Enabled {
		glb_arguments.Reporters = append(glb_arguments.Reporters, "Mattermost")
	}
	if glb_arguments.Reporter.Pushover.Enabled {
		glb_arguments.Reporters = append(glb_arguments.Reporters, "Pushover")
	}
	if glb_arguments.Reporter.Mail.Enabled {
		glb_arguments.Reporters = append(glb_arguments.Reporters, "Mail")
	}

}

func configureLogger(LogLevel string) {
	// Configure time/timestamp format
	zerolog.TimeFieldFormat = time.RFC1123Z

	// Change logging level when debug flag is set
	if LogLevel == "debug" {
		logger = zerolog.New(os.Stdout).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Str("service", "docker event monitor").
			Logger()
	} else {
		logger = zerolog.New(os.Stdout).
			Level(zerolog.InfoLevel).
			With().
			Str("service", "docker event monitor").
			Timestamp().
			Logger()
	}
}
