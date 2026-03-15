package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sipeed/picoclaw/pkg/skills"
)

type SearchAllTool struct {
	registry *ToolRegistry
	loader   *skills.SkillsLoader
	unlocked map[string]map[string]time.Time // chatID -> toolName -> expiry
	mu       sync.RWMutex
}

func NewSearchAllTool(r *ToolRegistry, l *skills.SkillsLoader) *SearchAllTool {
	return &SearchAllTool{
		registry: r,
		loader:   l,
		unlocked: make(map[string]map[string]time.Time),
	}
}

func (t *SearchAllTool) IsVisible(chatID, toolName string) bool {
	if toolName == t.Name() {
		return true
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	expiry, ok := t.unlocked[chatID][toolName]
	return ok && time.Now().Before(expiry)
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

	var basicTools []string
	var results []string

	// Search Registry Tools
	t.registry.mu.RLock()
	for name, entry := range t.registry.tools {
		// Skip searching the search tool itself to avoid noise
		if name == t.Name() {
			continue
		}
		desc := strings.ToLower(entry.Tool.Description())
		if query == "" || strings.Contains(strings.ToLower(name), query) || strings.Contains(desc, query) {
			basicTools = append(basicTools, name)
			results = append(results, fmt.Sprintf("[Basic Tool] %s: %s", name, entry.Tool.Description()))
		}
	}
	t.registry.mu.RUnlock()

	// Search Skills if loader is available
	if t.loader != nil {
		allSkills := t.loader.ListSkills()
		for _, s := range allSkills {
			name := strings.ToLower(s.Name)
			desc := strings.ToLower(s.Description)
			if query == "" || strings.Contains(name, query) || strings.Contains(desc, query) {
				results = append(results, fmt.Sprintf("[Skill] %s: %s (Path: %s)\n  To use this skill, read its SKILL.md for instructions.", s.Name, s.Description, s.Path))
			}
		}
	}

	// Record unlocked tools for this session
	chatID := ToolChatID(ctx)
	if chatID == "" {
		chatID = "default"
	}

	t.mu.Lock()
	if t.unlocked[chatID] == nil {
		t.unlocked[chatID] = make(map[string]time.Time)
	}
	expiry := time.Now().Add(5 * time.Minute) // Visible for 5 minutes
	for _, name := range basicTools {
		t.unlocked[chatID][name] = expiry
	}
	t.mu.Unlock()

	if len(results) == 0 {
		return NewToolResult("No matching tools or skills found.")
	}

	foundMsg := fmt.Sprintf("Found %d matching tools/skills. These tools have been unlocked and are now ready for use in your next turn.\n\n%s", len(results), strings.Join(results, "\n"))
	return NewToolResult(foundMsg)
}
