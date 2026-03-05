package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/agustinalbonico/feature-shaper/internal/store"
	"github.com/mark3labs/mcp-go/mcp"
)

func FeatureSave(featureStore *store.FeatureStore) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		_ = ctx

		projectSlug := mcp.ParseString(request, "projectSlug", "")
		slug := mcp.ParseString(request, "slug", "")
		title := mcp.ParseString(request, "title", "")
		typ := mcp.ParseString(request, "type", "")
		content := mcp.ParseString(request, "content", "")
		status := mcp.ParseString(request, "status", "draft")
		changelog := mcp.ParseString(request, "changelog", "")

		if projectSlug == "" || slug == "" || title == "" || typ == "" || content == "" {
			return mcp.NewToolResultError("projectSlug, slug, title, type and content are required"), nil
		}

		feature, err := featureStore.Save(projectSlug, slug, title, typ, content, status, changelog)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		payload, err := json.Marshal(feature)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cannot serialize feature: %v", err)), nil
		}

		return mcp.NewToolResultText(string(payload)), nil
	}
}

func FeatureGet(featureStore *store.FeatureStore) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		_ = ctx

		slug := mcp.ParseString(request, "slug", "")
		projectSlug := mcp.ParseString(request, "projectSlug", "")
		if slug == "" {
			return mcp.NewToolResultError("slug is required"), nil
		}

		feature, err := featureStore.Get(slug, projectSlug)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		payload, err := json.Marshal(feature)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cannot serialize feature: %v", err)), nil
		}

		return mcp.NewToolResultText(string(payload)), nil
	}
}

func FeatureSearch(featureStore *store.FeatureStore) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		_ = ctx

		query := mcp.ParseString(request, "query", "")
		projectSlug := mcp.ParseString(request, "projectSlug", "")
		if query == "" {
			return mcp.NewToolResultError("query is required"), nil
		}

		results, err := featureStore.Search(query, projectSlug)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		payload, err := json.Marshal(results)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cannot serialize search results: %v", err)), nil
		}

		return mcp.NewToolResultText(string(payload)), nil
	}
}

func FeatureCatalog(featureStore *store.FeatureStore) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		_ = ctx

		projectSlug := mcp.ParseString(request, "projectSlug", "")
		status := mcp.ParseString(request, "status", "")
		typ := mcp.ParseString(request, "type", "")
		if projectSlug == "" {
			return mcp.NewToolResultError("projectSlug is required"), nil
		}

		features, err := featureStore.Catalog(projectSlug, status, typ)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		payload, err := json.Marshal(features)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cannot serialize feature catalog: %v", err)), nil
		}

		return mcp.NewToolResultText(string(payload)), nil
	}
}

func FeatureVersions(featureStore *store.FeatureStore) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		_ = ctx

		slug := mcp.ParseString(request, "slug", "")
		projectSlug := mcp.ParseString(request, "projectSlug", "")
		if slug == "" || projectSlug == "" {
			return mcp.NewToolResultError("slug and projectSlug are required"), nil
		}

		feature, err := featureStore.Get(slug, projectSlug)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		versions, err := featureStore.Versions(slug, projectSlug)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		response := struct {
			FeatureID int64 `json:"featureId"`
			Versions  []struct {
				ID        int64  `json:"id"`
				FeatureID int64  `json:"featureId"`
				Version   int    `json:"version"`
				Content   string `json:"content"`
				Changelog string `json:"changelog"`
				CreatedAt string `json:"createdAt"`
			} `json:"versions"`
		}{
			FeatureID: feature.ID,
			Versions: make([]struct {
				ID        int64  `json:"id"`
				FeatureID int64  `json:"featureId"`
				Version   int    `json:"version"`
				Content   string `json:"content"`
				Changelog string `json:"changelog"`
				CreatedAt string `json:"createdAt"`
			}, 0, len(versions)),
		}

		for _, version := range versions {
			response.Versions = append(response.Versions, struct {
				ID        int64  `json:"id"`
				FeatureID int64  `json:"featureId"`
				Version   int    `json:"version"`
				Content   string `json:"content"`
				Changelog string `json:"changelog"`
				CreatedAt string `json:"createdAt"`
			}{
				ID:        version.ID,
				FeatureID: version.FeatureID,
				Version:   version.Version,
				Content:   version.Content,
				Changelog: version.Changelog,
				CreatedAt: version.CreatedAt,
			})
		}

		payload, err := json.Marshal(response)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cannot serialize versions: %v", err)), nil
		}

		return mcp.NewToolResultText(string(payload)), nil
	}
}

func FeatureGetVersion(featureStore *store.FeatureStore) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		_ = ctx

		featureID := int64(mcp.ParseInt(request, "featureId", 0))
		version := mcp.ParseInt(request, "version", 0)
		if featureID <= 0 || version <= 0 {
			return mcp.NewToolResultError("featureId and version are required"), nil
		}

		featureVersion, err := featureStore.GetVersion(featureID, version)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		payload, err := json.Marshal(featureVersion)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cannot serialize feature version: %v", err)), nil
		}

		return mcp.NewToolResultText(string(payload)), nil
	}
}

func ProjectRegister(projectStore *store.ProjectStore) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		_ = ctx

		slug := mcp.ParseString(request, "slug", "")
		name := mcp.ParseString(request, "name", "")
		path := mcp.ParseString(request, "path", "")
		if slug == "" || name == "" {
			return mcp.NewToolResultError("slug and name are required"), nil
		}

		if err := projectStore.Register(slug, name, path); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		payload, err := json.Marshal(map[string]string{"status": "ok"})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cannot serialize project register response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(payload)), nil
	}
}

func ProjectList(projectStore *store.ProjectStore) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		_ = ctx
		_ = request

		projects, err := projectStore.List()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		payload, err := json.Marshal(projects)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("cannot serialize projects: %v", err)), nil
		}

		return mcp.NewToolResultText(string(payload)), nil
	}
}
