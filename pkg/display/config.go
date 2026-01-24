package display

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Config holds global display configuration
type Config struct {
	// NoColor disables all ANSI color output
	NoColor bool
}

// GlobalConfig is the global display configuration
var GlobalConfig = Config{
	NoColor: false,
}

// InitConfig initializes display configuration from environment and flags
func InitConfig(noColor bool) {
	// NoColor: flag > env > default
	if noColor {
		GlobalConfig.NoColor = true
	} else if os.Getenv("NO_COLOR") != "" || os.Getenv("CI") != "" {
		GlobalConfig.NoColor = true
	}

	// Update the package-level NoColor for backward compatibility
	NoColor = GlobalConfig.NoColor

	// If NoColor is enabled, set lipgloss to ASCII mode
	if GlobalConfig.NoColor {
		lipgloss.SetColorProfile(termenv.Ascii)
	}
}
