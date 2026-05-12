package display

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerModel wraps a spinner for API loading states
type SpinnerModel struct {
	spinner spinner.Model
	message string
	done    bool
}

// NewSpinnerModel creates a new spinner model with a message
func NewSpinnerModel(message string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(Primary)
	return SpinnerModel{
		spinner: s,
		message: message,
	}
}

// Init initializes the spinner
func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles spinner updates
func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case spinnerDoneMsg:
		m.done = true
		return m, tea.Quit
	}
	return m, nil
}

// View renders the spinner
func (m SpinnerModel) View() string {
	if m.done {
		return ""
	}
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

// spinnerDoneMsg signals the spinner to stop
type spinnerDoneMsg struct{}

// Spinner provides a unified API for showing loading spinners
// Works in both COLOR and NO_COLOR modes, renders to stderr to not pollute stdout
type Spinner struct {
	program *tea.Program
	done    chan struct{}
}

// NewSpinner creates a new unified spinner that works in all modes
func NewSpinner(message string) *Spinner {
	model := NewSpinnerModel(message)
	// Render to stderr so stdout is clean for piping (important for --json)
	p := tea.NewProgram(model, tea.WithOutput(os.Stderr))

	s := &Spinner{
		program: p,
		done:    make(chan struct{}),
	}

	go func() {
		p.Run()
		close(s.done)
	}()

	// Give the spinner time to start
	time.Sleep(10 * time.Millisecond)

	return s
}

// Start begins the spinner animation (kept for backward compatibility with goroutine version)
// This is handled internally by NewSpinner now
func (s *Spinner) Start() {
	// Already started in NewSpinner
}

func (s *Spinner) Stop() {
	if s.program != nil {
		s.program.Send(spinnerDoneMsg{})
		<-s.done
		fmt.Fprint(os.Stderr, "\r\033[K")
	}
}

// Done stops the spinner with a success message
// Outputs: ✓ message (with color in COLOR mode, without in NO_COLOR)
func (s *Spinner) Done(message string) {
	s.Stop()
	fmt.Printf("%s %s\n", Icon("success"), message)
}

// Fail stops the spinner with an error message
// Outputs: ✗ message (with color in COLOR mode, without in NO_COLOR)
func (s *Spinner) Fail(message string) {
	s.Stop()
	fmt.Printf("%s %s\n", Icon("error"), message)
}

// WithSpinner executes a function while showing a spinner
// Returns the result of the function
func WithSpinner[T any](message string, fn func() (T, error)) (T, error) {
	s := NewSpinner(message)
	defer s.Stop()

	result, err := fn()

	// Clear the spinner line
	fmt.Fprint(os.Stderr, "\r\033[K")

	return result, err
}

// WithSpinnerNoResult executes a function while showing a spinner (no return value)
func WithSpinnerNoResult(message string, fn func() error) error {
	s := NewSpinner(message)
	defer s.Stop()

	err := fn()

	// Clear the spinner line
	fmt.Fprint(os.Stderr, "\r\033[K")

	return err
}
