package filenode

import (
	"path/filepath"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"

	"github.com/dlvhdr/diffnav/pkg/constants"
	"github.com/dlvhdr/diffnav/pkg/icons"
	"github.com/dlvhdr/diffnav/pkg/utils"
)

// Icon style constants.
const (
	IconsNerdFonts     = "nerd-fonts"
	IconsNerdFontsAlt  = "nerd-fonts-alt"
	IconsNerdFontsAlt2 = "nerd-fonts-alt2"
	IconsNerdFontsAlt3 = "nerd-fonts-alt3"
	IconsUnicode       = "unicode"
	IconsASCII         = "ascii"
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

	// nerd-fonts-alt3 has a special layout: [status icon] [filename] [file-type icon]
	if f.IconStyle == IconsNerdFontsAlt3 {
		return f.renderAlt3Layout(name)
	}

	// All other styles: [icon] [filename] with optional coloring
	return f.renderStandardLayout(name)
}

// renderStandardLayout renders: [icon] [filename]
// Used by nerd-fonts, nerd-fonts-alt, nerd-fonts-alt2, unicode, ascii.
func (f FileNode) renderStandardLayout(name string) string {
	icon := f.getIcon()
	depthWidth := f.Depth * 2
	iconWidth := lipgloss.Width(icon) + 1
	nameMaxWidth := constants.OpenFileTreeWidth - depthWidth - iconWidth
	truncatedName := utils.TruncateString(name, nameMaxWidth)

	// nerd-fonts-alt: icon not colored, but filename can be colored
	// All other styles: both icon and filename colored
	colorIcon := f.IconStyle != IconsNerdFontsAlt

	if f.Selected {
		bgStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(f.StatusColor()).
			Background(lipgloss.Color("#3a3a3a"))
		if f.PanelWidth > 0 {
			availableWidth := f.PanelWidth - iconWidth - (f.Depth * 2)
			if availableWidth > 0 {
				bgStyle = bgStyle.Width(availableWidth)
			}
		}
		displayIcon := icon
		if colorIcon {
			displayIcon = lipgloss.NewStyle().Foreground(f.StatusColor()).Render(icon)
		}
		return displayIcon + " " + bgStyle.Render(truncatedName)
	}

	if f.ColorFileNames {
		styledName := lipgloss.NewStyle().Foreground(f.StatusColor()).Render(truncatedName)
		if colorIcon {
			coloredIcon := lipgloss.NewStyle().Foreground(f.StatusColor()).Render(icon)
			return coloredIcon + " " + styledName
		}
		return icon + " " + styledName
	}

	if colorIcon {
		coloredIcon := lipgloss.NewStyle().Foreground(f.StatusColor()).Render(icon)
		return coloredIcon + " " + truncatedName
	}

	return icon + " " + truncatedName
}

// renderAlt3Layout renders: [status icon colored] [file-type icon] [filename]
// Status indicator colored, file-type icon and filename plain.
func (f FileNode) renderAlt3Layout(name string) string {
	statusIcon := f.getStatusIcon()
	fileIcon := icons.GetIcon(name, false)
	coloredStatus := lipgloss.NewStyle().Foreground(f.StatusColor()).Render(statusIcon)

	depthWidth := f.Depth * 2
	iconsWidth := lipgloss.Width(statusIcon) + 1 + lipgloss.Width(fileIcon) + 1
	nameMaxWidth := constants.OpenFileTreeWidth - depthWidth - iconsWidth
	truncatedName := utils.TruncateString(name, nameMaxWidth)

	if f.Selected {
		bgStyle := lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#3a3a3a"))
		if f.PanelWidth > 0 {
			iconWidth := lipgloss.Width(statusIcon) + 1 + lipgloss.Width(fileIcon) + 1
			availableWidth := f.PanelWidth - iconWidth - (f.Depth * 2)
			if availableWidth > 0 {
				bgStyle = bgStyle.Width(availableWidth)
			}
		}
		return coloredStatus + " " + fileIcon + " " + bgStyle.Render(truncatedName)
	}

	return coloredStatus + " " + fileIcon + " " + truncatedName
}

// getIcon returns the left icon based on the icon style.
func (f FileNode) getIcon() string {
	name := filepath.Base(f.Path())
	switch f.IconStyle {
	case IconsNerdFonts:
		if f.File.IsNew {
			return ""
		} else if f.File.IsDelete {
			return ""
		}
		return ""
	case IconsNerdFontsAlt:
		return ""
	case IconsNerdFontsAlt2:
		return icons.GetIcon(name, false) // File-type specific icon (colored by status)
	case IconsUnicode:
		if f.File.IsNew {
			return "+"
		} else if f.File.IsDelete {
			return "⛌"
		}
		return "●"
	default: // ascii (fallback for unknown values)
		if f.File.IsNew {
			return "+"
		} else if f.File.IsDelete {
			return "x"
		}
		return "*"
	}
}

// getStatusIcon returns the git status indicator icon (used by alt3 layout).
func (f FileNode) getStatusIcon() string {
	if f.File.IsNew {
		return ""
	} else if f.File.IsDelete {
		return ""
	}
	return ""
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
