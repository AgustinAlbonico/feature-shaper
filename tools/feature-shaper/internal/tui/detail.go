package tui

import (
	"fmt"

	"github.com/agustinalbonico/feature-shaper/internal/db"
	"github.com/charmbracelet/lipgloss"
)

// renderDetail genera la vista de detalle de una feature con viewport scrollable
func renderDetail(m Model) string {
	if len(m.features) == 0 || m.activeFeature >= len(m.features) {
		return renderEmptyDetail(m)
	}

	feature := m.features[m.activeFeature]

	// Header con información de la feature
	header := renderDetailHeader(m, feature)

	// Body con contenido scrollable
	body := renderDetailBody(m)

	// Footer con acciones
	footer := renderDetailFooter(m)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

// renderEmptyDetail muestra un mensaje cuando no hay feature seleccionada
func renderEmptyDetail(m Model) string {
	header := lipgloss.NewStyle().
		Background(colorPrimary).
		Foreground(colorWhite).
		Width(m.width).
		Padding(0, 1).
		Render(logoStyle.Render(LogoSmall) + "  Detalle de feature")

	emptyMsg := lipgloss.NewStyle().
		Foreground(colorDimWhite).
		Padding(2, 0).
		Render("No hay feature seleccionada")

	body := lipgloss.NewStyle().
		Width(m.width - 4).
		Height(m.height - 6).
		Align(lipgloss.Center, lipgloss.Center).
		Render(emptyMsg)

	footer := footerStyle.Width(m.width - 2).Render("[b] volver")

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

// renderDetailHeader genera el header de la vista de detalle
func renderDetailHeader(m Model, feature db.Feature) string {
	// Título de la feature
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite)
	titleText := titleStyle.Render(truncate(feature.Title, m.width/2))

	// Badges
	typeIcon := TypeIcon(feature.Type)
	typeBadge := TypeColor(feature.Type).Render(fmt.Sprintf("%s %s", typeIcon, feature.Type))
	statusBadge := StatusColor(feature.Status).Render(fmt.Sprintf("%s %s", StatusIcon(feature.Status), feature.Status))
	versionBadge := lipgloss.NewStyle().
		Foreground(colorAccent).
		Render(fmt.Sprintf("v%d", feature.Version))

	// Primera línea: título
	titleLine := lipgloss.NewStyle().
		Width(m.width - 4).
		Render(titleText)

	// Segunda línea: badges
	badgesLine := lipgloss.NewStyle().
		Width(m.width - 4).
		Render(fmt.Sprintf("%s  %s  %s", typeBadge, statusBadge, versionBadge))

	// Tercera línea: keybindings
	keysLine := lipgloss.NewStyle().
		Width(m.width - 4).
		Align(lipgloss.Right).
		Foreground(colorSubtle).
		Render("[b] volver  [h] historial  [e] exportar")

	content := lipgloss.JoinVertical(lipgloss.Left, titleLine, badgesLine, keysLine)

	return lipgloss.NewStyle().
		Background(colorPrimary).
		Width(m.width).
		Padding(0, 1).
		Render(content)
}

// renderDetailBody genera el cuerpo de la vista de detalle
func renderDetailBody(m Model) string {
	// Estilo del viewport
	viewportStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorSubtle).
		Padding(1, 2).
		Width(m.width - 6).
		Height(m.height - 8)

	return viewportStyle.Render(m.viewport.View())
}

// renderDetailFooter genera el footer de la vista de detalle
func renderDetailFooter(m Model) string {
	scrollPercent := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)

	hints := []keyHint{
		{"↑↓/jk", "scroll"},
		{"h", "historial"},
		{"e", "exportar"},
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
