package common

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// RenderScrollbar renders a vertical scrollbar track with a thumb indicator.
// It takes the viewport height, total number of content lines, and the current
// scroll offset (YOffset). Returns an empty string if all content fits.
func RenderScrollbar(viewHeight, totalLines, yOffset int) string {
	if totalLines <= viewHeight {
		return ""
	}

	trackHeight := viewHeight
	thumbSize := max(1, trackHeight*viewHeight/totalLines)

	scrollableLines := totalLines - viewHeight
	thumbPos := 0
	if scrollableLines > 0 {
		thumbPos = yOffset * (trackHeight - thumbSize) / scrollableLines
		if yOffset > 0 && thumbPos == 0 {
			thumbPos = 1
		}
	}

	track := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	thumb := lipgloss.NewStyle().Foreground(lipgloss.Blue)

	var sb strings.Builder
	for i := 0; i < trackHeight; i++ {
		if i > 0 {
			sb.WriteByte('\n')
		}
		if i >= thumbPos && i < thumbPos+thumbSize {
			sb.WriteString(thumb.Render("┃"))
		} else {
			sb.WriteString(track.Render("│"))
		}
	}
	return sb.String()
}
