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
	m := Model{
		t: tree.New(nil, constants.OpenFileTreeWidth, 0),
	}

	m.t.Enumerator(enumerator).Indenter(indenter)

	return m
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
	collapsed := collapseTree(t)
	// t, _ = truncateTree(t, 0, 0, 0)
	m.t.SetNodes(collapsed)
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

		// path does not exist from this point, need to create it
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

// Given a tree with nodes that have only one child, collapse the tree by
// merging these nodes with their parents, as long as the parent has only one child as well.
// For example, the tree:
// .
// ├── a
// │   └── b
// │       └── c
//
// will be collapsed to:
// .
// ├── a/b
// │   └── c
//
// This tree wouldn't be collapsed:
// ├── a
// │   ├── b
// │   │   └── c
// │   └── d
func collapseTree(t *tree.Node) *tree.Node {
	children := t.ChildNodes()
	newT := tree.Root(t.GivenValue())
	if len(children) == 0 {
		return newT
	}

	// recursively collapse children
	for _, child := range children {
		if child.Size() > 1 {
			collapsedChild := collapseTree(child)
			newT.Child(collapsedChild)
		} else {
			newT.Child(child.GivenValue())
		}
	}

	// all children are collapsed, now check if the parent can be collapsed
	newChildren := newT.ChildNodes()
	if len(newChildren) == 1 {
		child := newChildren[0]
		if t.GivenValue() == "." {
			return child
		}

		val := t.Value() + string(os.PathSeparator) + child.Value()
		collapsed := tree.Root(val).Child(child.ChildNodes())
		return collapsed
	}

	return newT
}

var enumerator = func(children ltree.Children, index int) string {
	return "│"
}

var indenter = func(children ltree.Children, index int) string {
	return "│"
}

// SetSize implements the Component interface.
func (m *Model) SetSize(width, height int) tea.Cmd {
	m.t.SetSize(width, height)
	return nil
}
