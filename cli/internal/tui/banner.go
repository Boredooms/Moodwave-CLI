package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Banner is the large ASCII art logo — block-pixel style like OpenCode.
const Banner = `
 ███╗   ███╗ ██████╗  ██████╗ ██████╗ ██╗    ██╗ █████╗ ██╗   ██╗███████╗
 ████╗ ████║██╔═══██╗██╔═══██╗██╔══██╗██║    ██║██╔══██╗██║   ██║██╔════╝
 ██╔████╔██║██║   ██║██║   ██║██║  ██║██║ █╗ ██║███████║██║   ██║█████╗  
 ██║╚██╔╝██║██║   ██║██║   ██║██║  ██║██║███╗██║██╔══██║╚██╗ ██╔╝██╔══╝  
 ██║ ╚═╝ ██║╚██████╔╝╚██████╔╝██████╔╝╚███╔███╔╝██║  ██║ ╚████╔╝ ███████╗
 ╚═╝     ╚═╝ ╚═════╝  ╚═════╝ ╚═════╝  ╚══╝╚══╝ ╚═╝  ╚═╝  ╚═══╝  ╚══════╝`

// CompactBanner is used when terminal width < 80.
const CompactBanner = `
 ╔╦╗╔═╗╔═╗╔╦╗╦ ╦╔═╗╦  ╦╔═╗
 ║║║║ ║║ ║ ║║║║║╠═╣╚╗╔╝║╣ 
 ╩ ╩╚═╝╚═╝═╩╝╚╩╝╩ ╩ ╚╝ ╚═╝`

// RenderBanner returns the styled banner with tagline.
func RenderBanner(width int, moodColor lipgloss.Color) string {
	var b strings.Builder

	bannerText := Banner
	if width < 80 {
		bannerText = CompactBanner
	}

	bannerStyle := lipgloss.NewStyle().
		Foreground(moodColor).
		Bold(true)

	taglineStyle := lipgloss.NewStyle().
		Foreground(ColorDim).
		Italic(true)

	rendered := bannerStyle.Render(bannerText)
	tagline := taglineStyle.Render("terminal mood music companion")

	// Center everything
	b.WriteString(rendered)
	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Center, tagline))

	return b.String()
}
