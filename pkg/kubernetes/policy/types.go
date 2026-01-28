package policy

// Level risk level
type Level string

const (
	LevelCritical Level = "critical"
	LevelHigh     Level = "high"
	LevelMedium   Level = "medium"
)

// DangerousPermission defines a single high-risk permission
type DangerousPermission struct {
	Resource    string   `yaml:"resource"`              // Resource, e.g. "pods", "secrets", "nodes"
	Subresource string   `yaml:"subresource,omitempty"` // Subresource, e.g. "exec", "proxy"
	Group       string   `yaml:"group,omitempty"`       // API Group, default "" (core)
	Verbs       []string `yaml:"verbs"`                 // Verbs, e.g. ["get"], ["create"], ["*"]
	Level       Level    `yaml:"level"`                 // Risk level
	Description string   `yaml:"description"`           // Risk description
	Reference   string   `yaml:"reference,omitempty"`   // Reference link
}

// DangerousPermissionsConfig YAML configuration structure
type DangerousPermissionsConfig struct {
	Permissions []DangerousPermission `yaml:"permissions"`
}

// FullResource returns full resource name (including subresource)
func (p *DangerousPermission) FullResource() string {
	if p.Subresource != "" {
		return p.Resource + "/" + p.Subresource
	}
	return p.Resource
}

// CheckResult detection result
type CheckResult struct {
	Permission  DangerousPermission
	Allowed     bool
	Namespace   string // "" indicates cluster level
	MatchedVerb string // Actual matched verb
}
