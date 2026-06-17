package templaterender

import (
	"fmt"
	"strings"
	"time"

	"github.com/woragis/management/backend/server/internal/models"
)

type Vars struct {
	Number            string
	Track             string
	ProblemTitle      string
	Difficulty        string
	LeetcodeURL       string
	SubmissionURL     string
	ShortDescription  string
	YoutubeURL        string
	Date              string
	ProblemList       string
	NextTheme         string
	InviteLink        string
}

func FromVideo(v *models.LeetcodeVideo, settings *models.LeetcodeChannelSettings) Vars {
	title := v.Title
	if v.ProblemTitle != nil && strings.TrimSpace(*v.ProblemTitle) != "" {
		title = strings.TrimSpace(*v.ProblemTitle)
	}
	num := ""
	if v.SeriesNumber != nil {
		num = fmt.Sprintf("%d", *v.SeriesNumber)
	} else if v.LeetcodeProblemNumber != nil {
		num = fmt.Sprintf("%d", *v.LeetcodeProblemNumber)
	}
	track := ""
	if v.TrackName != nil {
		track = strings.TrimSpace(*v.TrackName)
	}
	diff := ""
	if v.Difficulty != nil {
		diff = formatDifficulty(*v.Difficulty)
	}
	url := leetcodeURL(v)
	sub := ""
	if v.LeetcodeSubmissionURL != nil {
		sub = strings.TrimSpace(*v.LeetcodeSubmissionURL)
	}
	short := ""
	if v.ShortDescription != nil {
		short = strings.TrimSpace(*v.ShortDescription)
	}
	yt := ""
	if v.YoutubeURL != nil {
		yt = strings.TrimSpace(*v.YoutubeURL)
	}
	date := ""
	if v.ProblemDate != nil {
		date = v.ProblemDate.Format("02/01/2006")
	}
	invite := ""
	nextTheme := ""
	if settings != nil {
		if settings.InviteLink != nil {
			invite = strings.TrimSpace(*settings.InviteLink)
		}
		if settings.NextTheme != nil {
			nextTheme = strings.TrimSpace(*settings.NextTheme)
		}
	}
	return Vars{
		Number:           num,
		Track:            track,
		ProblemTitle:     title,
		Difficulty:       diff,
		LeetcodeURL:      url,
		SubmissionURL:    sub,
		ShortDescription: short,
		YoutubeURL:       yt,
		Date:             date,
		InviteLink:       invite,
		NextTheme:        nextTheme,
	}
}

func Render(body string, vars Vars) string {
	repl := map[string]string{
		"{{number}}":           vars.Number,
		"{{track}}":            vars.Track,
		"{{problemTitle}}":     vars.ProblemTitle,
		"{{difficulty}}":       vars.Difficulty,
		"{{leetcodeUrl}}":      vars.LeetcodeURL,
		"{{submissionUrl}}":    vars.SubmissionURL,
		"{{shortDescription}}": vars.ShortDescription,
		"{{youtubeUrl}}":       vars.YoutubeURL,
		"{{date}}":             vars.Date,
		"{{problemList}}":      vars.ProblemList,
		"{{nextTheme}}":        vars.NextTheme,
		"{{inviteLink}}":       vars.InviteLink,
	}
	out := body
	for k, v := range repl {
		out = strings.ReplaceAll(out, k, v)
	}
	return out
}

func leetcodeURL(v *models.LeetcodeVideo) string {
	if v.LeetcodeProblemURL != nil && strings.TrimSpace(*v.LeetcodeProblemURL) != "" {
		return strings.TrimSpace(*v.LeetcodeProblemURL)
	}
	if v.LeetcodeSubmissionURL != nil && strings.TrimSpace(*v.LeetcodeSubmissionURL) != "" {
		return strings.TrimSpace(*v.LeetcodeSubmissionURL)
	}
	if v.LeetcodeSlug != nil && strings.TrimSpace(*v.LeetcodeSlug) != "" {
		return "https://leetcode.com/problems/" + strings.TrimSpace(*v.LeetcodeSlug) + "/"
	}
	return ""
}

func formatDifficulty(d string) string {
	switch strings.ToLower(strings.TrimSpace(d)) {
	case "easy":
		return "Easy"
	case "medium":
		return "Medium"
	case "hard":
		return "Hard"
	default:
		return strings.TrimSpace(d)
	}
}

func FormatProblemList(videos []models.LeetcodeVideo) string {
	var lines []string
	for _, v := range videos {
		title := v.Title
		if v.ProblemTitle != nil && strings.TrimSpace(*v.ProblemTitle) != "" {
			title = strings.TrimSpace(*v.ProblemTitle)
		}
		prefix := "•"
		if v.SeriesNumber != nil {
			prefix = fmt.Sprintf("✅ #%d", *v.SeriesNumber)
		}
		lines = append(lines, fmt.Sprintf("%s %s", prefix, title))
	}
	return strings.Join(lines, "\n")
}

func ParseDateInTZ(dateStr, tz string) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}
	t, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func TodayInTZ(tz string) time.Time {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
}
