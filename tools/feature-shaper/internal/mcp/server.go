package mcp

import (
	"github.com/agustinalbonico/feature-shaper/internal/store"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewServer(featureStore *store.FeatureStore, projectStore *store.ProjectStore) *server.MCPServer {
	mcpServer := server.NewMCPServer("feature-shaper", "1.0.0", server.WithToolCapabilities(false))

	mcpServer.AddTool(mcp.NewTool("feature_save",
		mcp.WithDescription("Create or update a feature definition and store a version snapshot."),
		mcp.WithString("projectSlug", mcp.Required(), mcp.Description("Project slug")),
		mcp.WithString("slug", mcp.Required(), mcp.Description("Feature slug")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Feature title")),
		mcp.WithString("type", mcp.Required(), mcp.Description("Feature type")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Feature markdown content")),
		mcp.WithString("status", mcp.Description("Feature status")),
		mcp.WithString("changelog", mcp.Description("Change summary for this version")),
	), FeatureSave(featureStore))

	mcpServer.AddTool(mcp.NewTool("feature_get",
		mcp.WithDescription("Get a feature by slug, optionally scoped by project."),
		mcp.WithString("slug", mcp.Required(), mcp.Description("Feature slug")),
		mcp.WithString("projectSlug", mcp.Description("Project slug")),
	), FeatureGet(featureStore))

	mcpServer.AddTool(mcp.NewTool("feature_search",
		mcp.WithDescription("Search feature definitions using full-text search."),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
		mcp.WithString("projectSlug", mcp.Description("Project slug")),
	), FeatureSearch(featureStore))

	mcpServer.AddTool(mcp.NewTool("feature_catalog",
		mcp.WithDescription("List features for a project with optional filters."),
		mcp.WithString("projectSlug", mcp.Required(), mcp.Description("Project slug")),
		mcp.WithString("status", mcp.Description("Feature status filter")),
		mcp.WithString("type", mcp.Description("Feature type filter")),
	), FeatureCatalog(featureStore))

	mcpServer.AddTool(mcp.NewTool("feature_versions",
		mcp.WithDescription("List stored versions for a feature."),
		mcp.WithString("slug", mcp.Required(), mcp.Description("Feature slug")),
		mcp.WithString("projectSlug", mcp.Required(), mcp.Description("Project slug")),
	), FeatureVersions(featureStore))

	mcpServer.AddTool(mcp.NewTool("feature_get_version",
		mcp.WithDescription("Get a specific historical version of a feature."),
		mcp.WithNumber("featureId", mcp.Required(), mcp.Description("Feature numeric id")),
		mcp.WithNumber("version", mcp.Required(), mcp.Description("Version number")),
	), FeatureGetVersion(featureStore))

	mcpServer.AddTool(mcp.NewTool("project_register",
		mcp.WithDescription("Register or update a project in the store."),
		mcp.WithString("slug", mcp.Required(), mcp.Description("Project slug")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Project name")),
		mcp.WithString("path", mcp.Description("Project root path")),
	), ProjectRegister(projectStore))

	mcpServer.AddTool(mcp.NewTool("project_list",
		mcp.WithDescription("List all projects and feature counts."),
	), ProjectList(projectStore))

	return mcpServer
}

func ServeStdio(mcpServer *server.MCPServer) {
	server.ServeStdio(mcpServer)
}
