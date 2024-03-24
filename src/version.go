package main

import "os"

// version information, are injected during build process
var (
	version string = "n/a"
	commit  string = "n/a"
	date    string = "0"
	gitdate string = "0"
	branch  string = "n/a"
)

func printVersion() {
	logger.Info().
		Str("Version", version).
		Str("Branch", branch).
		Str("Commit", commit).
		Time("Compile_date", stringToUnix(date)).
		Time("Git_date", stringToUnix(gitdate)).
		Msg("Version Information")
	os.Exit(0)
}
