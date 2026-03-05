package tui

import "github.com/charmbracelet/lipgloss"

// Logo ASCII art para feature-shaper
const Logo = `
  ███████╗██████╗ ██╗ █████╗ ██████╗ ██████╗ ███████╗██████╗ 
  ██╔════╝██╔══██╗██║██╔══██╗██╔══██╗██╔══██╗██╔════╝██╔══██╗
  █████╗  ██████╔╝██║███████║██████╔╝██████╔╝█████╗  ██████╔╝
  ██╔══╝  ██╔══██╗██║██╔══██║██╔═══╝ ██╔═══╝ ██╔══╝  ██╔══██╗
  ██║     ██║  ██║██║██║  ██║██║     ██║     ███████╗██║  ██║
  ╚═╝     ╚═╝  ╚═╝╚═╝╚═╝  ╚═╝╚═╝     ╚═╝     ╚══════╝╚═╝  ╚═╝
`

// Logo pequeño para header compacto
const LogoSmall = "◈ feature-shaper"

// Colores por tipo de feature
var (
	colorProduct   = lipgloss.Color("#4FC3F7") // Azul cielo
	colorTechnical = lipgloss.Color("#81C784") // Verde
	colorBusiness  = lipgloss.Color("#FFB74D") // Naranja
)

// Colores por status
var (
	colorReady      = lipgloss.Color("#69F0AE") // Verde brillante
	colorDraft      = lipgloss.Color("#90A4AE") // Gris
	colorInProgress = lipgloss.Color("#FFD54F") // Amarillo
	colorDone       = lipgloss.Color("#A5D6A7") // Verde pastel
)

// Colores de UI - Paleta más moderna
var (
	colorPrimary       = lipgloss.Color("#9C27B0") // Púrpura
	colorSecondary     = lipgloss.Color("#7C4DFF") // Púrpura claro
	colorAccent        = lipgloss.Color("#E040FB") // Rosa
	colorActiveProject = lipgloss.Color("#CE93D8") // Lila
	colorActiveBorder  = lipgloss.Color("#9C27B0") // Púrpura
	colorSubtle        = lipgloss.Color("#626262") // Gris oscuro
	colorWhite         = lipgloss.Color("#FAFAFA") // Blanco
	colorDimWhite      = lipgloss.Color("#AAAAAA") // Gris claro
	colorBackground    = lipgloss.Color("#1A1A2E") // Azul muy oscuro
	colorSurface       = lipgloss.Color("#16213E") // Azul oscuro
	colorGradientStart = lipgloss.Color("#667eea") // Gradiente inicio
	colorGradientEnd   = lipgloss.Color("#764ba2") // Gradiente fin
)

// Estilos base
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite).
			Padding(0, 1)

	logoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite).
			Background(colorPrimary).
			Padding(0, 1).
			MarginBottom(0)

	footerStyle = lipgloss.NewStyle().
			Foreground(colorDimWhite).
			Padding(0, 1).
			MarginTop(0)

	panelBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSubtle).
			Padding(0, 1)

	activePanelBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorPrimary).
				BorderBackground(lipgloss.Color("")).
				Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(colorDimWhite)

	badgeStyle = lipgloss.NewStyle().
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(colorDimWhite).
			Italic(true).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorSubtle)

	dimStyle = lipgloss.NewStyle().
			Foreground(colorSubtle)

	// Nuevos estilos para mejorar la UI
	keyStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			Padding(0, 1)

	descStyle = lipgloss.NewStyle().
			Foreground(colorDimWhite)

	sectionHeaderStyle = lipgloss.NewStyle().
				Foreground(colorWhite).
				Bold(true).
				Padding(1, 0, 0, 0)

	dividerStyle = lipgloss.NewStyle().
			Foreground(colorSubtle)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(1, 2)

	// Estilo para el input de búsqueda
	searchInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorPrimary).
				Padding(0, 1).
				Foreground(colorWhite)
)

// TypeColor devuelve el estilo coloreado según el tipo de feature
func TypeColor(typ string) lipgloss.Style {
	switch typ {
	case "product":
		return badgeStyle.Foreground(colorProduct)
	case "technical":
		return badgeStyle.Foreground(colorTechnical)
	case "business":
		return badgeStyle.Foreground(colorBusiness)
	default:
		return badgeStyle.Foreground(colorDimWhite)
	}
}

// StatusColor devuelve el estilo coloreado según el status
func StatusColor(status string) lipgloss.Style {
	switch status {
	case "ready":
		return badgeStyle.Foreground(colorReady)
	case "draft":
		return badgeStyle.Foreground(colorDraft)
	case "in-progress":
		return badgeStyle.Foreground(colorInProgress)
	case "done":
		return badgeStyle.Foreground(colorDone)
	default:
		return badgeStyle.Foreground(colorDimWhite)
	}
}

// StatusIcon devuelve el icono correspondiente al status
func StatusIcon(status string) string {
	switch status {
	case "ready":
		return "●"
	case "draft":
		return "○"
	case "in-progress":
		return "◐"
	case "done":
		return "✓"
	default:
		return "◌"
	}
}

// TypeIcon devuelve el icono correspondiente al tipo
func TypeIcon(typ string) string {
	switch typ {
	case "product":
		return "◆"
	case "technical":
		return "⚙"
	case "business":
		return "◈"
	default:
		return "◇"
	}
}

// renderHelpKey genera un par clave-descripción estilizado
func renderHelpKey(key, desc string) string {
	return keyStyle.Render(key) + descStyle.Render(" "+desc)
}
