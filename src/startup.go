package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func buildStartupMessage(timestamp time.Time) string {
	var startup_message_builder strings.Builder

	startup_message_builder.WriteString("Docker event monitor started at " + timestamp.Format(time.RFC1123Z) + "\n")
	startup_message_builder.WriteString("Docker event monitor version: " + version + "\n")

	if glb_arguments.Reporter.Pushover.Enabled {
		startup_message_builder.WriteString("Pushover notification Enabled")
	} else {
		startup_message_builder.WriteString("Pushover notification disabled")
	}

	if glb_arguments.Reporter.Gotify.Enabled {
		startup_message_builder.WriteString("\nGotify notification Enabled")
	} else {
		startup_message_builder.WriteString("\nGotify notification disabled")
	}
	if glb_arguments.Reporter.Mail.Enabled {
		startup_message_builder.WriteString("\nE-Mail notification Enabled")
	} else {
		startup_message_builder.WriteString("\nE-Mail notification disabled")
	}

	if glb_arguments.Reporter.Mattermost.Enabled {
		startup_message_builder.WriteString("\nMattermost notification Enabled")
		if glb_arguments.Reporter.Mattermost.Channel != "" {
			startup_message_builder.WriteString("\nMattermost channel: " + glb_arguments.Reporter.Mattermost.Channel)
		}
		if glb_arguments.Reporter.Mattermost.User != "" {
			startup_message_builder.WriteString("\nMattermost username: " + glb_arguments.Reporter.Mattermost.User)
		}
	} else {
		startup_message_builder.WriteString("\nMattermost notification disabled")
	}

	if glb_arguments.Options.Delay > 0 {
		startup_message_builder.WriteString("\nUsing delay of " + glb_arguments.Options.Delay.String())
	} else {
		startup_message_builder.WriteString("\nDelay disabled")
	}

	startup_message_builder.WriteString("\nLog level: " + glb_arguments.Options.LogLevel)

	if glb_arguments.Options.ServerTag != "" {
		startup_message_builder.WriteString("\nServerTag: " + glb_arguments.Options.ServerTag)
	} else {
		startup_message_builder.WriteString("\nServerTag: none")
	}

	if len(glb_arguments.Options.FilterStrings) > 0 {
		startup_message_builder.WriteString("\nFilterStrings: " + strings.Join(glb_arguments.Options.FilterStrings, " "))
	} else {
		startup_message_builder.WriteString("\nFilterStrings: none")
	}

	if len(glb_arguments.Options.ExcludeStrings) > 0 {
		startup_message_builder.WriteString("\nExcludeStrings: " + strings.Join(glb_arguments.Options.ExcludeStrings, " "))
	} else {
		startup_message_builder.WriteString("\nExcludeStrings: none")
	}

	return startup_message_builder.String()
}

func logArguments() {
	logger.Info().
		Dict("options", zerolog.Dict().
			Dict("reporter", zerolog.Dict().
				Dict("Pushover", zerolog.Dict().
					Bool("Enabled", glb_arguments.Reporter.Pushover.Enabled).
					Str("PushoverAPIToken", glb_arguments.Reporter.Pushover.APIToken).
					Str("PushoverUserKey", glb_arguments.Reporter.Pushover.UserKey),
				).
				Dict("Gotify", zerolog.Dict().
					Bool("Enabled", glb_arguments.Reporter.Gotify.Enabled).
					Str("GotifyURL", glb_arguments.Reporter.Gotify.URL).
					Str("GotifyToken", glb_arguments.Reporter.Gotify.Token),
				).
				Dict("Mail", zerolog.Dict().
					Bool("Enabled", glb_arguments.Reporter.Mail.Enabled).
					Str("MailFrom", glb_arguments.Reporter.Mail.From).
					Str("MailTo", glb_arguments.Reporter.Mail.To).
					Str("MailHost", glb_arguments.Reporter.Mail.Host).
					Str("MailUser", glb_arguments.Reporter.Mail.User).
					Int("Port", glb_arguments.Reporter.Mail.Port),
				).
				Dict("Mattermost", zerolog.Dict().
					Bool("Enabled", glb_arguments.Reporter.Mattermost.Enabled).
					Str("MattermostURL", glb_arguments.Reporter.Mattermost.URL).
					Str("MattermostChannel", glb_arguments.Reporter.Mattermost.Channel).
					Str("MattermostUser", glb_arguments.Reporter.Mattermost.User),
				),
			).
			Str("Delay", glb_arguments.Options.Delay.String()).
			Str("Loglevel", glb_arguments.Options.LogLevel).
			Str("ServerTag", glb_arguments.Options.ServerTag).
			Str("Filter", strings.Join(glb_arguments.Options.FilterStrings, " ")).
			Str("Exclude", strings.Join(glb_arguments.Options.ExcludeStrings, " ")),
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
		logger.Fatal().Err(err).Msg("String to timestamp conversion failed")
	}
	tm := time.Unix(i, 0)
	return tm
}
