package templaterender

import (
	"encoding/json"
	"regexp"
	"strings"

	"gorm.io/datatypes"
)

var placeholderRe = regexp.MustCompile(`\{\{\s*([\w.]+)\s*\}\}`)

// DataSource tells the renderer which entity instance to load for a job.
type DataSource struct {
	Program     string `json:"program"`
	Resolver    string `json:"resolver"`
	Date        string `json:"date"`
	ProjectID   string `json:"projectId"`
	ProjectSlug string `json:"projectSlug"`
}

func ParseDataSource(raw datatypes.JSON) DataSource {
	if len(raw) == 0 {
		return DataSource{}
	}
	var ds DataSource
	_ = json.Unmarshal(raw, &ds)
	return ds
}

// Bindings map template placeholder (without braces) to field ref, e.g. "leetcode.problemTitle".
type Bindings map[string]string

func ParseBindings(raw datatypes.JSON) Bindings {
	if len(raw) == 0 {
		return Bindings{}
	}
	var b Bindings
	if err := json.Unmarshal(raw, &b); err != nil {
		return Bindings{}
	}
	return b
}

func RenderBody(body string, vars map[string]string) string {
	return placeholderRe.ReplaceAllStringFunc(body, func(match string) string {
		sub := placeholderRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return ""
		}
		key := sub[1]
		if v, ok := vars[key]; ok {
			return v
		}
		// support dotted keys stored flat
		if v, ok := vars[strings.ToLower(key)]; ok {
			return v
		}
		return ""
	})
}

func PlaceholdersInBody(body string) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range placeholderRe.FindAllStringSubmatch(body, -1) {
		if len(m) < 2 || seen[m[1]] {
			continue
		}
		seen[m[1]] = true
		out = append(out, m[1])
	}
	return out
}

// CatalogFields returns known binding targets per program slug.
func CatalogFields(program string) []CatalogField {
	switch strings.TrimSpace(strings.ToLower(program)) {
	case "leetcode":
		return leetcodeCatalog
	case "project":
		return projectCatalog
	default:
		return nil
	}
}

type CatalogField struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Binding     string `json:"binding"`
	Description string `json:"description,omitempty"`
}

var leetcodeCatalog = []CatalogField{
	{Key: "number", Label: "Number", Binding: "leetcode.number"},
	{Key: "track", Label: "Track", Binding: "leetcode.track"},
	{Key: "problemTitle", Label: "Problem title", Binding: "leetcode.problemTitle"},
	{Key: "difficulty", Label: "Difficulty", Binding: "leetcode.difficulty"},
	{Key: "leetcodeUrl", Label: "LeetCode URL", Binding: "leetcode.leetcodeUrl"},
	{Key: "submissionUrl", Label: "Submission URL", Binding: "leetcode.submissionUrl"},
	{Key: "shortDescription", Label: "Short description", Binding: "leetcode.shortDescription"},
	{Key: "youtubeUrl", Label: "YouTube URL", Binding: "leetcode.youtubeUrl"},
	{Key: "date", Label: "Date", Binding: "leetcode.date"},
	{Key: "problemList", Label: "Problem list", Binding: "leetcode.problemList"},
	{Key: "nextTheme", Label: "Next theme", Binding: "leetcode.nextTheme"},
	{Key: "inviteLink", Label: "Invite link", Binding: "leetcode.inviteLink"},
}

var projectCatalog = []CatalogField{
	{Key: "projectName", Label: "Name", Binding: "project.name"},
	{Key: "projectSlug", Label: "Slug", Binding: "project.slug"},
	{Key: "shortDescription", Label: "Short description", Binding: "project.shortDescription"},
	{Key: "demoUrl", Label: "Demo URL", Binding: "project.demoUrl"},
	{Key: "githubUrl", Label: "GitHub URL", Binding: "project.githubUrl"},
	{Key: "repoUrl", Label: "Repo URL", Binding: "project.repoUrl"},
	{Key: "stack", Label: "Stack", Binding: "project.stack"},
}

func DefaultBindings(program string) Bindings {
	fields := CatalogFields(program)
	out := Bindings{}
	for _, f := range fields {
		out[f.Key] = f.Binding
	}
	return out
}

func MergeBindings(explicit Bindings, program string, body string) Bindings {
	out := DefaultBindings(program)
	for k, v := range explicit {
		out[k] = v
	}
	for _, ph := range PlaceholdersInBody(body) {
		if _, ok := out[ph]; !ok && len(CatalogFields(program)) > 0 {
			out[ph] = program + "." + ph
		}
	}
	return out
}
