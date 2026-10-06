package diagnose

import (
	"fmt"

	"github.com/AnuragRaut08/dataquarry/internal/summary"
)

const (
	targetRowGroupSize = int64(128 * 1024 * 1024)
	minRowGroupSize    = targetRowGroupSize / 2
)

// CheckSmallRowGroups detects datasets whose average row-group size
// is significantly smaller than the recommended 128 MiB target.
func CheckSmallRowGroups(s summary.Summary) *Finding {
	if s.NumRowGroups <= 0 || s.UncompressedSize <= 0 {
		return nil
	}

	averageSize := s.UncompressedSize / int64(s.NumRowGroups)

	if averageSize >= minRowGroupSize {
		return nil
	}

	return findingForSmallRowGroups(averageSize, s.NumRowGroups)
}

func findingForSmallRowGroups(averageSize int64, numRowGroups int) *Finding {
	return &Finding{
		Title:    "Row groups are much smaller than recommended",
		Severity: SeverityHigh,
		Evidence: fmt.Sprintf(
			"Average row group size is %s across %d row groups.",
			formatBytes(averageSize),
			numRowGroups,
		),
		Fix: "Rewrite the file using a larger target row group size, such as 128 MiB.",
	}
}

func formatBytes(size int64) string {
	const (
		KiB = int64(1024)
		MiB = KiB * 1024
	)

	switch {
	case size >= MiB:
		return fmt.Sprintf("%.1f MiB", float64(size)/float64(MiB))
	case size >= KiB:
		return fmt.Sprintf("%.1f KiB", float64(size)/float64(KiB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
