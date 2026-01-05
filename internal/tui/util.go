package tui

func boolToIcon(b bool) string {
	i := "✗"
	if b {
		i = "✓"
	}

	return i
}

func shortenSourceType(s string) string {
	switch s {
	case "person":
		return "p"
	case "room":
		return "r"
	default:
		return "?"
	}
}
