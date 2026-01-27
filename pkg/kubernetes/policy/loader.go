package policy

import (
	_ "embed"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

//go:embed dangerous_permissions.yaml
var defaultPermissionsYAML []byte

var (
	// DefaultPermissions 默认的高危权限列表（从嵌入的 YAML 加载）
	DefaultPermissions []DangerousPermission
)

func init() {
	var err error
	DefaultPermissions, err = LoadPermissionsFromBytes(defaultPermissionsYAML)
	if err != nil {
		panic(fmt.Sprintf("failed to load default dangerous permissions: %v", err))
	}
}

// LoadPermissionsFromBytes 从字节数组加载权限配置
func LoadPermissionsFromBytes(data []byte) ([]DangerousPermission, error) {
	var config DangerousPermissionsConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return config.Permissions, nil
}

// LoadPermissionsFromFile 从文件加载权限配置
func LoadPermissionsFromFile(path string) ([]DangerousPermission, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return LoadPermissionsFromBytes(data)
}

// MergePermissions 合并多个权限列表（用于自定义扩展）
func MergePermissions(base, custom []DangerousPermission) []DangerousPermission {
	result := make([]DangerousPermission, len(base))
	copy(result, base)
	return append(result, custom...)
}

// FilterByLevel 按级别过滤权限
func FilterByLevel(permissions []DangerousPermission, minLevel Level) []DangerousPermission {
	levelOrder := map[Level]int{
		LevelCritical: 3,
		LevelHigh:     2,
		LevelMedium:   1,
	}

	minOrder := levelOrder[minLevel]
	var result []DangerousPermission
	for _, p := range permissions {
		if levelOrder[p.Level] >= minOrder {
			result = append(result, p)
		}
	}
	return result
}

// GroupResultsByLevel 按级别分组结果
func GroupResultsByLevel(results []CheckResult) map[Level][]CheckResult {
	grouped := make(map[Level][]CheckResult)
	for _, r := range results {
		grouped[r.Permission.Level] = append(grouped[r.Permission.Level], r)
	}
	return grouped
}
