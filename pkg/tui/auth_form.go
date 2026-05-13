package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/shiftwavedev/42-cli/pkg/credentials"
	"github.com/shiftwavedev/42-cli/pkg/display"
)

// AuthFormModel represents the auth login/update form
type AuthFormModel struct {
	inputs     []textinput.Model
	focusIndex int
	cursorMode textinput.EchoMode
	submitted  bool
	cancelled  bool
	validator  *credentials.Validator
	formTitle  string
	formType   string
}

// NewAuthForm creates a new auth form with empty fields
func NewAuthForm(formType string) *AuthFormModel {
	m := &AuthFormModel{
		inputs:     make([]textinput.Model, 3),
		cursorMode: textinput.EchoPassword,
		validator:  credentials.NewValidator(),
		formType:   formType,
	}

	if formType == "update" {
		m.formTitle = "Update Authentication Credentials"
	} else {
		m.formTitle = "42 API Authentication"
	}

	var fieldLabels = []string{"Login42", "Client ID", "Client Secret"}
	var echoModes = []textinput.EchoMode{
		textinput.EchoNormal,
		textinput.EchoNormal,
		textinput.EchoPassword,
	}

	for i := range m.inputs {
		t := textinput.New()
		t.Cursor.Style = lipgloss.NewStyle().Foreground(display.Primary)
		t.PromptStyle = lipgloss.NewStyle().Foreground(display.Primary)
		t.TextStyle = lipgloss.NewStyle().Foreground(display.Text)
		t.PlaceholderStyle = lipgloss.NewStyle().Foreground(display.Muted)
		t.EchoMode = echoModes[i]
		t.Placeholder = fieldLabels[i]

		m.inputs[i] = t
	}

	m.inputs[0].Focus()
	return m
}

// Init initializes the model
func (m *AuthFormModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles user input
func (m *AuthFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.cancelled = true
			return m, tea.Quit

		case tea.KeyEnter:
			if m.focusIndex == len(m.inputs)-1 {
				if m.validateForm() {
					m.submitted = true
					return m, tea.Quit
				}
			} else {
				m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
				m.updateInputFocus()
			}

		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}
			m.updateInputFocus()

		case tea.KeyTab, tea.KeyCtrlN:
			m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
			m.updateInputFocus()
		}

	case tea.WindowSizeMsg:
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

// View renders the form
func (m *AuthFormModel) View() string {
	var s strings.Builder

	s.WriteString(display.Header(m.formTitle, ""))
	s.WriteString("\n" + display.Divider(0) + "\n\n")

	fieldLabels := []string{"Login42", "Client ID", "Client Secret"}
	fieldDescriptions := []string{
		"Your 42 intranet login (eg: norminette)",
		"Your OAuth client ID (u-s4t2ud-xxxxx)",
		"Your OAuth client secret (s-s4t2ud-xxxxx)",
	}

	for i := range m.inputs {
		s.WriteString(renderFormField(fieldLabels[i], fieldDescriptions[i], m.inputs[i], i == m.focusIndex))
		s.WriteString("\n")
	}

	s.WriteString("\n" + display.Divider(0) + "\n\n")
	s.WriteString(display.Indent + display.RenderIf(display.Subtle, "Tab/Shift+Tab to navigate • Enter to confirm • Ctrl+C to cancel"))
	s.WriteString("\n")

	return s.String()
}

// GetCredentials returns the validated credentials if submitted
func (m *AuthFormModel) GetCredentials() *credentials.Credentials {
	if !m.submitted {
		return nil
	}

	return &credentials.Credentials{
		Login42:      m.inputs[0].Value(),
		ClientID:     m.inputs[1].Value(),
		ClientSecret: m.inputs[2].Value(),
	}
}

func (m *AuthFormModel) WasCancelled() bool {
	return m.cancelled
}

// Private helper functions

func (m *AuthFormModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *AuthFormModel) updateInputFocus() {
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.focusIndex {
			m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(display.Primary).Bold(true)
			m.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(display.Primary)
			m.inputs[i].Focus()
		} else {
			m.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(display.Muted)
			m.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(display.Text)
			m.inputs[i].Blur()
		}
	}
}

func (m *AuthFormModel) validateForm() bool {
	login := m.inputs[0].Value()
	clientID := m.inputs[1].Value()
	secret := m.inputs[2].Value()

	if login == "" || clientID == "" || secret == "" {
		return false
	}

	creds := &credentials.Credentials{
		Login42:      login,
		ClientID:     clientID,
		ClientSecret: secret,
	}

	return m.validator.Validate(creds)
}

func renderFormField(label, description string, input textinput.Model, focused bool) string {
	var s string

	labelStyle := lipgloss.NewStyle().
		Foreground(display.Text).
		Bold(true)
	s += labelStyle.Render(label) + "\n"

	s += display.Indent + display.RenderIf(display.Subtle, description) + "\n"

	inputStyle := lipgloss.NewStyle().
		Margin(0, 1)

	if focused {
		inputStyle = inputStyle.Border(lipgloss.RoundedBorder(), true).
			BorderForeground(display.Primary)
	}

	s += inputStyle.Render(input.View()) + "\n"

	return s
}

func RunAuthForm(formType string) (*credentials.Credentials, error) {
	model := NewAuthForm(formType)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("form error: %w", err)
	}

	if m, ok := finalModel.(*AuthFormModel); ok {
		if m.WasCancelled() {
			return nil, fmt.Errorf("form cancelled by user")
		}
		return m.GetCredentials(), nil
	}

	return nil, fmt.Errorf("unexpected form state")
}
