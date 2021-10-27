package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/kdeloach/platecalc"
)

var barWeight = flag.Int("bar", 45, "bar weight")
var platesFlag = flag.String("plates", "45,35,25,10,10,5,5,2.5", "available plates")
var maxDistance = flag.Int("maxdistance", 5, "maximum distance to search tree")

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: calc [weight:int]+\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	plates, err := parsePlates()
	if err != nil {
		log.Fatalf(err.Error())
	}

	setWeights, err := parseWeights()
	if err != nil {
		log.Fatalf(err.Error())
	}

	if len(setWeights) == 0 {
		log.Fatalf("one or more weights is required")
	}

	tree := platecalc.NewTree(nil, float32(*barWeight))
	for _, perm := range platecalc.Permutations(plates...) {
		tree.Add(perm...)
	}

	result := platecalc.BestSolution(tree, setWeights, *maxDistance)
	if result == nil {
		log.Fatalf("no solution found")
		return
	}

	for _, node := range result {
		fmt.Printf("%3v: %v\n", node.TotalWeight(), node)
	}
}

func parsePlates() ([]float32, error) {
	plates := []float32{}
	for _, s := range strings.Split(*platesFlag, ",") {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		plates = append(plates, float32(f))
	}
	return plates, nil
}

func parseWeights() ([]int, error) {
	weights := []int{}
	for _, s := range flag.Args() {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		weights = append(weights, platecalc.RoundUpToNearest(float32(i), 5))
	}
	return weights, nil
}
