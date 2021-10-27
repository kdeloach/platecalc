package platecalc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProps(t *testing.T) {
	tree := NewTree(nil, 45)
	node := tree.Add(45, 35, 25)
	assert.Equal(t, float32(25), node.Value)
	assert.Equal(t, 255, node.TotalWeight())
	assert.Equal(t, 3, node.Depth)
	assert.Equal(t, 190, node.Score())
	assert.Equal(t, "45, 35, 25", node.String())
}

func TestWalk(t *testing.T) {
	tree := NewTree(nil, 0)
	tree.Add(1, 2, 3)
	tree.Add(1, 4, 5)

	got := make([]string, 0)
	tree.Walk(func(node *Tree) {
		if node.Parent != nil {
			got = append(got, node.String())
		}
	})

	want := []string{
		"1",
		"1, 2",
		"1, 2, 3",
		"1, 4",
		"1, 4, 5",
	}
	assert.ElementsMatch(t, want, got)
}

func TestWalkNearby(t *testing.T) {
	tree := NewTree(nil, 0)
	tree.Add(1, 2, 3)
	tree.Add(1, 2, 4)

	node := tree.Find(1, 2)
	assert.NotNil(t, node)

	got := make([]string, 0)
	node.WalkNearby(1, func(node *Tree, dist int) {
		if node.Parent != nil {
			got = append(got, node.String())
		}
	})

	want := []string{
		"1, 2",
		"1, 2, 3",
		"1, 2, 4",
		"1",
	}
	assert.ElementsMatch(t, want, got)
}
