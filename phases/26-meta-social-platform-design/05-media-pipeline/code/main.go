package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// MediaType distinguishes photos from videos.
type MediaType string

const (
	MediaPhoto MediaType = "PHOTO"
	MediaVideo MediaType = "VIDEO"
)

// JobStatus represents processing state of a media job.
type JobStatus string

const (
	StatusPending    JobStatus = "PENDING"
	StatusProcessing JobStatus = "PROCESSING"
	StatusCompleted  JobStatus = "COMPLETED"
	StatusFailed     JobStatus = "FAILED"
)

// MediaJob represents a single upload that must be processed.
type MediaJob struct {
	ID       string    `json:"id"`
	Type     MediaType `json:"type"`
	Status   JobStatus `json:"status"`
	Variants []string  `json:"variants,omitempty"`
}

// photoVariants returns the standard set of processed sizes for a photo.
func photoVariants() []string {
	return []string{"thumbnail_150x150", "medium_600x600", "large_1200x1200", "original"}
}

// videoVariants returns all transcoded format+resolution combinations for a video.
func videoVariants() []string {
	codecs := []string{"h264", "vp9", "av1"}
	resolutions := []string{"360p", "720p", "1080p"}
	variants := make([]string, 0, len(codecs)*len(resolutions))
	for _, codec := range codecs {
		for _, res := range resolutions {
			variants = append(variants, fmt.Sprintf("%s_%s", codec, res))
		}
	}
	// AV1 also gets 4K
	variants = append(variants, "av1_4k")
	return variants
}

// ProcessJob simulates transcoding/resizing a media job into all required variants.
// It updates the job in place, setting Status to COMPLETED and filling Variants.
func ProcessJob(job *MediaJob) {
	job.Status = StatusProcessing

	// Simulate processing time proportional to type.
	switch job.Type {
	case MediaPhoto:
		time.Sleep(2 * time.Millisecond) // fast: resize only
		job.Variants = photoVariants()
	case MediaVideo:
		time.Sleep(10 * time.Millisecond) // slower: transcode multiple variants
		job.Variants = videoVariants()
	default:
		job.Status = StatusFailed
		return
	}

	job.Status = StatusCompleted
}

// StorageTier returns the storage tier for a media object based on its age in days.
// hot:  0–30 days   (SSD, high IOPS)
// warm: 31–180 days (HDD)
// cold: 181+ days   (erasure-coded, off-site)
func StorageTier(ageInDays int) string {
	switch {
	case ageInDays <= 30:
		return "hot"
	case ageInDays <= 180:
		return "warm"
	default:
		return "cold"
	}
}

// BatchResult summarises the outcome of processing a batch of media jobs.
type BatchResult struct {
	Total     int      `json:"total"`
	Completed int      `json:"completed"`
	Failed    int      `json:"failed"`
	Jobs      []MediaJob `json:"jobs"`
}

func main() {
	jobs := []*MediaJob{
		{ID: "photo-001", Type: MediaPhoto, Status: StatusPending},
		{ID: "photo-002", Type: MediaPhoto, Status: StatusPending},
		{ID: "video-001", Type: MediaVideo, Status: StatusPending},
		{ID: "reel-001", Type: MediaVideo, Status: StatusPending},
		{ID: "photo-003", Type: MediaPhoto, Status: StatusPending},
	}

	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(j *MediaJob) {
			defer wg.Done()
			ProcessJob(j)
		}(job)
	}
	wg.Wait()

	result := BatchResult{Total: len(jobs)}
	for _, j := range jobs {
		if j.Status == StatusCompleted {
			result.Completed++
		} else {
			result.Failed++
		}
		result.Jobs = append(result.Jobs, *j)
	}

	// Demonstrate storage tier logic.
	tierExamples := []struct {
		Age  int
		Tier string
	}{
		{Age: 1, Tier: StorageTier(1)},
		{Age: 30, Tier: StorageTier(30)},
		{Age: 31, Tier: StorageTier(31)},
		{Age: 180, Tier: StorageTier(180)},
		{Age: 181, Tier: StorageTier(181)},
		{Age: 365, Tier: StorageTier(365)},
	}

	output := map[string]any{
		"batch_result":   result,
		"tier_examples":  tierExamples,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
