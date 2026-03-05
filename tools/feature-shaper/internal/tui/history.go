package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderHistory genera la vista de historial de versiones
func renderHistory(m Model) string {
	featureTitle := ""
	if len(m.features) > 0 && m.activeFeature < len(m.features) {
		featureTitle = m.features[m.activeFeature].Title
	}

	header := renderHistoryHeader(m, featureTitle)
	body := renderHistoryBody(m)
	footer := renderHistoryFooter(m)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

// renderHistoryHeader genera el header de la vista de historial
func renderHistoryHeader(m Model, featureTitle string) string {
	// Título
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Render(fmt.Sprintf("📜 Historial — %s", truncate(featureTitle, m.width/3)))

	// Keybindings
	keys := lipgloss.NewStyle().
		Foreground(colorSubtle).
		Render("[b] volver")

	// Primera línea: título
	titleLine := lipgloss.NewStyle().
		Width(m.width - 4).
		Render(title)

	// Segunda línea: subtítulo
	subtitle := lipgloss.NewStyle().
		Foreground(colorDimWhite).
		Render(fmt.Sprintf("%d versiones guardadas", len(m.versions)))
	subtitleLine := lipgloss.NewStyle().
		Width(m.width - 12).
		Render(subtitle)

	keysLine := lipgloss.NewStyle().
		Width(10).
		Align(lipgloss.Right).
		Render(keys)

	secondLine := lipgloss.JoinHorizontal(lipgloss.Top, subtitleLine, keysLine)

	content := lipgloss.JoinVertical(lipgloss.Left, titleLine, secondLine)

	return lipgloss.NewStyle().
		Background(colorPrimary).
		Width(m.width).
		Padding(0, 1).
		Render(content)
}

// renderHistoryBody genera el cuerpo de la vista de historial
func renderHistoryBody(m Model) string {
	var sb strings.Builder

	if len(m.versions) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(colorDimWhite).
			Padding(2, 0).
			Render("No hay versiones anteriores.\nLa versión actual es la primera.")

		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSubtle).
			Padding(1, 2).
			Width(m.width - 6).
			Height(m.height - 8).
			Align(lipgloss.Center, lipgloss.Center).
			Render(emptyMsg)
	}

	// Cabecera de la tabla
	header := lipgloss.NewStyle().
		Foreground(colorDimWhite).
		Bold(true).
		Render(fmt.Sprintf("  %-6s  %-12s  %s", "VER", "FECHA", "CAMBIOS"))
	sb.WriteString(header)
	sb.WriteString("\n")
	sb.WriteString(dividerStyle.Render(strings.Repeat("─", m.width-8)))
	sb.WriteString("\n")

	for i, v := range m.versions {
		var line string

		// Extraer solo la fecha (primeros 10 chars del timestamp)
		date := v.CreatedAt
		if len(date) >= 10 {
			date = date[:10]
		}

		changelog := v.Changelog
		if changelog == "" {
			changelog = "(sin descripción)"
		}
		changelog = truncate(changelog, m.width-35)

		versionBadge := lipgloss.NewStyle().
			Foreground(colorAccent).
			Render(fmt.Sprintf("v%-4d", v.Version))

		if i == m.activeVersion {
			// Item seleccionado
			selectedVersion := lipgloss.NewStyle().
				Foreground(colorWhite).
				Bold(true).
				Render(versionBadge)
			line = fmt.Sprintf("▶ %s  %s  %s", selectedVersion, date, changelog)
		} else {
			line = fmt.Sprintf("  %s  %s  %s", versionBadge, date, changelog)
		}

		sb.WriteString(line)
		if i < len(m.versions)-1 {
			sb.WriteString("\n")
		}
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorSubtle).
		Padding(1, 2).
		Width(m.width - 6).
		Height(m.height - 8).
		Render(sb.String())
}

// renderHistoryFooter genera el footer de la vista de historial
func renderHistoryFooter(m Model) string {
	hints := []keyHint{
		{"↑↓", "navegar"},
		{"enter", "ver versión"},
		{"b", "volver"},
	}

	return footerStyle.Width(m.width - 2).Render(renderKeyHints(hints))
}

// renderVersionDetail genera la vista de detalle de una versión específica
func renderVersionDetail(m Model) string {
	versionLabel := ""
	if len(m.versions) > 0 && m.activeVersion < len(m.versions) {
		versionLabel = fmt.Sprintf("v%d", m.versions[m.activeVersion].Version)
	}

	featureTitle := ""
	if len(m.features) > 0 && m.activeFeature < len(m.features) {
		featureTitle = m.features[m.activeFeature].Title
	}

	header := renderVersionDetailHeader(m, featureTitle, versionLabel)
	body := renderVersionDetailBody(m)
	footer := renderVersionDetailFooter(m)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

// renderVersionDetailHeader genera el header de la vista de versión
func renderVersionDetailHeader(m Model, featureTitle, versionLabel string) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Render(fmt.Sprintf("%s — %s", truncate(featureTitle, m.width/3), versionLabel))

	keys := lipgloss.NewStyle().
		Foreground(colorSubtle).
		Render("[b] volver")

	titleLine := lipgloss.NewStyle().
		Width(m.width - 4).
		Render(title)

	keysLine := lipgloss.NewStyle().
		Width(m.width - 4).
		Align(lipgloss.Right).
		Render(keys)

	content := lipgloss.JoinVertical(lipgloss.Left, titleLine, keysLine)

	return lipgloss.NewStyle().
		Background(colorPrimary).
		Width(m.width).
		Padding(0, 1).
		Render(content)
}

// renderVersionDetailBody genera el cuerpo de la vista de versión
func renderVersionDetailBody(m Model) string {
	viewportStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorSubtle).
		Padding(1, 2).
		Width(m.width - 6).
		Height(m.height - 8)

	return viewportStyle.Render(m.viewport.View())
}

// renderVersionDetailFooter genera el footer de la vista de versión
func renderVersionDetailFooter(m Model) string {
	scrollPercent := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)

	hints := []keyHint{
		{"↑↓/jk", "scroll"},
		{"b", "volver"},
	}

	leftPart := renderKeyHints(hints)
	rightPart := lipgloss.NewStyle().
		Foreground(colorAccent).
		Render(scrollPercent)

	footerContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(m.width - 10).Render(leftPart),
		lipgloss.NewStyle().Width(8).Align(lipgloss.Right).Render(rightPart),
	)

	return footerStyle.Width(m.width - 2).Render(footerContent)
}
