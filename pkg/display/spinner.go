package display

import (
	"fmt"
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

// NewSpinner creates a new spinner with a message
func NewSpinner(message string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	if !NoColor {
		s.Style = lipgloss.NewStyle().Foreground(Primary)
	}
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

// Spinner provides a simple API for showing loading spinners
type Spinner struct {
	program *tea.Program
	done    chan struct{}
}

// StartSpinner starts a spinner with a message and returns a Spinner that can be stopped
func StartSpinner(message string) *Spinner {
	model := NewSpinner(message)
	p := tea.NewProgram(model, tea.WithOutput(nil))

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

// Stop stops the spinner
func (s *Spinner) Stop() {
	if s.program != nil {
		s.program.Send(spinnerDoneMsg{})
		<-s.done
	}
}

// WithSpinner executes a function while showing a spinner
// Returns the result of the function
func WithSpinner[T any](message string, fn func() (T, error)) (T, error) {
	s := StartSpinner(message)
	defer s.Stop()

	result, err := fn()

	// Clear the spinner line
	fmt.Print("\r\033[K")

	return result, err
}

// WithSpinnerNoResult executes a function while showing a spinner (no return value)
func WithSpinnerNoResult(message string, fn func() error) error {
	s := StartSpinner(message)
	defer s.Stop()

	err := fn()

	// Clear the spinner line
	fmt.Print("\r\033[K")

	return err
}

// SimpleSpinner provides a non-bubbletea spinner for simpler use cases
type SimpleSpinner struct {
	message string
	done    chan struct{}
	frames  []string
}

// NewSimpleSpinner creates a basic spinner without bubbletea
func NewSimpleSpinner(message string) *SimpleSpinner {
	return &SimpleSpinner{
		message: message,
		done:    make(chan struct{}),
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

// Start begins the spinner animation
func (s *SimpleSpinner) Start() {
	go func() {
		i := 0
		for {
			select {
			case <-s.done:
				return
			default:
				frame := s.frames[i%len(s.frames)]
				if NoColor {
					fmt.Printf("\r%s %s", frame, s.message)
				} else {
					styledFrame := lipgloss.NewStyle().Foreground(Primary).Render(frame)
					fmt.Printf("\r%s %s", styledFrame, s.message)
				}
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

// Stop stops the spinner and clears the line
func (s *SimpleSpinner) Stop() {
	close(s.done)
	fmt.Print("\r\033[K") // Clear the line
}

// StopWithMessage stops the spinner and prints a final message
func (s *SimpleSpinner) StopWithMessage(message string) {
	close(s.done)
	fmt.Print("\r\033[K") // Clear the line
	fmt.Println(message)
}

// StopWithSuccess stops the spinner with a success message
func (s *SimpleSpinner) StopWithSuccess(message string) {
	close(s.done)
	fmt.Print("\r\033[K")
	fmt.Println(Badge("✓", Success) + " " + message)
}

// StopWithError stops the spinner with an error message
func (s *SimpleSpinner) StopWithError(message string) {
	close(s.done)
	fmt.Print("\r\033[K")
	fmt.Println(Badge("✗", Error) + " " + message)
}
