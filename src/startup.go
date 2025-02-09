package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func buildStartupMessage(timestamp time.Time) string {
	var startup_message_builder strings.Builder

	startup_message_builder.WriteString("Docker event monitor started at " + timestamp.Format(time.RFC1123Z) + "\n")
	startup_message_builder.WriteString("Docker event monitor version: " + version + "\n")

	startup_message_builder.WriteString("\nLog level: " + config.Options.LogLevel)

	if config.Options.ServerTag != "" {
		startup_message_builder.WriteString("\nServerTag: " + config.Options.ServerTag)
	} else {
		startup_message_builder.WriteString("\nServerTag: none")
	}

	if len(config.Options.FilterStrings) > 0 {
		startup_message_builder.WriteString("\nFilterStrings: " + strings.Join(config.Options.FilterStrings, " "))
	} else {
		startup_message_builder.WriteString("\nFilterStrings: none")
	}

	if len(config.Options.ExcludeStrings) > 0 {
		startup_message_builder.WriteString("\nExcludeStrings: " + strings.Join(config.Options.ExcludeStrings, " "))
	} else {
		startup_message_builder.WriteString("\nExcludeStrings: none")
	}

	return startup_message_builder.String()
}

func logArguments() {
	log.Info().
		Interface("options", config.Options).
		Dict("version", zerolog.Dict().
			Str("Version", version).
			Str("Branch", branch).
			Str("Commit", commit).
			Time("Compile_date", stringToUnix(date)).
			Time("Git_date", stringToUnix(gitdate)),
		).
		Msg("Docker event monitor started")
}

func stringToUnix(str string) time.Time {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msg("String to timestamp conversion failed")
	}
	tm := time.Unix(i, 0)
	return tm
}
