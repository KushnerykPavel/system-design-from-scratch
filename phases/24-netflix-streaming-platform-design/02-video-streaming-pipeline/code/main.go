package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Variant represents one encoded stream profile.
type Variant struct {
	Codec      string `json:"codec"`
	Resolution string `json:"resolution"`
	BitrateKbps int   `json:"bitrate_kbps"`
}

// Segment is a single chunk in a variant stream.
type Segment struct {
	Index      int     `json:"index"`
	DurationMs int     `json:"duration_ms"`
	S3Key      string  `json:"s3_key"`
}

// Manifest holds the master playlist for a title.
type Manifest struct {
	TitleID  string    `json:"title_id"`
	Variants []Variant `json:"variants"`
	Segments map[string][]Segment `json:"segments"` // keyed by variant key
}

// VariantKey returns a stable string key for a variant.
func VariantKey(v Variant) string {
	return fmt.Sprintf("%s/%s/%d", v.Codec, v.Resolution, v.BitrateKbps)
}

// BuildManifest constructs a manifest for the given title and encoding ladder.
// titleDurationMs is the total video duration. segmentDurationMs is the target chunk size.
func BuildManifest(titleID string, ladder []Variant, titleDurationMs, segmentDurationMs int) Manifest {
	m := Manifest{
		TitleID:  titleID,
		Variants: ladder,
		Segments: make(map[string][]Segment),
	}
	for _, v := range ladder {
		key := VariantKey(v)
		idx := 0
		remaining := titleDurationMs
		for remaining > 0 {
			dur := segmentDurationMs
			if remaining < dur {
				dur = remaining
			}
			m.Segments[key] = append(m.Segments[key], Segment{
				Index:      idx,
				DurationMs: dur,
				S3Key:      fmt.Sprintf("titles/%s/%s/seg-%04d.ts", titleID, key, idx),
			})
			remaining -= dur
			idx++
		}
	}
	return m
}

func main() {
	ladder := []Variant{
		{Codec: "h264", Resolution: "320x240", BitrateKbps: 235},
		{Codec: "h264", Resolution: "1280x720", BitrateKbps: 2350},
		{Codec: "h264", Resolution: "1920x1080", BitrateKbps: 5800},
	}
	// 10-minute title, 4-second segments
	m := BuildManifest("tt-001", ladder, 10*60*1000, 4000)

	summary := map[string]any{
		"title_id":        m.TitleID,
		"variant_count":   len(m.Variants),
		"segments_per_variant": len(m.Segments[VariantKey(ladder[0])]),
		"total_segments":  len(m.Variants) * len(m.Segments[VariantKey(ladder[0])]),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(summary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
