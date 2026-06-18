package main

import "testing"

func TestBuildManifestSegmentCount(t *testing.T) {
	ladder := []Variant{
		{Codec: "h264", Resolution: "1280x720", BitrateKbps: 2350},
	}
	// 20 seconds, 4-second segments => 5 segments
	m := BuildManifest("test-title", ladder, 20000, 4000)
	key := VariantKey(ladder[0])
	if len(m.Segments[key]) != 5 {
		t.Fatalf("expected 5 segments, got %d", len(m.Segments[key]))
	}
}

func TestBuildManifestPartialLastSegment(t *testing.T) {
	ladder := []Variant{
		{Codec: "h264", Resolution: "320x240", BitrateKbps: 235},
	}
	// 10 seconds with 4-second segments: segments of 4, 4, 2 ms
	m := BuildManifest("test-partial", ladder, 10000, 4000)
	key := VariantKey(ladder[0])
	segs := m.Segments[key]
	if len(segs) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(segs))
	}
	last := segs[len(segs)-1]
	if last.DurationMs != 2000 {
		t.Fatalf("expected last segment duration 2000ms, got %d", last.DurationMs)
	}
}

func TestBuildManifestVariantCount(t *testing.T) {
	ladder := []Variant{
		{Codec: "h264", Resolution: "320x240", BitrateKbps: 235},
		{Codec: "h264", Resolution: "1280x720", BitrateKbps: 2350},
		{Codec: "h265", Resolution: "1920x1080", BitrateKbps: 4000},
	}
	m := BuildManifest("multi", ladder, 8000, 4000)
	if len(m.Variants) != 3 {
		t.Fatalf("expected 3 variants, got %d", len(m.Variants))
	}
	if len(m.Segments) != 3 {
		t.Fatalf("expected 3 segment maps, got %d", len(m.Segments))
	}
}
