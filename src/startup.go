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

	if config.Reporter.Pushover.Enabled {
		startup_message_builder.WriteString("Pushover notification Enabled")
	} else {
		startup_message_builder.WriteString("Pushover notification disabled")
	}

	if config.Reporter.Gotify.Enabled {
		startup_message_builder.WriteString("\nGotify notification Enabled")
	} else {
		startup_message_builder.WriteString("\nGotify notification disabled")
	}
	if config.Reporter.Mail.Enabled {
		startup_message_builder.WriteString("\nE-Mail notification Enabled")
	} else {
		startup_message_builder.WriteString("\nE-Mail notification disabled")
	}

	if config.Reporter.Mattermost.Enabled {
		startup_message_builder.WriteString("\nMattermost notification Enabled")
		if config.Reporter.Mattermost.Channel != "" {
			startup_message_builder.WriteString("\nMattermost channel: " + config.Reporter.Mattermost.Channel)
		}
		if config.Reporter.Mattermost.User != "" {
			startup_message_builder.WriteString("\nMattermost username: " + config.Reporter.Mattermost.User)
		}
	} else {
		startup_message_builder.WriteString("\nMattermost notification disabled")
	}

	if config.Options.Delay > 0 {
		startup_message_builder.WriteString("\nUsing delay of " + config.Options.Delay.String())
	} else {
		startup_message_builder.WriteString("\nDelay disabled")
	}

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
		Dict("options", zerolog.Dict().
			Dict("reporter", zerolog.Dict().
				Dict("Pushover", zerolog.Dict().
					Bool("Enabled", config.Reporter.Pushover.Enabled).
					Str("PushoverAPIToken", config.Reporter.Pushover.APIToken).
					Str("PushoverUserKey", config.Reporter.Pushover.UserKey),
				).
				Dict("Gotify", zerolog.Dict().
					Bool("Enabled", config.Reporter.Gotify.Enabled).
					Str("GotifyURL", config.Reporter.Gotify.URL).
					Str("GotifyToken", config.Reporter.Gotify.Token),
				).
				Dict("Mail", zerolog.Dict().
					Bool("Enabled", config.Reporter.Mail.Enabled).
					Str("MailFrom", config.Reporter.Mail.From).
					Str("MailTo", config.Reporter.Mail.To).
					Str("MailHost", config.Reporter.Mail.Host).
					Str("MailUser", config.Reporter.Mail.User).
					Int("Port", config.Reporter.Mail.Port),
				).
				Dict("Mattermost", zerolog.Dict().
					Bool("Enabled", config.Reporter.Mattermost.Enabled).
					Str("MattermostURL", config.Reporter.Mattermost.URL).
					Str("MattermostChannel", config.Reporter.Mattermost.Channel).
					Str("MattermostUser", config.Reporter.Mattermost.User),
				),
			).
			Str("Delay", config.Options.Delay.String()).
			Str("Loglevel", config.Options.LogLevel).
			Str("ServerTag", config.Options.ServerTag).
			Str("Filter", strings.Join(config.Options.FilterStrings, " ")).
			Str("Exclude", strings.Join(config.Options.ExcludeStrings, " ")),
		).
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
