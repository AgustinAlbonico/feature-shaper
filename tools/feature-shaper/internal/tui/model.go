package tui

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/agustinalbonico/feature-shaper/internal/db"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// View representa la vista activa de la TUI
type View int

const (
	CatalogView View = iota
	DetailView
	HistoryView
	VersionDetailView
	SearchView
)

// Model es el modelo principal de Bubble Tea (Elm Architecture)
type Model struct {
	database      *sql.DB
	view          View
	projects      []db.ProjectWithCount
	features      []db.Feature
	versions      []db.FeatureVersion
	searchResults []db.FeatureSearchResult
	activeProject int
	activeFeature int
	activeVersion int
	focusLeft     bool
	searchInput   textinput.Model
	viewport      viewport.Model
	width         int
	height        int
	ready         bool
	showHelp      bool
	confirmDelete bool
	statusMsg     string
	err           error
}

// Mensajes personalizados para carga de datos
type projectsLoadedMsg struct{ projects []db.ProjectWithCount }
type featuresLoadedMsg struct{ features []db.Feature }
type versionsLoadedMsg struct{ versions []db.FeatureVersion }
type searchResultsMsg struct{ results []db.FeatureSearchResult }
type featureDeletedMsg struct{ slug string }
type errMsg struct{ err error }

// NewModel crea un nuevo modelo TUI
func NewModel(database *sql.DB) Model {
	ti := textinput.New()
	ti.Placeholder = "Buscar features..."
	ti.CharLimit = 100

	vp := viewport.New(80, 20)

	return Model{
		database:    database,
		view:        CatalogView,
		focusLeft:   true,
		searchInput: ti,
		viewport:    vp,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadProjects()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 6
		m.ready = true
		return m, nil

	case projectsLoadedMsg:
		m.projects = msg.projects
		m.activeProject = 0
		if len(m.projects) > 0 {
			// Si hay un solo proyecto, arrancar con foco en features
			if len(m.projects) == 1 {
				m.focusLeft = false
			}
			return m, m.loadFeatures(m.projects[0].Slug)
		}
		m.features = nil
		return m, nil

	case featuresLoadedMsg:
		m.features = msg.features
		m.activeFeature = 0
		return m, nil

	case versionsLoadedMsg:
		m.versions = msg.versions
		m.activeVersion = 0
		return m, nil

	case searchResultsMsg:
		m.searchResults = msg.results
		return m, nil

	case featureDeletedMsg:
		m.statusMsg = fmt.Sprintf("Feature '%s' eliminada", msg.slug)
		m.confirmDelete = false
		if len(m.projects) > 0 {
			return m, tea.Batch(m.loadProjects(), m.loadFeatures(m.projects[m.activeProject].Slug))
		}
		return m, m.loadProjects()

	case errMsg:
		m.err = msg.err
		m.statusMsg = fmt.Sprintf("Error: %v", msg.err)
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Actualizar el viewport si estamos en vista de detalle
	if m.view == DetailView || m.view == VersionDetailView {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	// Actualizar el input de búsqueda
	if m.view == SearchView {
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Salir siempre funciona
	if key == "ctrl+c" {
		return m, tea.Quit
	}

	// En vista de búsqueda, el input captura teclas
	if m.view == SearchView {
		return m.handleSearchKey(msg)
	}

	// Keybindings globales
	switch key {
	case "q":
		return m, tea.Quit
	case "?":
		m.showHelp = !m.showHelp
		return m, nil
	case "/":
		m.view = SearchView
		m.searchInput.Focus()
		m.searchResults = nil
		return m, textinput.Blink
	}

	// Keybindings por vista
	switch m.view {
	case CatalogView:
		return m.handleCatalogKey(msg)
	case DetailView:
		return m.handleDetailKey(msg)
	case HistoryView:
		return m.handleHistoryKey(msg)
	case VersionDetailView:
		return m.handleVersionDetailKey(msg)
	}

	return m, nil
}

func (m Model) handleCatalogKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "tab", "left", "right":
		m.focusLeft = !m.focusLeft
		m.confirmDelete = false
		return m, nil

	case "up", "k":
		m.confirmDelete = false
		if m.focusLeft {
			if m.activeProject > 0 {
				m.activeProject--
				return m, m.loadFeatures(m.projects[m.activeProject].Slug)
			}
		} else {
			if m.activeFeature > 0 {
				m.activeFeature--
			}
		}
		return m, nil

	case "down", "j":
		m.confirmDelete = false
		if m.focusLeft {
			if m.activeProject < len(m.projects)-1 {
				m.activeProject++
				return m, m.loadFeatures(m.projects[m.activeProject].Slug)
			}
		} else {
			if m.activeFeature < len(m.features)-1 {
				m.activeFeature++
			}
		}
		return m, nil

	case "enter":
		if m.focusLeft {
			// Enter en panel izquierdo: seleccionar proyecto y mover foco a features
			if len(m.projects) > 0 {
				m.focusLeft = false
				return m, m.loadFeatures(m.projects[m.activeProject].Slug)
			}
		} else if len(m.features) > 0 {
			feature := m.features[m.activeFeature]
			m.view = DetailView
			m.viewport.SetContent(feature.Content)
			m.viewport.GotoTop()
		}
		return m, nil

	case "h":
		if !m.focusLeft && len(m.features) > 0 {
			m.view = HistoryView
			return m, m.loadVersions(m.features[m.activeFeature].ID)
		}
		return m, nil

	case "e":
		if !m.focusLeft && len(m.features) > 0 {
			return m, m.exportFeature(m.features[m.activeFeature])
		}
		return m, nil

	case "d":
		if !m.focusLeft && len(m.features) > 0 {
			if m.confirmDelete {
				feature := m.features[m.activeFeature]
				return m, m.deleteFeature(feature.ID, feature.Slug)
			}
			m.confirmDelete = true
			m.statusMsg = "Presiona 'd' de nuevo para confirmar eliminación"
			return m, nil
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleDetailKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "b", "esc":
		m.view = CatalogView
		return m, nil
	case "h":
		if len(m.features) > 0 {
			m.view = HistoryView
			return m, m.loadVersions(m.features[m.activeFeature].ID)
		}
		return m, nil
	case "e":
		if len(m.features) > 0 {
			return m, m.exportFeature(m.features[m.activeFeature])
		}
		return m, nil
	case "up", "k":
		m.viewport.LineUp(1)
		return m, nil
	case "down", "j":
		m.viewport.LineDown(1)
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) handleHistoryKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "b", "esc":
		m.view = CatalogView
		return m, nil
	case "up", "k":
		if m.activeVersion > 0 {
			m.activeVersion--
		}
		return m, nil
	case "down", "j":
		if m.activeVersion < len(m.versions)-1 {
			m.activeVersion++
		}
		return m, nil
	case "enter":
		if len(m.versions) > 0 {
			version := m.versions[m.activeVersion]
			m.view = VersionDetailView
			m.viewport.SetContent(version.Content)
			m.viewport.GotoTop()
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleVersionDetailKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "b", "esc":
		m.view = HistoryView
		return m, nil
	case "up", "k":
		m.viewport.LineUp(1)
		return m, nil
	case "down", "j":
		m.viewport.LineDown(1)
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "esc":
		m.view = CatalogView
		m.searchInput.Blur()
		m.searchInput.Reset()
		return m, nil
	case "enter":
		query := strings.TrimSpace(m.searchInput.Value())
		if query != "" {
			return m, m.searchFeatures(query)
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if !m.ready {
		return "Cargando..."
	}

	if m.showHelp {
		return m.renderHelp()
	}

	var content string
	switch m.view {
	case CatalogView:
		content = renderCatalog(m)
	case DetailView:
		content = renderDetail(m)
	case HistoryView:
		content = renderHistory(m)
	case VersionDetailView:
		content = renderVersionDetail(m)
	case SearchView:
		content = renderSearch(m)
	}

	return content
}

func (m Model) renderHelp() string {
	// Logo centrado
	logoRendered := lipgloss.NewStyle().
		Foreground(colorPrimary).
		Render(LogoSmall)

	// Título
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Render("Ayuda y Atajos de Teclado")

	// Secciones de ayuda
	navSection := renderHelpSection("Navegación", []keyHint{
		{"↑/↓ ó j/k", "Mover cursor"},
		{"Tab / ←/→", "Cambiar foco entre paneles"},
		{"Enter", "Abrir feature / ver versión"},
		{"b / Esc", "Volver a la vista anterior"},
	})

	actionSection := renderHelpSection("Acciones", []keyHint{
		{" /", "Buscar features (FTS)"},
		{"h", "Ver historial de versiones"},
		{"e", "Exportar feature a .md"},
		{"d", "Eliminar feature (doble d para confirmar)"},
	})

	generalSection := renderHelpSection("General", []keyHint{
		{"?", "Mostrar/ocultar esta ayuda"},
		{"q", "Salir"},
	})

	// Leyenda de tipos
	typeLegend := renderTypeLegend()

	// Leyenda de estados
	statusLegend := renderStatusLegend()

	// Footer
	footer := lipgloss.NewStyle().
		Foreground(colorSubtle).
		Render("Presiona ? para volver al catálogo")

	// Construir el contenido
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		logoRendered,
		"",
		title,
		"",
		dividerStyle.Render(strings.Repeat("─", m.width-8)),
		"",
		navSection,
		"",
		actionSection,
		"",
		generalSection,
		"",
		dividerStyle.Render(strings.Repeat("─", m.width-8)),
		"",
		typeLegend,
		"",
		statusLegend,
		"",
		footer,
	)

	// Envolver en un box con borde
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(1, 2).
		Width(m.width - 4).
		Height(m.height - 2).
		Render(content)
}

// renderHelpSection genera una sección de ayuda con título y key hints
func renderHelpSection(title string, hints []keyHint) string {
	var sb strings.Builder

	sb.WriteString(sectionHeaderStyle.Render(title))
	sb.WriteString("\n")

	for _, h := range hints {
		sb.WriteString(renderHelpKey(h.key, h.desc))
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderTypeLegend genera la leyenda de tipos de feature
func renderTypeLegend() string {
	var sb strings.Builder

	sb.WriteString(sectionHeaderStyle.Render("Tipos de Feature"))
	sb.WriteString("\n")

	types := []struct {
		name  string
		color lipgloss.Color
		icon  string
	}{
		{"product", colorProduct, "◆"},
		{"technical", colorTechnical, "⚙"},
		{"business", colorBusiness, "◈"},
	}

	for i, t := range types {
		badge := lipgloss.NewStyle().
			Foreground(t.color).
			Render(fmt.Sprintf("%s %s", t.icon, t.name))
		sb.WriteString(fmt.Sprintf("  %s", badge))
		if i < len(types)-1 {
			sb.WriteString("   ")
		}
	}

	return sb.String()
}

// renderStatusLegend genera la leyenda de estados
func renderStatusLegend() string {
	var sb strings.Builder

	sb.WriteString(sectionHeaderStyle.Render("Estados"))
	sb.WriteString("\n")

	statuses := []struct {
		name  string
		color lipgloss.Color
		icon  string
	}{
		{"draft", colorDraft, "○"},
		{"in-progress", colorInProgress, "◐"},
		{"ready", colorReady, "●"},
		{"done", colorDone, "✓"},
	}

	for i, s := range statuses {
		badge := lipgloss.NewStyle().
			Foreground(s.color).
			Render(fmt.Sprintf("%s %s", s.icon, s.name))
		sb.WriteString(fmt.Sprintf("  %s", badge))
		if i < len(statuses)-1 {
			sb.WriteString("   ")
		}
	}

	return sb.String()
}

// Comandos para carga de datos
func (m Model) loadProjects() tea.Cmd {
	return func() tea.Msg {
		projects, err := db.ListProjects(m.database)
		if err != nil {
			return errMsg{err}
		}
		return projectsLoadedMsg{projects}
	}
}

func (m Model) loadFeatures(projectSlug string) tea.Cmd {
	return func() tea.Msg {
		features, err := db.ListFeatures(m.database, projectSlug, "", "")
		if err != nil {
			return errMsg{err}
		}
		return featuresLoadedMsg{features}
	}
}

func (m Model) loadVersions(featureID int64) tea.Cmd {
	return func() tea.Msg {
		versions, err := db.ListFeatureVersions(m.database, featureID)
		if err != nil {
			return errMsg{err}
		}
		return versionsLoadedMsg{versions}
	}
}

func (m Model) searchFeatures(query string) tea.Cmd {
	return func() tea.Msg {
		projectSlug := ""
		if len(m.projects) > 0 {
			projectSlug = m.projects[m.activeProject].Slug
		}
		results, err := db.SearchFeatures(m.database, query, projectSlug)
		if err != nil {
			return errMsg{err}
		}
		return searchResultsMsg{results}
	}
}

func (m Model) deleteFeature(featureID int64, slug string) tea.Cmd {
	return func() tea.Msg {
		_, err := m.database.Exec(`DELETE FROM features WHERE id = ?`, featureID)
		if err != nil {
			return errMsg{fmt.Errorf("no se pudo eliminar la feature: %w", err)}
		}
		return featureDeletedMsg{slug}
	}
}

func (m Model) exportFeature(feature db.Feature) tea.Cmd {
	return func() tea.Msg {
		dir := "docs/features"
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return errMsg{fmt.Errorf("no se pudo crear directorio: %w", err)}
		}

		path := filepath.Join(dir, feature.Slug+".md")
		if err := os.WriteFile(path, []byte(feature.Content), 0o644); err != nil {
			return errMsg{fmt.Errorf("no se pudo escribir archivo: %w", err)}
		}

		m.statusMsg = fmt.Sprintf("Exportado a %s", path)
		return errMsg{err: nil}
	}
}
