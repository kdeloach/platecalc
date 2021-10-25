package platecalc

import (
	"fmt"
	"strings"
)

type Node struct {
	Parent *Node
	Plate  float32
	Weight int
	Depth  int
	Score  int
	Nodes  map[float32]*Node
}

func NewTree(weight int) *Node {
	return &Node{
		Parent: nil,
		Plate:  0,
		Weight: weight,
		Depth:  0,
		Score:  0,
		Nodes:  make(map[float32]*Node),
	}
}

func NewNode(parent *Node, plate float32) *Node {
	weight := parent.Weight + int(plate*2)
	depth := parent.Depth + 1
	score := parent.Score + depth*int(plate)
	return &Node{
		Parent: parent,
		Plate:  plate,
		Weight: weight,
		Depth:  depth,
		Score:  score,
		Nodes:  make(map[float32]*Node),
	}
}

func (node *Node) Find(plates ...float32) *Node {
	if len(plates) == 0 {
		return nil
	}

	plate, rest := plates[0], plates[1:]
	next, ok := node.Nodes[plate]
	if !ok {
		return nil
	}

	if len(rest) > 0 {
		return next.Find(rest...)
	} else {
		return next
	}
}

func (node *Node) Add(plates ...float32) *Node {
	if len(plates) == 0 {
		return nil
	}

	plate, rest := plates[0], plates[1:]
	next, ok := node.Nodes[plate]
	if !ok {
		next = NewNode(node, plate)
		node.Nodes[plate] = next
	}

	if len(rest) > 0 {
		return next.Add(rest...)
	} else {
		return next
	}
}

type WalkNodeFn func(*Node)

func (node *Node) Walk(fn WalkNodeFn) {
	fn(node)
	for _, child := range node.Nodes {
		child.Walk(fn)
	}
}

type WalkNearbyNodeFn func(*Node, int)

func (node *Node) WalkNearby(maxDistance int, fn WalkNearbyNodeFn) {
	seen := make(map[*Node]bool)

	var walk func(*Node, int)
	walk = func(node *Node, distance int) {
		if node == nil || distance > maxDistance {
			return
		}
		if _, ok := seen[node]; ok {
			return
		}

		fn(node, distance)
		seen[node] = true

		for _, child := range node.Nodes {
			walk(child, distance+1)
		}
		walk(node.Parent, distance+1)
	}

	walk(node, 0)
}

func (node *Node) String() string {
	plates := make([]string, 0)
	for parent := node; parent != nil; parent = parent.Parent {
		if parent.Parent != nil {
			plates = append([]string{fmt.Sprintf("%v", parent.Plate)}, plates...)
		}
	}
	return strings.Join(plates, ", ")
}
