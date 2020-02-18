package swap

import (
	"path/filepath"
)

// Terraform specific paths
const (
	// TerraformProviderPrefix is prepended to terraform plugin artifacts registered
	// by a provider.
	TerraformProviderPrefix = "terraform-provider-"

	// TerraformWorkspaceDirectory is where terraform related artifacts are stored.
	TerraformWorkspaceDirectory = ".terraform"
)

var (
	// TerraformPluginPath is where plugins are stored, by platform/arch.
	TerraformPluginPath = filepath.Join(TerraformWorkspaceDirectory, "plugins")
)

// GetDefaultTerraformPlatformPath attempts to find a default platform based to modify plugin
// state in.
func GetDefaultTerraformPlatformPath() string {
	dirs, _ := getDirectories(TerraformPluginPath)
	if len(dirs) > 0 {
		return dirs[0]
	}
	return ""
}

// GetTerraformProviderArtifactPath attempts to find a Terraform provider's artifact.
func GetTerraformProviderArtifactPath(provider, pluginsPath string) string {
	pattern := filepath.Join(pluginsPath, TerraformProviderPrefix+provider+"*")
	matches, _ := filepath.Glob(pattern)
	if len(matches) == 1 {
		return matches[0]
	}
	return ""
}
