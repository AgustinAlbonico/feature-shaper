package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderSearch genera la vista de búsqueda FTS
func renderSearch(m Model) string {
	header := renderSearchHeader(m)
	input := renderSearchInput(m)
	results := renderSearchResults(m)
	footer := renderSearchFooter(m)

	return lipgloss.JoinVertical(lipgloss.Left, header, input, results, footer)
}

// renderSearchHeader genera el header de la vista de búsqueda
func renderSearchHeader(m Model) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Render("🔍 Búsqueda de features")

	// Contador de resultados
	resultCount := ""
	if m.searchInput.Value() != "" {
		resultCount = lipgloss.NewStyle().
			Foreground(colorAccent).
			Render(fmt.Sprintf("%d resultados", len(m.searchResults)))
	}

	keys := lipgloss.NewStyle().
		Foreground(colorSubtle).
		Render("[esc] cancelar")

	titleLine := lipgloss.NewStyle().
		Width(m.width/2 - 2).
		Render(title)

	countLine := lipgloss.NewStyle().
		Width(m.width/2 - 2).
		Align(lipgloss.Right).
		Render(resultCount + "  " + keys)

	content := lipgloss.JoinHorizontal(lipgloss.Top, titleLine, countLine)

	return lipgloss.NewStyle().
		Background(colorPrimary).
		Width(m.width).
		Padding(0, 1).
		Render(content)
}

// renderSearchInput genera el input de búsqueda estilizado
func renderSearchInput(m Model) string {
	// Estilo del input
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorAccent).
		Padding(0, 1).
		Width(m.width - 6).
		Background(lipgloss.Color("#1A1A2E"))

	// Placeholder cuando está vacío
	if m.searchInput.Value() == "" {
		placeholder := lipgloss.NewStyle().
			Foreground(colorSubtle).
			Render("Escribe para buscar en todas las features...")
		return inputStyle.Render(placeholder)
	}

	return inputStyle.Foreground(colorWhite).Render(m.searchInput.View())
}

// renderSearchResults genera la lista de resultados
func renderSearchResults(m Model) string {
	var sb strings.Builder

	if len(m.searchResults) == 0 {
		var emptyMsg string
		if m.searchInput.Value() != "" {
			emptyMsg = lipgloss.NewStyle().
				Foreground(colorDimWhite).
				Render("No se encontraron resultados para \"" + m.searchInput.Value() + "\"")
		} else {
			emptyMsg = lipgloss.NewStyle().
				Foreground(colorSubtle).
				Render("Escribe para buscar en títulos, contenido, tipos y estados...")
		}

		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSubtle).
			Padding(1, 2).
			Width(m.width - 6).
			Height(m.height - 12).
			Align(lipgloss.Center, lipgloss.Center).
			Render(emptyMsg)
	}

	// Anchos de columna dinámicos
	titleWidth := m.width / 3
	if titleWidth < 20 {
		titleWidth = 20
	}

	for i, r := range m.searchResults {
		var line string

		// Título con icono de tipo
		typeIcon := TypeIcon(r.Type)
		title := lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite).
			Render(fmt.Sprintf("%s %s", typeIcon, truncate(r.Title, titleWidth)))

		// Badges
		projectBadge := dimStyle.Render(fmt.Sprintf("[%s]", r.ProjectSlug))
		typ := TypeColor(r.Type).Render(r.Type)
		status := StatusColor(r.Status).Render(fmt.Sprintf("%s %s", StatusIcon(r.Status), r.Status))
		ver := dimStyle.Render(fmt.Sprintf("v%d", r.Version))

		line = fmt.Sprintf("  %s  %s  %s  %s  %s", title, projectBadge, typ, status, ver)
		sb.WriteString(line)

		// Preview del contenido
		if r.Preview != "" {
			preview := truncate(strings.ReplaceAll(r.Preview, "\n", " "), m.width-8)
			previewLine := lipgloss.NewStyle().
				Foreground(colorSubtle).
				Italic(true).
				Render(fmt.Sprintf("    %s", preview))
			sb.WriteString("\n")
			sb.WriteString(previewLine)
		}

		if i < len(m.searchResults)-1 {
			sb.WriteString("\n\n")
		}
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorSubtle).
		Padding(1, 2).
		Width(m.width - 6).
		Height(m.height - 12).
		Render(sb.String())
}

// renderSearchFooter genera el footer de la vista de búsqueda
func renderSearchFooter(m Model) string {
	hints := []keyHint{
		{"enter", "buscar"},
		{"esc", "cancelar"},
	}

	return footerStyle.Width(m.width - 2).Render(renderKeyHints(hints))
}
