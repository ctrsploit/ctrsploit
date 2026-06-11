package image

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/log"
)

// LayerReport holds the result of a shared-layer analysis across images.
type LayerReport struct {
	// Images lists the input image references in order.
	Images []string
	// Shared maps each shared diff_id hex to the set of images containing it.
	Shared map[string][]string
}

// Analyze resolves diff_ids for every image ref and reports which
// layers are shared among them.
func Analyze(refs []string) (*LayerReport, error) {
	// Deduplicate input references.
	seen := make(map[string]bool, len(refs))
	unique := make([]string, 0, len(refs))
	for _, ref := range refs {
		if seen[ref] {
			log.Logger.Warnf("duplicate reference %q ignored", ref)
			continue
		}
		seen[ref] = true
		unique = append(unique, ref)
	}

	if len(unique) < 2 {
		return nil, fmt.Errorf("need at least two unique image references to compare, got %d", len(unique))
	}

	// hex → set of imageRefs
	index := make(map[string]map[string]struct{})

	for _, ref := range unique {
		diffIDs, err := ResolveDiffIDs(ref)
		if err != nil {
			return nil, fmt.Errorf("resolve %q: %w", ref, err)
		}

		for _, h := range diffIDs {
			hex := h.Hex
			if index[hex] == nil {
				index[hex] = make(map[string]struct{})
			}
			index[hex][ref] = struct{}{}
		}
	}

	report := &LayerReport{
		Images: unique,
		Shared: make(map[string][]string),
	}
	for hex, owners := range index {
		if len(owners) > 1 {
			names := make([]string, 0, len(owners))
			for ref := range owners {
				names = append(names, ref)
			}
			sort.Strings(names)
			report.Shared[hex] = names
		}
	}

	return report, nil
}

// Print outputs the analysis report in a human-readable table.
func (r *LayerReport) Print() {
	if r == nil || len(r.Images) == 0 {
		return
	}

	log.Logger.Info("")
	log.Logger.Info("=== Image Layer Analysis ===")
	log.Logger.Infof("Compared images: %s", strings.Join(r.Images, ", "))
	log.Logger.Infof("")

	if len(r.Shared) == 0 {
		log.Logger.Warn("No shared layers found between the given images.")
		return
	}

	log.Logger.Warnf("Found %d shared layer(s):", len(r.Shared))
	log.Logger.Warnf("")

	// Sort shared layers by the number of images (most shared first)
	type entry struct {
		hex    string
		images []string
	}
	entries := make([]entry, 0, len(r.Shared))
	for hex, imgs := range r.Shared {
		entries = append(entries, entry{hex, imgs})
	}
	sort.Slice(entries, func(i, j int) bool {
		if len(entries[i].images) != len(entries[j].images) {
			return len(entries[i].images) > len(entries[j].images)
		}
		return entries[i].hex < entries[j].hex
	})

	for _, e := range entries {
		log.Logger.Warnf("  sha256:%s", e.hex)
		log.Logger.Infof("    shared by (%d): %s", len(e.images), strings.Join(e.images, ", "))
		log.Logger.Infof("")
	}

	log.Logger.Warn("Polluting a file in one of these layers will affect all listed images.")
}
