package diagnose

import (
	"strings"
	"testing"

	"github.com/AnuragRaut08/dataquarry/internal/summary"
)

func TestCheckSmallRowGroupsHealthy(t *testing.T) {
	s := summary.Summary{
		NumRowGroups:     4,
		UncompressedSize: 512 * 1024 * 1024, // 128 MiB average
	}

	finding := CheckSmallRowGroups(s)

	if finding != nil {
		t.Fatalf("expected no finding, got %+v", finding)
	}
}

func TestCheckSmallRowGroupsDetectsSmallGroups(t *testing.T) {
	s := summary.Summary{
		NumRowGroups:     40,
		UncompressedSize: 40 * 96 * 1024, // 96 KiB average
	}

	finding := CheckSmallRowGroups(s)

	if finding == nil {
		t.Fatal("expected finding, got nil")
	}

	if finding.Title != "Row groups are much smaller than recommended" {
		t.Fatalf("unexpected title: %q", finding.Title)
	}

	if finding.Severity != SeverityHigh {
		t.Fatalf("unexpected severity: %q", finding.Severity)
	}

	if !strings.Contains(finding.Evidence, "96.0 KiB") {
		t.Fatalf("unexpected evidence: %q", finding.Evidence)
	}

	if finding.Fix == "" {
		t.Fatal("expected fix to be non-empty")
	}
}

func TestCheckSmallRowGroupsBoundary(t *testing.T) {
	s := summary.Summary{
		NumRowGroups:     2,
		UncompressedSize: 2 * 64 * 1024 * 1024, // exactly 64 MiB average
	}

	finding := CheckSmallRowGroups(s)

	if finding != nil {
		t.Fatalf("expected no finding at threshold, got %+v", finding)
	}
}

func TestCheckSmallRowGroupsZeroRowGroups(t *testing.T) {
	s := summary.Summary{
		NumRowGroups:     0,
		UncompressedSize: 100 * 1024 * 1024,
	}

	finding := CheckSmallRowGroups(s)

	if finding != nil {
		t.Fatalf("expected no finding, got %+v", finding)
	}
}
