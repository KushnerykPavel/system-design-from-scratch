package main

import (
	"encoding/json"
	"flag"
	"os"
)

type ObjectStoragePolicy struct {
	Name                 string `json:"name"`
	MultipartEnabled     bool   `json:"multipart_enabled"`
	ChecksumRequired     bool   `json:"checksum_required"`
	VersioningEnabled    bool   `json:"versioning_enabled"`
	MetadataFinalize     bool   `json:"metadata_finalize"`
	RetentionAwareDelete bool   `json:"retention_aware_delete"`
	OrphanSweeper        bool   `json:"orphan_sweeper"`
	StorageClass         string `json:"storage_class"`
	ReplicationRegions   int    `json:"replication_regions"`
}

func ValidateObjectStoragePolicy(cfg ObjectStoragePolicy) []string {
	var issues []string
	if !cfg.MultipartEnabled {
		issues = append(issues, "multipart_enabled should be true for large or unstable-client uploads")
	}
	if !cfg.ChecksumRequired {
		issues = append(issues, "checksum_required should be true for integrity guarantees")
	}
	if !cfg.MetadataFinalize {
		issues = append(issues, "metadata_finalize should be explicit to avoid phantom objects")
	}
	if !cfg.RetentionAwareDelete {
		issues = append(issues, "retention_aware_delete should be enabled for safe deletion semantics")
	}
	if !cfg.OrphanSweeper {
		issues = append(issues, "orphan_sweeper should be enabled to reconcile failed uploads")
	}
	if cfg.StorageClass != "hot" && cfg.StorageClass != "standard" && cfg.StorageClass != "archive" {
		issues = append(issues, "storage_class must be hot, standard, or archive")
	}
	if cfg.ReplicationRegions < 1 {
		issues = append(issues, "replication_regions must be at least 1")
	}
	return issues
}

func main() {
	name := flag.String("name", "user-object-platform", "name of the object storage policy")
	flag.Parse()

	cfg := ObjectStoragePolicy{
		Name:                 *name,
		MultipartEnabled:     true,
		ChecksumRequired:     true,
		VersioningEnabled:    true,
		MetadataFinalize:     true,
		RetentionAwareDelete: true,
		OrphanSweeper:        true,
		StorageClass:         "standard",
		ReplicationRegions:   2,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"policy": cfg,
		"issues": ValidateObjectStoragePolicy(cfg),
	})
}
