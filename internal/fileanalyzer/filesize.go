package fileanalyzer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var re = regexp.MustCompile("[0-9]*(.[0-9]+)?")

func ParseSize(sizeStr string) (int64, error) {
	match := re.FindStringIndex(sizeStr)

	if match == nil || len(match) != 2 {
		return 0, fmt.Errorf("invalid size: %s", sizeStr)
	}

	size, err := strconv.ParseFloat(sizeStr[:match[1]], 64)

	if err != nil {
		return 0, fmt.Errorf("invalid size: %s", sizeStr)
	}

	unit := sizeStr[match[1]:]

	switch strings.ToUpper(unit) {
	case "B":
		return int64(size), nil
	case "KB", "K":
		return int64(size * 1_000), nil
	case "MB", "M":
		return int64(size * 1_000_000), nil
	case "GB", "G":
		return int64(size * 1_000_000_000), nil
	case "TB", "T":
		return int64(size * 1_000_000_000_000), nil
	default:
		return 0, fmt.Errorf("invalid size unit: %s", unit)
	}
}

func FormatSize(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
