package policy

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/kubernetes/policy"
	"github.com/ctrsploit/sploit-spec/pkg/log"
)

// Check checks for dangerous permissions and prints the analysis results.
func Check(results []policy.CheckResult) error {
	Report(results)
	return nil
}

// Report prints the analysis results for the given dangerous permissions.
func Report(results []policy.CheckResult) {
	if len(results) == 0 {
		log.Logger.Infof("No dangerous permissions found")
		return
	}

	grouped := policy.GroupResultsByLevel(results)

	log.Logger.Infof("")
	log.Logger.Infof("=== ServiceAccount Token Dangerous Permissions ===")

	// Critical level
	if critical := grouped[policy.LevelCritical]; len(critical) > 0 {
		log.Logger.Errorf("")
		log.Logger.Errorf("[CRITICAL] %d critical permissions found:", len(critical))
		for _, r := range critical {
			printResult(r)
		}
	}

	// High level
	if high := grouped[policy.LevelHigh]; len(high) > 0 {
		log.Logger.Warnf("")
		log.Logger.Warnf("[HIGH] %d high-risk permissions found:", len(high))
		for _, r := range high {
			printResult(r)
		}
	}

	// Medium level
	if medium := grouped[policy.LevelMedium]; len(medium) > 0 {
		log.Logger.Infof("")
		log.Logger.Infof("[MEDIUM] %d medium-risk permissions found:", len(medium))
		for _, r := range medium {
			printResult(r)
		}
	}

	// Summary
	log.Logger.Infof("")
	log.Logger.Infof("=== Summary ===")
	log.Logger.Infof("Critical: %d, High: %d, Medium: %d",
		len(grouped[policy.LevelCritical]),
		len(grouped[policy.LevelHigh]),
		len(grouped[policy.LevelMedium]))
}

func printResult(r policy.CheckResult) {
	resource := r.Permission.FullResource()
	scope := "Cluster-Wide"
	if r.Namespace != "" {
		scope = fmt.Sprintf("Namespace: %s", r.Namespace)
	}
	log.Logger.Infof("  - %s [%s] %s", resource, r.MatchedVerb, scope)
	log.Logger.Infof("    Risk: %s", r.Permission.Description)
	if r.Permission.Reference != "" {
		log.Logger.Infof("    Ref: %s", r.Permission.Reference)
	}
}
