package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/moodwave/moodwave/internal/config"
)

// DoctorCheckStatus represents the result of a single diagnostic check.
type DoctorCheckStatus int

const (
	CheckPending DoctorCheckStatus = iota
	CheckRunning
	CheckPass
	CheckFail
	CheckWarn
	CheckFixing
	CheckFixed
)

// DoctorCheck is a single diagnostic item.
type DoctorCheck struct {
	Category string
	Label    string
	Status   DoctorCheckStatus
	Detail   string
	CanFix   bool
	FixLabel string
}

// DoctorModel is the Bubble Tea model for the doctor TUI.
type DoctorModel struct {
	Width    int
	Height   int
	Frame    int
	Checks   []DoctorCheck
	Phase    int // which batch of checks is running
	Done     bool
	Quitting bool
	FixingID int // index of check being fixed, -1 if none

	// Callbacks
	RunChecks func() []DoctorCheck
	FixCheck  func(idx int) error
}

// DoctorTickMsg advances the animation frame.
type DoctorTickMsg time.Time

// DoctorChecksCompleteMsg signals checks are done.
type DoctorChecksCompleteMsg struct {
	Checks []DoctorCheck
}

// DoctorFixCompleteMsg signals a fix attempt finished.
type DoctorFixCompleteMsg struct {
	Idx int
	Err error
}

// NewDoctorModel creates the doctor TUI model.
func NewDoctorModel() DoctorModel {
	return DoctorModel{
		FixingID: -1,
	}
}

func (m DoctorModel) Init() tea.Cmd {
	return tea.Batch(
		doctorTickCmd(),
		tea.WindowSize(),
		m.runChecksCmd(),
	)
}

func doctorTickCmd() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return DoctorTickMsg(t)
	})
}

func (m DoctorModel) runChecksCmd() tea.Cmd {
	return func() tea.Msg {
		if m.RunChecks == nil {
			return DoctorChecksCompleteMsg{Checks: nil}
		}
		checks := m.RunChecks()
		return DoctorChecksCompleteMsg{Checks: checks}
	}
}

func (m DoctorModel) fixCheckCmd(idx int) tea.Cmd {
	return func() tea.Msg {
		if m.FixCheck == nil {
			return DoctorFixCompleteMsg{Idx: idx, Err: fmt.Errorf("no fix handler")}
		}
		err := m.FixCheck(idx)
		return DoctorFixCompleteMsg{Idx: idx, Err: err}
	}
}

func (m DoctorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case DoctorTickMsg:
		m.Frame++
		return m, doctorTickCmd()

	case DoctorChecksCompleteMsg:
		m.Checks = msg.Checks
		m.Done = true
		return m, nil

	case DoctorFixCompleteMsg:
		if msg.Idx >= 0 && msg.Idx < len(m.Checks) {
			if msg.Err == nil {
				m.Checks[msg.Idx].Status = CheckFixed
				m.Checks[msg.Idx].Detail = "Fixed successfully"
			} else {
				m.Checks[msg.Idx].Status = CheckFail
				m.Checks[msg.Idx].Detail = "Fix failed: " + msg.Err.Error()
			}
		}
		m.FixingID = -1
		return m, nil

	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "q", "esc", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "f":
			// Fix the first fixable failed check
			if m.FixingID >= 0 {
				return m, nil // already fixing
			}
			for i, c := range m.Checks {
				if c.CanFix && (c.Status == CheckFail || c.Status == CheckWarn) {
					m.FixingID = i
					m.Checks[i].Status = CheckFixing
					m.Checks[i].Detail = "Installing..."
					return m, m.fixCheckCmd(i)
				}
			}
		case "enter":
			if m.Done {
				m.Quitting = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m DoctorModel) View() string {
	if m.Quitting {
		return ""
	}
	if m.Width == 0 {
		return "Loading..."
	}

	var sections []string

	// Header
	header := m.renderHeader()
	sections = append(sections, header)
	sections = append(sections, "")

	if !m.Done {
		// Loading state
		sections = append(sections, m.renderLoading())
	} else {
		// Results
		sections = append(sections, m.renderResults())
	}

	// Footer hints
	sections = append(sections, "")
	sections = append(sections, m.renderFooter())

	content := strings.Join(sections, "\n")

	// Fill remaining height
	contentLines := strings.Count(content, "\n") + 1
	if m.Height > contentLines+1 {
		content += strings.Repeat("\n", m.Height-contentLines-2)
	}
	// Status bar
	content += "\n" + m.renderStatusBar()

	return content
}

func (m DoctorModel) renderHeader() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(ColorHighlight).
		Bold(true)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(ColorDim)

	title := titleStyle.Render("  ◇ Moodwave Doctor")
	subtitle := subtitleStyle.Render("  System diagnostics & auto-repair")

	borderStyle := lipgloss.NewStyle().
		Foreground(ColorBorder)
	line := borderStyle.Render("  " + strings.Repeat("─", m.Width-4))

	return title + "\n" + subtitle + "\n" + line
}

func (m DoctorModel) renderLoading() string {
	spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spin := spinChars[m.Frame%len(spinChars)]

	style := lipgloss.NewStyle().Foreground(ColorHighlight)
	dimStyle := lipgloss.NewStyle().Foreground(ColorDim)

	return style.Render(fmt.Sprintf("  %s  Running diagnostics...", spin)) + "\n" +
		dimStyle.Render("     Checking system, audio, network...")
}

func (m DoctorModel) renderResults() string {
	var lines []string
	currentCategory := ""

	passIcon := lipgloss.NewStyle().Foreground(ColorSuccess).Render("●")
	failIcon := lipgloss.NewStyle().Foreground(ColorError).Render("●")
	warnIcon := lipgloss.NewStyle().Foreground(ColorWarm).Render("●")
	fixingIcon := lipgloss.NewStyle().Foreground(ColorHighlight).Render("◌")
	fixedIcon := lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("◉")

	catStyle := lipgloss.NewStyle().Foreground(ColorFg).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(ColorFg)
	detailStyle := lipgloss.NewStyle().Foreground(ColorDim)
	fixStyle := lipgloss.NewStyle().Foreground(ColorHighlight).Italic(true)

	for _, c := range m.Checks {
		if c.Category != currentCategory {
			if currentCategory != "" {
				lines = append(lines, "")
			}
			currentCategory = c.Category
			lines = append(lines, "  "+catStyle.Render(c.Category))
		}

		var icon string
		switch c.Status {
		case CheckPass:
			icon = passIcon
		case CheckFail:
			icon = failIcon
		case CheckWarn:
			icon = warnIcon
		case CheckFixing:
			spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
			icon = lipgloss.NewStyle().Foreground(ColorHighlight).
				Render(spinChars[m.Frame%len(spinChars)])
		case CheckFixed:
			icon = fixedIcon
		default:
			icon = fixingIcon
		}

		line := fmt.Sprintf("    %s %s", icon, labelStyle.Render(c.Label))
		if c.Detail != "" {
			line += " " + detailStyle.Render("— "+c.Detail)
		}
		if c.CanFix && (c.Status == CheckFail || c.Status == CheckWarn) {
			line += " " + fixStyle.Render("[press F to fix]")
		}
		if c.Status == CheckFixing {
			line += " " + fixStyle.Render(c.Detail)
		}
		lines = append(lines, line)
	}

	// Summary
	lines = append(lines, "")
	pass, fail, warn := 0, 0, 0
	for _, c := range m.Checks {
		switch c.Status {
		case CheckPass, CheckFixed:
			pass++
		case CheckFail:
			fail++
		case CheckWarn:
			warn++
		}
	}

	sumStyle := lipgloss.NewStyle().Foreground(ColorDim)
	pStyle := lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	fStyle := lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	wStyle := lipgloss.NewStyle().Foreground(ColorWarm).Bold(true)

	summary := fmt.Sprintf("  %s %s  %s %s  %s %s",
		pStyle.Render(fmt.Sprintf("%d", pass)), sumStyle.Render("passed"),
		fStyle.Render(fmt.Sprintf("%d", fail)), sumStyle.Render("failed"),
		wStyle.Render(fmt.Sprintf("%d", warn)), sumStyle.Render("warnings"))
	lines = append(lines, summary)

	// Render in a panel
	content := strings.Join(lines, "\n")
	panelWidth := m.Width - 4
	if panelWidth > 90 {
		panelWidth = 90
	}
	panel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 1).
		Width(panelWidth).
		Render(content)

	return lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, panel)
}

func (m DoctorModel) renderFooter() string {
	k := lipgloss.NewStyle().Foreground(ColorHighlight).Bold(true)
	d := lipgloss.NewStyle().Foreground(ColorDim)
	sep := d.Render("  ·  ")

	hasFixable := false
	for _, c := range m.Checks {
		if c.CanFix && (c.Status == CheckFail || c.Status == CheckWarn) {
			hasFixable = true
			break
		}
	}

	hints := k.Render("[Q/Esc]") + d.Render(" Quit")
	if hasFixable {
		hints = k.Render("[F]") + d.Render(" Auto-Fix") + sep + hints
	}
	if m.Done {
		hints = k.Render("[Enter]") + d.Render(" Done") + sep + hints
	}

	return lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, hints)
}

func (m DoctorModel) renderStatusBar() string {
	left := lipgloss.NewStyle().Foreground(ColorDim).Render("  ◈ moodwave doctor")
	// Real build version, not a hardcoded literal (see app.go renderStatusBar).
	version := lipgloss.NewStyle().Foreground(ColorDim).
		Render("v" + strings.TrimPrefix(config.Version, "v"))

	state := "diagnosing..."
	if m.Done {
		state = "complete"
	}
	if m.FixingID >= 0 {
		state = "fixing..."
	}
	stateRender := lipgloss.NewStyle().Foreground(ColorHighlight).Render(state)

	gap := m.Width - lipgloss.Width(left) - lipgloss.Width(stateRender) - lipgloss.Width(version) - 6
	if gap < 0 {
		gap = 0
	}

	barStyle := lipgloss.NewStyle().Background(lipgloss.Color("#111111")).Width(m.Width)
	bar := left + strings.Repeat(" ", gap) + stateRender + "  " + version + "  "
	return barStyle.Render(bar)
}

// RunDoctor launches the doctor TUI. The caller provides:
// - runChecks: function that performs all diagnostic checks and returns results
// - fixCheck: function that attempts to fix a specific check by index
func RunDoctor(runChecks func() []DoctorCheck, fixCheck func(idx int) error) error {
	m := NewDoctorModel()
	m.RunChecks = runChecks
	m.FixCheck = fixCheck
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
