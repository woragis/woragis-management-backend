package templaterender

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/models"
)

func TestRenderReplacesVariables(t *testing.T) {
	body := "Hello {{number}} {{problemTitle}} ({{difficulty}}) {{leetcodeUrl}}"
	out := Render(body, Vars{
		Number:       "42",
		ProblemTitle: "Two Sum",
		Difficulty:   "Easy",
		LeetcodeURL:  "https://leetcode.com/problems/two-sum/",
	})
	want := "Hello 42 Two Sum (Easy) https://leetcode.com/problems/two-sum/"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestFromVideoPrefersProblemTitleAndSeriesNumber(t *testing.T) {
	series := 7
	problemNum := 1
	problemTitle := "Two Sum"
	track := "Arrays"
	diff := "easy"
	url := "https://leetcode.com/problems/two-sum/"
	day := time.Date(2026, 6, 16, 0, 0, 0, 0, time.UTC)

	v := &models.LeetcodeVideo{
		ID:                    uuid.New(),
		Title:                 "fallback title",
		SeriesNumber:          &series,
		TrackName:             &track,
		ProblemTitle:          &problemTitle,
		LeetcodeProblemNumber: &problemNum,
		Difficulty:            &diff,
		LeetcodeProblemURL:    &url,
		ProblemDate:           &day,
	}
	invite := "https://chat.whatsapp.com/test"
	settings := &models.LeetcodeChannelSettings{InviteLink: &invite}

	vars := FromVideo(v, settings)
	if vars.Number != "7" {
		t.Fatalf("number: got %q", vars.Number)
	}
	if vars.ProblemTitle != "Two Sum" {
		t.Fatalf("problemTitle: got %q", vars.ProblemTitle)
	}
	if vars.Difficulty != "Easy" {
		t.Fatalf("difficulty: got %q", vars.Difficulty)
	}
	if vars.LeetcodeURL != url {
		t.Fatalf("leetcode url: got %q", vars.LeetcodeURL)
	}
	if vars.Date != "16/06/2026" {
		t.Fatalf("date: got %q", vars.Date)
	}
	if vars.InviteLink != invite {
		t.Fatalf("invite: got %q", vars.InviteLink)
	}
}

func TestLeetcodeURLFallsBackToSlug(t *testing.T) {
	slug := "two-sum"
	v := &models.LeetcodeVideo{LeetcodeSlug: &slug}
	if got := leetcodeURL(v); got != "https://leetcode.com/problems/two-sum/" {
		t.Fatalf("got %q", got)
	}
}

func TestFormatProblemList(t *testing.T) {
	series := 3
	title := "Linked List"
	v := models.LeetcodeVideo{SeriesNumber: &series, ProblemTitle: &title}
	out := FormatProblemList([]models.LeetcodeVideo{v})
	if !strings.Contains(out, "✅ #3 Linked List") {
		t.Fatalf("unexpected list: %q", out)
	}
}

func TestParseDateInTZ(t *testing.T) {
	day, err := ParseDateInTZ("2026-06-16", "America/Sao_Paulo")
	if err != nil {
		t.Fatal(err)
	}
	if day.Format("2006-01-02") != "2026-06-16" {
		t.Fatalf("got %s", day.Format("2006-01-02"))
	}
}
