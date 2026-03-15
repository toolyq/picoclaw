package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/sipeed/picoclaw/pkg/skills"
)

type SearchAllTool struct {
	registry *ToolRegistry
	loader   *skills.SkillsLoader
}

func NewSearchAllTool(r *ToolRegistry, l *skills.SkillsLoader) *SearchAllTool {
	return &SearchAllTool{
		registry: r,
		loader:   l,
	}
}

func (t *SearchAllTool) Name() string {
	return "search_tools_and_skills"
}

func (t *SearchAllTool) Description() string {
	return "Search all available local tools and skills by query. This search includes both built-in agent tools and custom skills registered in the workspace."
}

func (t *SearchAllTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"query": map[string]any{
				"type":        "string",
				"description": "Keywords to search for in tool/skill name or description.",
			},
		},
	}
}

func (t *SearchAllTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	query, _ := args["query"].(string)
	query = strings.ToLower(strings.TrimSpace(query))

	var results []string

	// Search Registry Tools
	t.registry.mu.RLock()
	for name, entry := range t.registry.tools {
		desc := strings.ToLower(entry.Tool.Description())
		if query == "" || strings.Contains(strings.ToLower(name), query) || strings.Contains(desc, query) {
			results = append(results, fmt.Sprintf("[Basic Tool] %s: %s", name, entry.Tool.Description()))
		}
	}
	t.registry.mu.RUnlock()

	// Search Skills
	if t.loader != nil {
		allSkills := t.loader.ListSkills()
		for _, s := range allSkills {
			desc := strings.ToLower(s.Description)
			if query == "" || strings.Contains(strings.ToLower(s.Name), query) || strings.Contains(desc, query) {
				results = append(results, fmt.Sprintf("[Custom Skill] %s: %s (Location: %s)", s.Name, s.Description, s.Path))
			}
		}
	}

	if len(results) == 0 {
		return NewToolResult("No matching tools or skills found.")
	}

	return NewToolResult(fmt.Sprintf("Found %d tools/skills:\n\n%s", len(results), strings.Join(results, "\n")))
}
