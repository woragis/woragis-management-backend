package service

import (
	"encoding/json"
	"strings"

	"gorm.io/datatypes"
)

func normalizeIntent(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "commercial", "academic", "personal_tool", "portfolio", "hobby", "nonprofit":
		return strings.TrimSpace(strings.ToLower(v))
	default:
		return "portfolio"
	}
}

func normalizeMonetization(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "subscription", "one_time", "ads", "services", "indirect", "none":
		return strings.TrimSpace(strings.ToLower(v))
	default:
		return "none"
	}
}

func normalizeMaturity(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "idea", "building", "mvp", "launched", "maintenance", "sunset":
		return strings.TrimSpace(strings.ToLower(v))
	default:
		return "idea"
	}
}

func normalizeVisibilityGoal(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "revenue", "job_hunting", "academic_credit", "community", "private":
		return strings.TrimSpace(strings.ToLower(v))
	default:
		return "private"
	}
}

var validDistribution = map[string]struct{}{
	"web":           {},
	"play_store":    {},
	"app_store":     {},
	"desktop":       {},
	"internal_only": {},
}

func normalizeDistribution(items []string) []string {
	if len(items) == 0 {
		return []string{}
	}
	seen := make(map[string]struct{})
	out := make([]string, 0, len(items))
	for _, raw := range items {
		v := strings.TrimSpace(strings.ToLower(raw))
		if v == "" {
			continue
		}
		if _, ok := validDistribution[v]; !ok {
			continue
		}
		if _, dup := seen[v]; dup {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func distributionJSON(items []string) datatypes.JSON {
	normalized := normalizeDistribution(items)
	if normalized == nil {
		normalized = []string{}
	}
	b, _ := json.Marshal(normalized)
	return datatypes.JSON(b)
}

func parseDistribution(raw datatypes.JSON) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil {
		return []string{}
	}
	return normalizeDistribution(out)
}
