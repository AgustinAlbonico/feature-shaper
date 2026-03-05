package tui

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
)

// Start inicia la TUI con la base de datos proporcionada
func Start(database *sql.DB) error {
	m := NewModel(database)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
