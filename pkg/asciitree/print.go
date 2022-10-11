package asciitree

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/mgutz/ansi"
)

// taken from https://github.com/Tufin/asciitree

type Tree map[string]Tree

func (tree Tree) Add(path string) {
	frags := strings.Split(path, "#")
	tree.add(frags)
}

func (tree Tree) add(frags []string) {
	if len(frags) == 0 {
		return
	}

	nextTree, ok := tree[frags[0]]
	if !ok {
		nextTree = Tree{}
		tree[frags[0]] = nextTree
	}

	nextTree.add(frags[1:])
}

func (tree Tree) Fprint(w io.Writer, root bool, padding string) {
	if tree == nil {
		return
	}

	index := 0

	keys := make([]string, 0)
	for k := range tree {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := tree[k]
		fmt.Fprintf(w, "%s%s\n", padding+getPadding(root, getBoxType(index, len(tree))), k)
		v.Fprint(w, false, padding+getPadding(root, getBoxTypeExternal(index, len(tree))))
		index++
	}
}

type BoxType int

const (
	Regular BoxType = iota
	Last
	AfterLast
	Between
)

func (boxType BoxType) String() string {
	phosphorize := ansi.ColorFunc("blue+h:black")
	switch boxType {
	case Regular:
		return phosphorize("\u251c") // ├
	case Last:
		return phosphorize("\u2514") // └
	case AfterLast:
		return " "
	case Between:
		return phosphorize("\u2502") // │
	default:
		panic("invalid box type")
	}
}

func getBoxType(index int, len int) BoxType {
	if index+1 == len {
		return Last
	} else if index+1 > len {
		return AfterLast
	}
	return Regular
}

func getBoxTypeExternal(index int, len int) BoxType {
	if index+1 == len {
		return AfterLast
	}
	return Between
}

func getPadding(root bool, boxType BoxType) string {
	if root {
		return ""
	}

	return boxType.String() + " "
}
