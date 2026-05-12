package display

import (
	"os"

	"golang.org/x/term"
)

type TerminalDimensions struct {
	Width  int
	Height int
}

// GetTerminalDimensions detects terminal size with intelligent fallbacks
func GetTerminalDimensions() TerminalDimensions {
	if width, height, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		return TerminalDimensions{Width: width, Height: height}
	}

	if width, height, err := term.GetSize(int(os.Stderr.Fd())); err == nil {
		return TerminalDimensions{Width: width, Height: height}
	}

	return TerminalDimensions{Width: 80, Height: 24}
}

// GlobalTerminalDimensions caches the terminal size (set at startup)
var GlobalTerminalDimensions = GetTerminalDimensions()

// RefreshTerminalDimensions updates the cached size
func RefreshTerminalDimensions() {
	GlobalTerminalDimensions = GetTerminalDimensions()
}
