package platecalc

import (
	"fmt"
	"math"
)

// Permutations returns every possible combination as tuples of length 1 to N.
// Ex: [1, 2, 3] -> [[1], [1,2], [1,2,3], [1,3,2], [2], [2,1], ...]
func Permutations(plates ...float32) [][]float32 {
	tuples := make([][]float32, 0)
	if len(plates) == 0 {
		return tuples
	}
	for i, p := range plates {
		tuples = append(tuples, []float32{p})

		// create new list with p removed
		newPlates := make([]float32, 0, len(plates)-1)
		newPlates = append(newPlates, plates[:i]...)
		newPlates = append(newPlates, plates[i+1:]...)

		// prefix each child tuple with p
		for _, tup := range Permutations(newPlates...) {
			newTuple := append([]float32{p}, tup...)
			tuples = append(tuples, newTuple)
		}
	}
	return tuples
}

// BestSolution returns the optimal sequence of plate changes for setWeights
// by walking the permutation tree and selecting the closest nodes with the
// lowest combined score.
func BestSolution(tree *Tree, setWeights []int, maxDistance int, debug bool) []*Tree {
	if len(setWeights) == 0 {
		return nil
	}

	bestScore := math.MaxInt32
	var solution []*Tree

	foundSolution := func(score int, nodes []*Tree, direct bool) {
		if score < bestScore {
			bestScore = score
			solution = nodes
			if debug {
				for _, n := range nodes {
					fmt.Printf("%3v: %v (score=%v)\n", n.TotalWeight(), n, n.Score())
				}
				fmt.Printf("total=%v direct=%v\n\n", score, direct)
			}
		}
	}
	nextFn := foundSolution

	head, tail := setWeights[0], setWeights[1:]

	for i := len(tail) - 1; i >= 0; i-- {
		weight := tail[i]
		oldNextFn := nextFn
		nextFn = func(prevScore int, prevNodes []*Tree, prevDirect bool) {
			prevNode := prevNodes[len(prevNodes)-1]
			prevNode.WalkNearby(maxDistance, func(node *Tree, dist int, isDirectRelative bool) {
				if node.TotalWeight() == weight {
					nodes := make([]*Tree, len(prevNodes))
					copy(nodes, prevNodes)
					nodes = append(nodes, node)

					handicapFactor := 1
					if !isDirectRelative {
						handicapFactor = 2
					}

					// Optimize for number of plates added/removed by weight
					score := prevScore + node.Score()*dist*handicapFactor

					// Optimize for number of plates added/removed by count
					// score := prevScore + dist

					oldNextFn(score, nodes, isDirectRelative)
				}
			})
		}
	}

	tree.Walk(func(node *Tree) {
		if node.TotalWeight() == head {
			nodes := []*Tree{node}
			nextFn(node.Score(), nodes, true)
		}
	})

	return solution
}

// SimpleSolution returns the best plate arrangement for each individual weight
// in setWeights.
func SimpleSolution(tree *Tree, setWeights []int, debug bool) []*Tree {
	solution := make([]*Tree, 0)

	for _, weight := range setWeights {
		best := BestSolution(tree, []int{weight}, 0, debug)
		if best == nil {
			return nil
		}
		solution = append(solution, best[0])
	}

	return solution

}

func RoundUpToNearest(n float32, inc int) int {
	return inc * int(math.Ceil(float64(n)/float64(inc)))
}

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
