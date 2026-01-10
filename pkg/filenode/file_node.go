package filenode

import (
	"path/filepath"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"

	"github.com/dlvhdr/diffnav/pkg/constants"
	"github.com/dlvhdr/diffnav/pkg/utils"
)

// Icon style constants.
const (
	IconsNerdFonts = "nerd-fonts"
	IconsUnicode   = "unicode"
	IconsASCII     = "ascii"
)

type FileNode struct {
	File      *gitdiff.File
	Depth     int
	YOffset   int
	IconStyle string
}

func (f FileNode) Path() string {
	return GetFileName(f.File)
}

func (f FileNode) Value() string {
	icon := " "
	status := " " + f.getStatusIcon()

	depthWidth := f.Depth * 2
	iconsWidth := lipgloss.Width(icon) + lipgloss.Width(status)
	nameMaxWidth := constants.OpenFileTreeWidth - depthWidth - iconsWidth
	base := filepath.Base(f.Path())
	name := utils.TruncateString(base, nameMaxWidth)

	spacerWidth := constants.OpenFileTreeWidth - lipgloss.Width(name) - iconsWidth - depthWidth
	if len(name) < len(base) {
		spacerWidth = spacerWidth - 1
	}
	spacer := ""
	if spacerWidth > 0 {
		spacer = strings.Repeat(" ", spacerWidth)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, icon, name, spacer, status)
}

func (f FileNode) getStatusIcon() string {
	color := f.StatusColor()
	style := lipgloss.NewStyle().Foreground(color)

	switch f.IconStyle {
	case IconsNerdFonts:
		if f.File.IsNew {
			return style.Render("")
		} else if f.File.IsDelete {
			return style.Render("")
		}
		return style.Render("")
	case IconsUnicode:
		if f.File.IsNew {
			return style.Render("+")
		} else if f.File.IsDelete {
			return style.Render("⛌")
		}
		return style.Render("●")
	default: // ascii
		if f.File.IsNew {
			return style.Render("+")
		} else if f.File.IsDelete {
			return style.Render("x")
		}
		return style.Render("*")
	}
}

// StatusColor returns the color for this file based on its git status.
func (f FileNode) StatusColor() lipgloss.Color {
	if f.File.IsNew {
		return lipgloss.Color("2") // green
	} else if f.File.IsDelete {
		return lipgloss.Color("1") // red
	}
	return lipgloss.Color("3") // yellow
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
