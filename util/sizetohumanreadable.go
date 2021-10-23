package util

import (
	"fmt"
)

// Return a human-readable representation of a file size.
func SizeToHumanReadableSize(size int64) string {
	sizes := "KMGPTE"
	if size < 1024 {
		return fmt.Sprintf("%d\u202fB", size)
	}
	fsize := float64(size)
	for _, prefix := range sizes {
		fsize /= 1024
		if fsize < 1024 {
			return fmt.Sprintf("%.2f\u202f%ciB", fsize, prefix)
		}
	}
	return fmt.Sprintf("%.2f\u202fEiB", fsize)
}
