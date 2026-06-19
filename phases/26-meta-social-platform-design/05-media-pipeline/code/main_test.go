package main

import (
	"strings"
	"testing"
)

func TestProcessJobPhoto(t *testing.T) {
	job := &MediaJob{ID: "test-photo", Type: MediaPhoto, Status: StatusPending}
	ProcessJob(job)

	if job.Status != StatusCompleted {
		t.Fatalf("expected COMPLETED, got %s", job.Status)
	}
	if len(job.Variants) == 0 {
		t.Fatal("expected variants to be populated for photo job")
	}
	// Photo must include an original variant.
	found := false
	for _, v := range job.Variants {
		if v == "original" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected 'original' variant in photo job variants")
	}
}

func TestProcessJobVideo(t *testing.T) {
	job := &MediaJob{ID: "test-video", Type: MediaVideo, Status: StatusPending}
	ProcessJob(job)

	if job.Status != StatusCompleted {
		t.Fatalf("expected COMPLETED, got %s", job.Status)
	}
	// Video must produce variants for at least 3 codecs.
	codecs := map[string]bool{"h264": false, "vp9": false, "av1": false}
	for _, v := range job.Variants {
		for codec := range codecs {
			if strings.HasPrefix(v, codec) {
				codecs[codec] = true
			}
		}
	}
	for codec, seen := range codecs {
		if !seen {
			t.Errorf("expected variant for codec %s but not found in %v", codec, job.Variants)
		}
	}
}

func TestProcessJobVideoVariantCount(t *testing.T) {
	job := &MediaJob{ID: "test-video-count", Type: MediaVideo, Status: StatusPending}
	ProcessJob(job)

	// 3 codecs × 3 resolutions + 1 AV1 4K = 10
	if len(job.Variants) != 10 {
		t.Fatalf("expected 10 video variants, got %d: %v", len(job.Variants), job.Variants)
	}
}

func TestProcessJobUnknownType(t *testing.T) {
	job := &MediaJob{ID: "test-unknown", Type: "AUDIO", Status: StatusPending}
	ProcessJob(job)

	if job.Status != StatusFailed {
		t.Fatalf("expected FAILED for unknown type, got %s", job.Status)
	}
}

func TestStorageTierBoundaries(t *testing.T) {
	cases := []struct {
		age  int
		want string
	}{
		{0, "hot"},
		{1, "hot"},
		{30, "hot"},
		{31, "warm"},
		{180, "warm"},
		{181, "cold"},
		{365, "cold"},
		{3650, "cold"},
	}
	for _, tc := range cases {
		got := StorageTier(tc.age)
		if got != tc.want {
			t.Errorf("StorageTier(%d) = %q, want %q", tc.age, got, tc.want)
		}
	}
}

func TestPhotoVariantsContainThumbnail(t *testing.T) {
	variants := photoVariants()
	found := false
	for _, v := range variants {
		if strings.HasPrefix(v, "thumbnail") {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("photoVariants should include a thumbnail variant")
	}
}

func TestVideoVariantsContain4K(t *testing.T) {
	variants := videoVariants()
	found := false
	for _, v := range variants {
		if strings.Contains(v, "4k") {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("videoVariants should include a 4K variant")
	}
}

func TestProcessJobStatusTransition(t *testing.T) {
	job := &MediaJob{ID: "status-check", Type: MediaPhoto, Status: StatusPending}
	// Before processing the status must be PENDING.
	if job.Status != StatusPending {
		t.Fatalf("expected initial status PENDING, got %s", job.Status)
	}
	ProcessJob(job)
	if job.Status != StatusCompleted {
		t.Fatalf("expected final status COMPLETED, got %s", job.Status)
	}
}
