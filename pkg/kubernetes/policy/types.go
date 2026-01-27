package policy

// Level 危险级别
type Level string

const (
	LevelCritical Level = "critical"
	LevelHigh     Level = "high"
	LevelMedium   Level = "medium"
)

// DangerousPermission 定义单个高危权限
type DangerousPermission struct {
	Resource    string   `yaml:"resource"`              // 资源，如 "pods", "secrets", "nodes"
	Subresource string   `yaml:"subresource,omitempty"` // 子资源，如 "exec", "proxy"
	Group       string   `yaml:"group,omitempty"`       // API Group，默认 "" (core)
	Verbs       []string `yaml:"verbs"`                 // 动词，如 ["get"], ["create"], ["*"]
	Level       Level    `yaml:"level"`                 // 危险级别
	Description string   `yaml:"description"`           // 风险描述
	Reference   string   `yaml:"reference,omitempty"`   // 参考链接
}

// DangerousPermissionsConfig YAML 配置结构
type DangerousPermissionsConfig struct {
	Permissions []DangerousPermission `yaml:"permissions"`
}

// FullResource 返回完整资源名（包含子资源）
func (p *DangerousPermission) FullResource() string {
	if p.Subresource != "" {
		return p.Resource + "/" + p.Subresource
	}
	return p.Resource
}

// CheckResult 检测结果
type CheckResult struct {
	Permission  DangerousPermission
	Allowed     bool
	Namespace   string // "" 表示集群级别
	MatchedVerb string // 实际匹配的动词
}
