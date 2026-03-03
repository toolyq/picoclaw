package tools

import (
	"context"
	"fmt"
	"github.com/sipeed/picoclaw/pkg/skills"
)

// ListBuiltinSkillsTool allows the LLM agent to list available local skills
// without having them all in the system prompt.
type ListBuiltinSkillsTool struct {
	skillsLoader *skills.SkillsLoader
}

// NewListBuiltinSkillsTool creates a new ListBuiltinSkillsTool.
func NewListBuiltinSkillsTool(skillsLoader *skills.SkillsLoader) *ListBuiltinSkillsTool {
	return &ListBuiltinSkillsTool{
		skillsLoader: skillsLoader,
	}
}

func (t *ListBuiltinSkillsTool) Name() string {
	return "list_builtin_skills"
}

func (t *ListBuiltinSkillsTool) Description() string {
	return "List all available internal/local skills with their descriptions and file paths. Use this to discover what local capabilities are available before reading the skill's SKILL.md for usage details."
}

func (t *ListBuiltinSkillsTool) Parameters() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
}

func (t *ListBuiltinSkillsTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	summary := t.skillsLoader.BuildSkillsSummary()
	if summary == "" {
		return SilentResult("No local skills found.")
	}

	return SilentResult(fmt.Sprintf("Available local skills:\n\n%s", summary))
}
