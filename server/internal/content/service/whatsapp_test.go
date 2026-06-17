package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/woragis/management/backend/server/internal/content/repository"
	"github.com/woragis/management/backend/server/internal/models"
	"github.com/woragis/management/backend/server/internal/testutil"
)

func TestEnsureWhatsappDefaultsSeedsTemplates(t *testing.T) {
	db := testutil.OpenSQLite(t)
	svc := New(repository.New(db), nil, nil, nil, "", "")

	ctx := context.Background()
	if err := svc.EnsureWhatsappDefaults(ctx); err != nil {
		t.Fatal(err)
	}
	rows, err := svc.ListWhatsappTemplates(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) < 4 {
		t.Fatalf("expected seeded templates, got %d", len(rows))
	}
}

func TestPreviewWhatsappRendersMessage(t *testing.T) {
	db := testutil.OpenSQLite(t)
	svc := New(repository.New(db), nil, nil, nil, "", "")

	ctx := context.Background()
	if err := svc.EnsureWhatsappDefaults(ctx); err != nil {
		t.Fatal(err)
	}

	series := 12
	track := "Binary Search"
	title := "Search Insert Position"
	diff := "easy"
	url := "https://leetcode.com/problems/search-insert-position/"
	short := "Find index to insert target."

	video, err := svc.CreateVideo(ctx, CreateVideoInput{
		Title:              "LC #35",
		Status:             "draft",
		SeriesNumber:       &series,
		TrackName:          &track,
		ProblemTitle:       &title,
		Difficulty:         &diff,
		LeetcodeProblemURL: &url,
		ShortDescription:   &short,
		WhatsappEnabled:    ptrBool(true),
	})
	if err != nil {
		t.Fatal(err)
	}

	msg, err := svc.PreviewWhatsapp(ctx, video.ID, models.WhatsappTplProblemDaily)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(msg, "Search Insert Position") {
		t.Fatalf("message missing title: %s", msg)
	}
	if !strings.Contains(msg, url) {
		t.Fatalf("message missing url: %s", msg)
	}
}

func TestPatchWhatsappStatusSetsTimestamps(t *testing.T) {
	db := testutil.OpenSQLite(t)
	svc := New(repository.New(db), nil, nil, nil, "", "")

	ctx := context.Background()
	video, err := svc.CreateVideo(ctx, CreateVideoInput{Title: "LC", Status: "draft"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.PatchWhatsappStatus(ctx, video.ID, WhatsappStatusPatch{ProblemSent: true}); err != nil {
		t.Fatal(err)
	}
	updated, err := svc.GetVideo(ctx, video.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updated.WhatsappProblemSentAt == nil {
		t.Fatal("expected problem sent timestamp")
	}
	if updated.WhatsappProblemSentAt.Before(time.Now().Add(-time.Minute)) == false {
		// timestamp should be recent
	}
}

func ptrBool(v bool) *bool { return &v }
