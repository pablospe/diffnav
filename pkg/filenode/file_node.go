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
	IconsNerdStatus   = "nerd-fonts-status"
	IconsNerdSimple   = "nerd-fonts-simple"
	IconsNerdFiletype = "nerd-fonts-filetype"
	IconsNerdFull     = "nerd-fonts-full"
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

	// full has a special layout: [status icon] [filename] [file-type icon]
	if f.IconStyle == IconsNerdFull {
		return f.renderFullLayout(name)
	}

	// All other styles: [icon] [filename] with optional coloring
	return f.renderStandardLayout(name)
}

// renderStandardLayout renders: [icon colored] [filename]
// Used by status, simple, filetype, unicode, ascii.
func (f FileNode) renderStandardLayout(name string) string {
	icon := f.getIcon()
	depthWidth := f.Depth * 2
	iconWidth := lipgloss.Width(icon) + 1
	nameMaxWidth := constants.OpenFileTreeWidth - depthWidth - iconWidth
	truncatedName := utils.TruncateString(name, nameMaxWidth)
	coloredIcon := lipgloss.NewStyle().Foreground(f.StatusColor()).Render(icon)

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
		return coloredIcon + " " + bgStyle.Render(truncatedName)
	}

	if f.ColorFileNames {
		styledName := lipgloss.NewStyle().Foreground(f.StatusColor()).Render(truncatedName)
		return coloredIcon + " " + styledName
	}

	return coloredIcon + " " + truncatedName
}

// renderFullLayout renders: [status icon colored] [file-type icon colored] [filename colored]
// All elements colored by git status.
func (f FileNode) renderFullLayout(name string) string {
	statusIcon := f.getStatusIcon()
	fileIcon := icons.GetIcon(name, false)
	statusColor := f.StatusColor()
	coloredStatus := lipgloss.NewStyle().Foreground(statusColor).Render(statusIcon)
	coloredFileIcon := lipgloss.NewStyle().Foreground(statusColor).Render(fileIcon)

	depthWidth := f.Depth * 2
	iconsWidth := lipgloss.Width(statusIcon) + 1 + lipgloss.Width(fileIcon) + 1
	nameMaxWidth := constants.OpenFileTreeWidth - depthWidth - iconsWidth
	truncatedName := utils.TruncateString(name, nameMaxWidth)

	if f.Selected {
		bgStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(statusColor).
			Background(lipgloss.Color("#3a3a3a"))
		if f.PanelWidth > 0 {
			iconWidth := lipgloss.Width(statusIcon) + 1 + lipgloss.Width(fileIcon) + 1
			availableWidth := f.PanelWidth - iconWidth - (f.Depth * 2)
			if availableWidth > 0 {
				bgStyle = bgStyle.Width(availableWidth)
			}
		}
		return coloredStatus + " " + coloredFileIcon + " " + bgStyle.Render(truncatedName)
	}

	if f.ColorFileNames {
		coloredName := lipgloss.NewStyle().Foreground(statusColor).Render(truncatedName)
		return coloredStatus + " " + coloredFileIcon + " " + coloredName
	}

	return coloredStatus + " " + coloredFileIcon + " " + truncatedName
}

// getIcon returns the left icon based on the icon style.
func (f FileNode) getIcon() string {
	name := filepath.Base(f.Path())
	switch f.IconStyle {
	case IconsNerdStatus:
		if f.File.IsNew {
			return ""
		} else if f.File.IsDelete {
			return ""
		}
		return ""
	case IconsNerdSimple:
		return ""
	case IconsNerdFiletype:
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

// getStatusIcon returns the git status indicator icon (used by full layout).
// Uses the same boxed icons as status style.
func (f FileNode) getStatusIcon() string {
	if f.File.IsNew {
		return "\uf457" //
	} else if f.File.IsDelete {
		return "\ueadf" //
	}
	return "\uf459" //
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
