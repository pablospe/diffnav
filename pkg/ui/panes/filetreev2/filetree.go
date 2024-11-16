package filetreev2

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/charmbracelet/bubbles/tree"
	tea "github.com/charmbracelet/bubbletea"
	ltree "github.com/charmbracelet/lipgloss/tree"

	"github.com/dlvhdr/diffnav/pkg/constants"
	"github.com/dlvhdr/diffnav/pkg/filenode"
)

type Model struct {
	t     tree.Model
	files []*gitdiff.File
}

func New() Model {
	return Model{
		t: tree.New(nil, constants.OpenFileTreeWidth, 50),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return m.t.View()
}

func (m Model) SetFiles(files []*gitdiff.File) Model {
	m.files = files
	t := buildFullFileTree(files)
	m.t.SetNodes(t)
	return m
}

func (m Model) SetCursor(cursor int) Model {
	return m
}

func buildFullFileTree(files []*gitdiff.File) *tree.Node {
	t := tree.Root(".")
	for _, file := range files {
		subTree := t

		name := filenode.GetFileName(file)
		dir := filepath.Dir(name)
		parts := strings.Split(dir, string(os.PathSeparator))
		path := ""

		// walk the tree to find existing path
		for _, part := range parts {
			found := false
			children := subTree.ChildNodes()
			for j := 0; j < len(children); j++ {
				child := children[j]
				if child.Value() == part {
					subTree = child
					path = path + part + string(os.PathSeparator)
					found = true
					break
				}
			}
			if !found {
				break
			}
		}

		// path does not exist from this point, need to creat it
		leftover := strings.TrimPrefix(name, path)
		parts = strings.Split(leftover, string(os.PathSeparator))
		for i, part := range parts {
			var c *tree.Node
			if i == len(parts)-1 {
				subTree.Child(filenode.FileNode{File: file})
			} else {
				c = tree.Root(part)
				subTree.Child(c)
				subTree = c
			}
		}
	}

	return t
}

var enumerator = func(children ltree.Children, index int) string {
	return "│"
}

var indenter = func(children ltree.Children, index int) string {
	if children.Length()-1 == index {
		return " "
	}
	return "│"
}

// SetSize implements the Component interface.
func (m *Model) SetSize(width, height int) tea.Cmd {
	return nil
}
