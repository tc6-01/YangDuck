package generator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

type ToolInfo struct {
	Name        string
	Description string
	Version     string
	Homepage    string
	Examples    []string
	ReadmeParts []string
	EnvVars     []string
}

func CollectBrewInfo(pkg string) (*ToolInfo, error) {
	out, err := exec.Command("brew", "info", "--json=v2", pkg).Output()
	if err != nil {
		return nil, fmt.Errorf("brew info %s: %w", pkg, err)
	}
	var result struct {
		Formulae []struct {
			Name     string `json:"name"`
			Desc     string `json:"desc"`
			Homepage string `json:"homepage"`
			Versions struct {
				Stable string `json:"stable"`
			} `json:"versions"`
		} `json:"formulae"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, err
	}
	if len(result.Formulae) == 0 {
		return nil, fmt.Errorf("no formula found for %s", pkg)
	}
	f := result.Formulae[0]
	info := &ToolInfo{
		Name:        f.Name,
		Description: f.Desc,
		Version:     f.Versions.Stable,
		Homepage:    f.Homepage,
	}

	if tldr, err := exec.Command("tldr", pkg).Output(); err == nil {
		info.Examples = parseTLDRExamples(string(tldr))
	}

	return info, nil
}

func CollectNPMInfo(pkg string) (*ToolInfo, error) {
	out, err := exec.Command("npm", "info", pkg, "--json").Output()
	if err != nil {
		return nil, fmt.Errorf("npm info %s: %w", pkg, err)
	}
	var result struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Version     string `json:"version"`
		Homepage    string `json:"homepage"`
		Repository  struct {
			URL string `json:"url"`
		} `json:"repository"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, err
	}
	info := &ToolInfo{
		Name:        result.Name,
		Description: result.Description,
		Version:     result.Version,
		Homepage:    result.Homepage,
	}

	repoURL := extractGitHubURL(result.Repository.URL)
	if repoURL != "" {
		if readme, err := fetchGitHubREADME(repoURL); err == nil {
			info.ReadmeParts = []string{readme}
			info.EnvVars = extractEnvVars(readme)
		}
	}

	return info, nil
}

func parseTLDRExamples(tldr string) []string {
	var examples []string
	lines := strings.Split(tldr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			examples = append(examples, strings.TrimPrefix(line, "- "))
		}
	}
	return examples
}

func extractGitHubURL(repoURL string) string {
	repoURL = strings.TrimPrefix(repoURL, "git+")
	repoURL = strings.TrimSuffix(repoURL, ".git")
	repoURL = strings.Replace(repoURL, "git://", "https://", 1)
	repoURL = strings.Replace(repoURL, "ssh://git@", "https://", 1)
	if strings.Contains(repoURL, "github.com") {
		return repoURL
	}
	return ""
}

func fetchGitHubREADME(repoURL string) (string, error) {
	parts := strings.Split(strings.TrimPrefix(repoURL, "https://github.com/"), "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid github url")
	}
	apiURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/README.md", parts[0], parts[1])
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}
	buf := make([]byte, 8192)
	n, _ := resp.Body.Read(buf)
	return string(buf[:n]), nil
}

func extractEnvVars(text string) []string {
	var vars []string
	seen := make(map[string]bool)
	for _, line := range strings.Split(text, "\n") {
		for _, word := range strings.Fields(line) {
			if len(word) > 3 && word == strings.ToUpper(word) && strings.Contains(word, "_") {
				clean := strings.Trim(word, "`\"',;:()[]{}$")
				if !seen[clean] && len(clean) > 3 {
					vars = append(vars, clean)
					seen[clean] = true
				}
			}
		}
	}
	return vars
}
