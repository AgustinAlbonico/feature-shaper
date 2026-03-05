package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderCatalog genera la vista de catálogo con dos paneles (proyectos + features)
func renderCatalog(m Model) string {
	leftWidth := m.width/4 - 2
	if leftWidth < 24 {
		leftWidth = 24
	}
	rightWidth := m.width - leftWidth - 6

	// Header con logo y estadísticas
	header := renderCatalogHeader(m)

	// Panel izquierdo: proyectos
	leftContent := renderProjectList(m, leftWidth)
	leftPanel := panelBorder.Width(leftWidth)
	if m.focusLeft {
		leftPanel = activePanelBorder.Width(leftWidth)
	}

	// Panel derecho: features
	rightContent := renderFeatureList(m, rightWidth)
	rightPanel := panelBorder.Width(rightWidth)
	if !m.focusLeft {
		rightPanel = activePanelBorder.Width(rightWidth)
	}

	// Paneles lado a lado
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel.Render(leftContent),
		rightPanel.Render(rightContent),
	)

	// Footer contextual
	footer := renderCatalogFooter(m)

	return lipgloss.JoinVertical(lipgloss.Left, header, panels, footer)
}

// renderCatalogHeader genera el header principal del catálogo
func renderCatalogHeader(m Model) string {
	// Logo pequeño a la izquierda
	logo := logoStyle.Render(LogoSmall)

	// Estadísticas a la derecha
	projectCount := len(m.projects)
	var featureCount int64
	if len(m.projects) > 0 && m.activeProject < len(m.projects) {
		featureCount = m.projects[m.activeProject].FeatureCount
	}

	stats := lipgloss.NewStyle().
		Foreground(colorDimWhite).
		Render(fmt.Sprintf("%d proyectos • %d features", projectCount, featureCount))

	// Keybindings en el header
	keys := lipgloss.NewStyle().
		Foreground(colorSubtle).
		Render("[/] buscar  [?] ayuda  [q] salir")

	// Primera línea: logo | stats
	leftPart := lipgloss.NewStyle().Width(m.width/2 - 2).Render(logo)
	rightPart := lipgloss.NewStyle().Width(m.width/2 - 2).Align(lipgloss.Right).Render(stats)
	firstLine := lipgloss.JoinHorizontal(lipgloss.Top, leftPart, rightPart)

	// Segunda línea: keybindings
	secondLine := lipgloss.NewStyle().
		Width(m.width - 2).
		Align(lipgloss.Right).
		Render(keys)

	// Combinar con fondo
	headerContent := lipgloss.JoinVertical(lipgloss.Left, firstLine, secondLine)

	return lipgloss.NewStyle().
		Background(colorPrimary).
		Width(m.width).
		Padding(0, 1).
		Render(headerContent)
}

// renderCatalogFooter genera el footer contextual
func renderCatalogFooter(m Model) string {
	var footerText string
	if m.confirmDelete {
		footerText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Render("⚠ Presiona 'd' de nuevo para confirmar eliminación") +
			dimStyle.Render("  [esc] cancelar")
	} else if m.statusMsg != "" {
		footerText = lipgloss.NewStyle().
			Foreground(colorReady).
			Render("✓ " + m.statusMsg)
	} else if m.focusLeft {
		footerText = renderKeyHints([]keyHint{
			{"enter/→", "ir a features"},
			{"↑↓", "navegar"},
			{"/", "buscar"},
			{"?", "ayuda"},
			{"q", "salir"},
		})
	} else {
		footerText = renderKeyHints([]keyHint{
			{"enter", "abrir"},
			{"←", "proyectos"},
			{"h", "historial"},
			{"e", "exportar"},
			{"d", "eliminar"},
			{"q", "salir"},
		})
	}

	return footerStyle.Width(m.width - 2).Render(footerText)
}

type keyHint struct {
	key  string
	desc string
}

func renderKeyHints(hints []keyHint) string {
	var parts []string
	for _, h := range hints {
		parts = append(parts, keyStyle.Render("["+h.key+"]")+dimStyle.Render(" "+h.desc))
	}
	return strings.Join(parts, "  ")
}

func renderProjectList(m Model, width int) string {
	var sb strings.Builder

	// Título de la sección con icono
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Render("📁 PROYECTOS")
	sb.WriteString(title)
	sb.WriteString("\n")
	sb.WriteString(dividerStyle.Render(strings.Repeat("─", width-4)))
	sb.WriteString("\n")

	if len(m.projects) == 0 {
		sb.WriteString("\n")
		sb.WriteString(dimStyle.Render("  No hay proyectos"))
		sb.WriteString("\n")
		sb.WriteString(dimStyle.Render("  Usa /shape para crear"))
		return sb.String()
	}

	for i, p := range m.projects {
		var line string
		if i == m.activeProject {
			// Item seleccionado con estilo destacado
			icon := "▶"
			name := lipgloss.NewStyle().
				Foreground(colorWhite).
				Bold(true).
				Render(truncate(p.Name, width-12))
			count := lipgloss.NewStyle().
				Foreground(colorAccent).
				Render(fmt.Sprintf("(%d)", p.FeatureCount))
			line = fmt.Sprintf("%s %s %s", icon, name, count)
		} else {
			// Item normal
			name := truncate(p.Name, width-10)
			count := dimStyle.Render(fmt.Sprintf("(%d)", p.FeatureCount))
			line = fmt.Sprintf("  %s %s", name, count)
		}
		sb.WriteString(line)
		if i < len(m.projects)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func renderFeatureList(m Model, width int) string {
	var sb strings.Builder

	// Título de la sección con el proyecto actual
	projectName := "todos"
	if len(m.projects) > 0 && m.activeProject < len(m.projects) {
		projectName = m.projects[m.activeProject].Name
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Render(fmt.Sprintf("✨ FEATURES — %s", projectName))
	sb.WriteString(title)
	sb.WriteString("\n")
	sb.WriteString(dividerStyle.Render(strings.Repeat("─", width-4)))
	sb.WriteString("\n")

	if len(m.features) == 0 {
		sb.WriteString("\n")
		sb.WriteString(dimStyle.Render("  No hay features"))
		sb.WriteString("\n")
		sb.WriteString(dimStyle.Render("  Usa /shape para crear"))
		return sb.String()
	}

	// Anchos de columna dinámicos
	titleWidth := width / 2
	if titleWidth < 16 {
		titleWidth = 16
	}

	for i, f := range m.features {
		var line string

		// Icono y estado
		statusIcon := StatusIcon(f.Status)
		typeIcon := TypeIcon(f.Type)

		// Título truncado
		title := truncate(f.Title, titleWidth)

		// Badges
		typ := TypeColor(f.Type).Render(f.Type)
		status := StatusColor(f.Status).Render(fmt.Sprintf("%s %s", statusIcon, f.Status))
		ver := dimStyle.Render(fmt.Sprintf("v%d", f.Version))

		if i == m.activeFeature && !m.focusLeft {
			// Item seleccionado
			selectedTitle := lipgloss.NewStyle().
				Foreground(colorWhite).
				Bold(true).
				Render(fmt.Sprintf("%s %s", typeIcon, title))
			line = fmt.Sprintf("▶ %s  %s  %s  %s", selectedTitle, typ, status, ver)
		} else {
			// Item normal
			line = fmt.Sprintf("  %s %s  %s  %s  %s", typeIcon, title, typ, status, ver)
		}

		sb.WriteString(line)
		if i < len(m.features)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
