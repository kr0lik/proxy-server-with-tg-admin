package helper

import "fmt"

const (
	kb = 1 << 10 // 1024
	mb = 1 << 20 // 1024 * 1024
	gb = 1 << 30 // 1024 * 1024 * 1024
)

func BytesFormat(bytes uint64) string {
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.2fGB", float64(bytes)/gb)
	case bytes >= mb:
		return fmt.Sprintf("%.2fMB", float64(bytes)/mb)
	case bytes >= kb:
		return fmt.Sprintf("%.2fKB", float64(bytes)/kb)
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}
