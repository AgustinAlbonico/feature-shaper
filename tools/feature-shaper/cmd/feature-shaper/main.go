package main

import (
	"fmt"
	"os"

	"github.com/agustinalbonico/feature-shaper/internal/db"
	featuremcp "github.com/agustinalbonico/feature-shaper/internal/mcp"
	"github.com/agustinalbonico/feature-shaper/internal/store"
	"github.com/agustinalbonico/feature-shaper/internal/tui"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]
	var err error

	switch subcommand {
	case "migrate":
		err = runMigrate()
	case "mcp":
		err = runMCP()
	case "tui":
		err = runTUI()
	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runMigrate() error {
	database, err := db.Migrate()
	if err != nil {
		return err
	}
	defer database.Close()

	fmt.Println("database migrated successfully")
	return nil
}

func runMCP() error {
	database, err := db.Migrate()
	if err != nil {
		return err
	}
	defer database.Close()

	featureStore := store.NewFeatureStore(database)
	projectStore := store.NewProjectStore(database)
	mcpServer := featuremcp.NewServer(featureStore, projectStore)

	featuremcp.ServeStdio(mcpServer)
	return nil
}

func runTUI() error {
	database, err := db.Migrate()
	if err != nil {
		return err
	}
	defer database.Close()
	return tui.Start(database)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: feature-shaper <mcp|tui|migrate>")
}
