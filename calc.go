package platecalc

import (
	"math"
)

func Permutations(plates ...float32) [][]float32 {
	platesColl := make([][]float32, 0)
	if len(plates) == 0 {
		return platesColl
	}
	for i, p := range plates {
		subPlates := make([]float32, 0, len(plates)-1)
		subPlates = append(subPlates, plates[:i]...)
		subPlates = append(subPlates, plates[i+1:]...)

		platesColl = append(platesColl, []float32{p})

		for _, pz := range Permutations(subPlates...) {
			perm := append([]float32{p}, pz...)
			platesColl = append(platesColl, perm)
		}
	}
	return platesColl
}

func BestSolution(tree *Tree, setWeights []int, maxDistance int) []*Tree {
	if len(setWeights) == 0 {
		return nil
	}

	bestScore := math.MaxInt32
	var solution []*Tree

	foundSolution := func(score int, nodes []*Tree) {
		if score < bestScore {
			bestScore = score
			solution = nodes
			// for _, n := range nodes {
			// 	fmt.Printf("%3v: %v (score=%v)\n", n.Weight, n, n.Score())
			// }
			// fmt.Printf("total=%v\n\n", score)
		}
	}
	nextFn := foundSolution

	head, tail := setWeights[0], setWeights[1:]

	for i := len(tail) - 1; i >= 0; i-- {
		weight := tail[i]
		oldNextFn := nextFn
		nextFn = func(prevScore int, prevNodes []*Tree) {
			prevNode := prevNodes[len(prevNodes)-1]
			prevNode.WalkNearby(maxDistance, func(node *Tree, dist int) {
				if node.TotalWeight() == weight {
					nodes := make([]*Tree, len(prevNodes))
					copy(nodes, prevNodes)
					nodes = append(nodes, node)

					// score := prevScore + node.Score()*dist
					// fmt.Printf("%v + %v * %v = %v\n", prevScore, node.Score(), dist, score)

					// Optimize for number of plates added/removed
					score := prevScore + dist

					oldNextFn(score, nodes)
				}
			})
		}
	}

	tree.Walk(func(node *Tree) {
		if node.TotalWeight() == head {
			nodes := []*Tree{node}
			nextFn(node.Score(), nodes)
		}
	})

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