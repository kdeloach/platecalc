package platecalc

import (
	"fmt"
	"math"
	"strings"
)

type Tree struct {
	Parent   *Tree
	Depth    int
	Children map[float32]*Tree
	Value    float32
}

func NewTree(parent *Tree, value float32) *Tree {
	depth := 0
	if parent != nil {
		depth = parent.Depth + 1
	}
	return &Tree{
		Parent:   parent,
		Depth:    depth,
		Children: make(map[float32]*Tree),
		Value:    value,
	}
}

func (t *Tree) Score(preferLessPlates bool) int {
	if t.Parent == nil {
		return 0
	}

	var scale float32 = 1

	if !preferLessPlates {
		// scale up heavier plates; prefer lighter plates
		scale = float32(math.Round(float64(t.Value) / 10))
		if scale < 1 {
			scale = 1
		}
	}

	return t.Parent.Score(preferLessPlates) + t.Depth*int(t.Value*scale)
}

func (t *Tree) TotalWeight() int {
	if t.Parent == nil {
		return int(t.Value)
	}
	return t.Parent.TotalWeight() + int(t.Value*2)
}

func (t *Tree) Find(plates ...float32) *Tree {
	if len(plates) == 0 {
		return nil
	}

	plate, rest := plates[0], plates[1:]
	next, ok := t.Children[plate]
	if !ok {
		return nil
	}

	if len(rest) > 0 {
		return next.Find(rest...)
	} else {
		return next
	}
}

func (t *Tree) Add(plates ...float32) *Tree {
	if len(plates) == 0 {
		return nil
	}

	plate, rest := plates[0], plates[1:]
	next, ok := t.Children[plate]
	if !ok {
		next = NewTree(t, plate)
		t.Children[plate] = next
	}

	if len(rest) > 0 {
		return next.Add(rest...)
	} else {
		return next
	}
}

type WalkTreeFn func(*Tree)

func (t *Tree) Walk(fn WalkTreeFn) {
	fn(t)
	for _, child := range t.Children {
		child.Walk(fn)
	}
}

type WalkNearbyTreeFn func(*Tree, int)

func (t *Tree) WalkNearby(maxDistance int, fn WalkNearbyTreeFn) {
	seen := make(map[*Tree]bool)

	var walk func(*Tree, int)
	walk = func(t *Tree, distance int) {
		if t == nil || distance > maxDistance {
			return
		}
		if _, ok := seen[t]; ok {
			return
		}

		fn(t, distance)
		seen[t] = true

		for _, child := range t.Children {
			walk(child, distance+1)
		}
		walk(t.Parent, distance+1)
	}

	walk(t, 0)
}

func (t *Tree) String() string {
	plates := make([]string, 0)
	for parent := t; parent != nil; parent = parent.Parent {
		if parent.Parent != nil {
			plates = append([]string{fmt.Sprintf("%v", parent.Value)}, plates...)
		}
	}
	return strings.Join(plates, ", ")
}
