package platecalc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermutations(t *testing.T) {
	got := Permutations(1, 2, 3)
	want := [][]float32{
		{1},
		{1, 2},
		{1, 2, 3},
		{1, 3},
		{1, 3, 2},
		{2},
		{2, 1},
		{2, 1, 3},
		{2, 3},
		{2, 3, 1},
		{3},
		{3, 1},
		{3, 1, 2},
		{3, 2},
		{3, 2, 1},
	}
	assert.Equal(t, want, got)
}

func TestBestSolution(t *testing.T) {
	tree := NewTree(nil, 45)
	for _, p := range Permutations(5, 5, 10, 10, 2.5) {
		tree.Add(p...)
	}

	sets := []int{55, 65, 75, 55}
	result := BestSolution(tree, sets, 5, false)

	got := make([]string, 0)
	for _, node := range result {
		got = append(got, node.String())
	}

	want := []string{
		"5",
		"5, 5",
		"5, 10",
		"5",
	}
	assert.Equal(t, want, got)
}

func TestSimpleSolution(t *testing.T) {
	tree := NewTree(nil, 45)
	for _, p := range Permutations(5, 5, 10, 10, 2.5) {
		tree.Add(p...)
	}

	sets := []int{55, 65, 75, 55}
	result := SimpleSolution(tree, sets, false)

	got := make([]string, 0)
	for _, node := range result {
		got = append(got, node.String())
	}

	want := []string{
		"5",
		"10",
		"10, 5",
		"5",
	}
	assert.Equal(t, want, got)
}

func TestRoundUpToNearest(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{4, 5},
		{6, 10},
		{9, 10},
		{11, 15},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			assert.Equal(t, tc.want, RoundUpToNearest(float32(tc.n), 5))
		})
	}
}
