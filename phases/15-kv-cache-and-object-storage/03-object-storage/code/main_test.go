package main

import "testing"

func TestValidateObjectStoragePolicyAcceptsStrongDefaults(t *testing.T) {
	cfg := ObjectStoragePolicy{
		Name:                 "strong",
		MultipartEnabled:     true,
		ChecksumRequired:     true,
		VersioningEnabled:    true,
		MetadataFinalize:     true,
		RetentionAwareDelete: true,
		OrphanSweeper:        true,
		StorageClass:         "archive",
		ReplicationRegions:   2,
	}
	if issues := ValidateObjectStoragePolicy(cfg); len(issues) != 0 {
		t.Fatalf("ValidateObjectStoragePolicy returned issues: %v", issues)
	}
}

func TestValidateObjectStoragePolicyRejectsUnsafeWorkflow(t *testing.T) {
	cfg := ObjectStoragePolicy{
		Name:               "unsafe",
		MultipartEnabled:   false,
		ChecksumRequired:   false,
		MetadataFinalize:   false,
		StorageClass:       "coldish",
		ReplicationRegions: 0,
	}
	if issues := ValidateObjectStoragePolicy(cfg); len(issues) < 5 {
		t.Fatalf("ValidateObjectStoragePolicy returned too few issues: %v", issues)
	}
}
