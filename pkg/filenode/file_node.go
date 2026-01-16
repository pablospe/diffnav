package filenode

import (
	"path/filepath"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"

	"github.com/dlvhdr/diffnav/pkg/constants"
	"github.com/dlvhdr/diffnav/pkg/icons"
	"github.com/dlvhdr/diffnav/pkg/utils"
)

// Icon style constants.
const (
	IconsNerdFonts    = "nerd-fonts"
	IconsNerdFontsAlt = "nerd-fonts-alt"
	IconsUnicode      = "unicode"
	IconsASCII        = "ascii"
)

type FileNode struct {
	File           *gitdiff.File
	Depth          int
	YOffset        int
	IconStyle      string
	Selected       bool
	ColorFileNames bool
	PanelWidth     int
}

func (f FileNode) Path() string {
	return GetFileName(f.File)
}

func (f FileNode) Value() string {
	name := filepath.Base(f.Path())
	icon := f.getIcon()
	statusIcon := f.getStatusIcon()

	// Color the status icon by git status
	coloredStatus := lipgloss.NewStyle().Foreground(f.StatusColor()).Render(statusIcon)

	depthWidth := f.Depth * 2
	iconsWidth := lipgloss.Width(icon) + 1 + lipgloss.Width(coloredStatus) + 1 // icon + space + status + space
	nameMaxWidth := constants.OpenFileTreeWidth - depthWidth - iconsWidth
	truncatedName := utils.TruncateString(name, nameMaxWidth)

	// Calculate spacer to push status icon to the right
	spacerWidth := constants.OpenFileTreeWidth - lipgloss.Width(truncatedName) - iconsWidth - depthWidth
	if len(truncatedName) < len(name) {
		spacerWidth = spacerWidth - 1
	}
	spacer := ""
	if spacerWidth > 0 {
		spacer = strings.Repeat(" ", spacerWidth)
	}

	if f.Selected {
		// Apply background with fixed width to extend to panel edge
		bgStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(f.StatusColor()).
			Background(lipgloss.Color("#3a3a3a"))
		if f.PanelWidth > 0 {
			iconWidth := lipgloss.Width(icon) + 1
			availableWidth := f.PanelWidth - iconWidth - (f.Depth * 2)
			if availableWidth > 0 {
				bgStyle = bgStyle.Width(availableWidth)
			}
		}
		return icon + " " + bgStyle.Render(truncatedName+spacer+" "+coloredStatus)
	}

	if f.ColorFileNames {
		styledName := lipgloss.NewStyle().Foreground(f.StatusColor()).Render(truncatedName)
		return icon + " " + styledName + spacer + " " + coloredStatus
	}

	return icon + " " + truncatedName + spacer + " " + coloredStatus
}

// getIcon returns the file-type icon based on the icon style.
func (f FileNode) getIcon() string {
	name := filepath.Base(f.Path())
	switch f.IconStyle {
	case IconsNerdFonts:
		// Use file-type specific icons from the icons package
		return icons.GetIcon(name, false)
	case IconsNerdFontsAlt:
		return ""
	case IconsUnicode:
		return "●"
	default: // ascii
		return "*"
	}
}

// getStatusIcon returns the git status indicator icon.
func (f FileNode) getStatusIcon() string {
	switch f.IconStyle {
	case IconsNerdFonts:
		if f.File.IsNew {
			return ""
		} else if f.File.IsDelete {
			return ""
		}
		return ""
	case IconsNerdFontsAlt:
		if f.File.IsNew {
			return ""
		} else if f.File.IsDelete {
			return ""
		}
		return ""
	case IconsUnicode:
		if f.File.IsNew {
			return "+"
		} else if f.File.IsDelete {
			return "⛌"
		}
		return "●"
	default: // ascii
		if f.File.IsNew {
			return "+"
		} else if f.File.IsDelete {
			return "x"
		}
		return "*"
	}
}

// StatusColor returns the color for this file based on its git status.
func (f FileNode) StatusColor() lipgloss.Color {
	if f.File.IsNew {
		return lipgloss.Color("2") // green
	} else if f.File.IsDelete {
		return lipgloss.Color("1") // red
	}
	return lipgloss.Color("3") // yellow/orange
}

func (f FileNode) String() string {
	return f.Value()
}

func (f FileNode) Children() tree.Children {
	return tree.NodeChildren(nil)
}

func (f FileNode) Hidden() bool {
	return false
}

func (f FileNode) SetHidden(bool) {}

func (f FileNode) SetValue(any) {}

func GetFileName(file *gitdiff.File) string {
	if file.NewName != "" {
		return file.NewName
	}
	return file.OldName
}
