package llm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TemplateConfig represents a customizable prompt template
type TemplateConfig struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Persona     string            `json:"persona"`
	ToneStyle   string            `json:"tone_style"`
	Focus       []string          `json:"focus"`
	Keywords    map[string]string `json:"keywords"`
	OutputHints []string          `json:"output_hints"`
}

// BuiltinTemplates contains pre-defined templates for common use cases
var BuiltinTemplates = map[string]TemplateConfig{
	"default": {
		Name:        "Default",
		Description: "Balanced template suitable for most tech roles",
		Persona:     "You are an expert technical writer specializing in resume optimization.",
		ToneStyle:   "professional, concise, achievement-focused",
		Focus:       []string{"technical impact", "quantifiable results", "problem-solving"},
		Keywords: map[string]string{
			"implement":  "Engineered",
			"fix":        "Resolved",
			"add":        "Developed",
			"update":     "Enhanced",
			"refactor":   "Optimized",
			"improve":    "Improved",
			"create":     "Architected",
			"build":      "Built",
		},
		OutputHints: []string{
			"Start with strong action verbs",
			"Include metrics when possible",
			"Focus on business impact",
		},
	},
	"startup": {
		Name:        "Startup",
		Description: "Fast-paced startup environment emphasis",
		Persona:     "You are a startup-focused career coach who values speed, impact, and versatility.",
		ToneStyle:   "dynamic, results-driven, entrepreneurial",
		Focus:       []string{"rapid delivery", "cross-functional impact", "innovation", "scalability"},
		Keywords: map[string]string{
			"implement":  "Shipped",
			"fix":        "Unblocked",
			"add":        "Launched",
			"update":     "Iterated",
			"refactor":   "Scaled",
			"improve":    "Accelerated",
			"create":     "Pioneered",
			"build":      "Spearheaded",
		},
		OutputHints: []string{
			"Emphasize speed to market",
			"Highlight cross-team collaboration",
			"Show ownership and initiative",
		},
	},
	"enterprise": {
		Name:        "Enterprise",
		Description: "Large corporation and enterprise focus",
		Persona:     "You are a senior technical writer experienced with Fortune 500 companies.",
		ToneStyle:   "formal, process-oriented, compliance-aware",
		Focus:       []string{"reliability", "security", "compliance", "stakeholder management"},
		Keywords: map[string]string{
			"implement":  "Implemented",
			"fix":        "Remediated",
			"add":        "Introduced",
			"update":     "Modernized",
			"refactor":   "Streamlined",
			"improve":    "Optimized",
			"create":     "Established",
			"build":      "Developed",
		},
		OutputHints: []string{
			"Reference industry standards",
			"Include compliance considerations",
			"Mention stakeholder alignment",
		},
	},
	"backend": {
		Name:        "Backend Engineer",
		Description: "Backend/infrastructure engineering focus",
		Persona:     "You are a backend engineering expert who understands distributed systems.",
		ToneStyle:   "technical, precise, systems-focused",
		Focus:       []string{"performance", "scalability", "reliability", "data integrity"},
		Keywords: map[string]string{
			"implement":  "Engineered",
			"fix":        "Debugged",
			"add":        "Integrated",
			"update":     "Upgraded",
			"refactor":   "Re-architected",
			"improve":    "Optimized",
			"create":     "Designed",
			"build":      "Constructed",
		},
		OutputHints: []string{
			"Include performance metrics (latency, throughput)",
			"Mention scale (requests/sec, data volume)",
			"Reference architectural patterns",
		},
	},
	"frontend": {
		Name:        "Frontend Engineer",
		Description: "Frontend/UI engineering focus",
		Persona:     "You are a frontend engineering expert focused on user experience.",
		ToneStyle:   "user-centric, visual, accessibility-aware",
		Focus:       []string{"user experience", "performance", "accessibility", "design systems"},
		Keywords: map[string]string{
			"implement":  "Crafted",
			"fix":        "Resolved",
			"add":        "Introduced",
			"update":     "Refined",
			"refactor":   "Modernized",
			"improve":    "Enhanced",
			"create":     "Designed",
			"build":      "Built",
		},
		OutputHints: []string{
			"Include UX metrics (load time, LCP, CLS)",
			"Mention accessibility improvements",
			"Reference user engagement impact",
		},
	},
	"devops": {
		Name:        "DevOps/SRE",
		Description: "DevOps and Site Reliability Engineering focus",
		Persona:     "You are an SRE expert focused on reliability and automation.",
		ToneStyle:   "operational, metrics-driven, automation-focused",
		Focus:       []string{"reliability", "automation", "observability", "incident response"},
		Keywords: map[string]string{
			"implement":  "Deployed",
			"fix":        "Mitigated",
			"add":        "Automated",
			"update":     "Upgraded",
			"refactor":   "Consolidated",
			"improve":    "Hardened",
			"create":     "Architected",
			"build":      "Provisioned",
		},
		OutputHints: []string{
			"Include reliability metrics (uptime, MTTR)",
			"Mention automation benefits (time saved)",
			"Reference infrastructure scale",
		},
	},
	"data": {
		Name:        "Data Engineer",
		Description: "Data engineering and analytics focus",
		Persona:     "You are a data engineering expert focused on data pipelines and analytics.",
		ToneStyle:   "analytical, data-driven, precision-focused",
		Focus:       []string{"data quality", "pipeline efficiency", "analytics enablement", "data governance"},
		Keywords: map[string]string{
			"implement":  "Built",
			"fix":        "Corrected",
			"add":        "Ingested",
			"update":     "Migrated",
			"refactor":   "Optimized",
			"improve":    "Accelerated",
			"create":     "Designed",
			"build":      "Constructed",
		},
		OutputHints: []string{
			"Include data volume metrics",
			"Mention processing time improvements",
			"Reference data quality improvements",
		},
	},
}

// TemplateManager handles loading and applying templates
type TemplateManager struct {
	templates map[string]TemplateConfig
	current   string
}

// NewTemplateManager creates a new template manager with builtin templates
func NewTemplateManager() *TemplateManager {
	tm := &TemplateManager{
		templates: make(map[string]TemplateConfig),
		current:   "default",
	}

	// Load builtin templates
	for name, tmpl := range BuiltinTemplates {
		tm.templates[name] = tmpl
	}

	return tm
}

// LoadCustomTemplate loads a custom template from a JSON file
func (tm *TemplateManager) LoadCustomTemplate(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	var config TemplateConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if config.Name == "" {
		config.Name = filepath.Base(path)
	}

	tm.templates[strings.ToLower(config.Name)] = config
	return nil
}

// SetTemplate sets the current active template
func (tm *TemplateManager) SetTemplate(name string) error {
	name = strings.ToLower(name)
	if _, ok := tm.templates[name]; !ok {
		return fmt.Errorf("template not found: %s", name)
	}
	tm.current = name
	return nil
}

// GetTemplate returns the current template
func (tm *TemplateManager) GetTemplate() TemplateConfig {
	return tm.templates[tm.current]
}

// ListTemplates returns all available template names
func (tm *TemplateManager) ListTemplates() []string {
	var names []string
	for name := range tm.templates {
		names = append(names, name)
	}
	return names
}

// BuildSystemPrompt generates a system prompt based on the current template
func (tm *TemplateManager) BuildSystemPrompt() string {
	tmpl := tm.GetTemplate()

	var sb strings.Builder

	// Persona
	sb.WriteString(tmpl.Persona)
	sb.WriteString("\n\n")

	// Tone and style
	sb.WriteString(fmt.Sprintf("Write in a %s tone.\n\n", tmpl.ToneStyle))

	// Focus areas
	sb.WriteString("Focus on these aspects:\n")
	for _, focus := range tmpl.Focus {
		sb.WriteString(fmt.Sprintf("- %s\n", focus))
	}
	sb.WriteString("\n")

	// STAR method instructions
	sb.WriteString(`Apply the STAR method (Situation, Task, Action, Result) to create impactful bullet points.

Guidelines:
- Each bullet point should be 1-2 sentences
- Start with a strong action verb
- Include quantifiable impact when possible
- Focus on business value and outcomes
`)

	// Output hints
	if len(tmpl.OutputHints) > 0 {
		sb.WriteString("\nAdditional guidelines:\n")
		for _, hint := range tmpl.OutputHints {
			sb.WriteString(fmt.Sprintf("- %s\n", hint))
		}
	}

	// Output format
	sb.WriteString(`
Output Format:
Return a JSON array with objects containing:
- "hash": first 7 characters of commit hash
- "category": one of "Feature", "Fix", "Refactor", "Test", "Docs", "Chore"
- "summary": the resume bullet point

Example output:
[
  {"hash": "abc1234", "category": "Feature", "summary": "Engineered real-time notification system serving 10K concurrent users"}
]
`)

	return sb.String()
}

// TransformVerb transforms a common verb using the template's keyword mapping
func (tm *TemplateManager) TransformVerb(verb string) string {
	tmpl := tm.GetTemplate()
	verb = strings.ToLower(verb)
	if transformed, ok := tmpl.Keywords[verb]; ok {
		return transformed
	}
	return verb
}
