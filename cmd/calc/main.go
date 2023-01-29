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
var platesFlag = flag.String("plates", "35,25,10,10,5,5,2.5,1.25", "available plates")
var maxDistance = flag.Int("maxdistance", 5, "maximum distance to search tree")
var debug = flag.Bool("debug", false, "display debug output")
var simple = flag.Bool("simple", false, "use simple plate orderings")
var preferLess = flag.Bool("less", false, "prefer less/heavier plates")

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

	opts := &platecalc.SolutionOpts{
		Debug:            *debug,
		PreferLessPlates: *preferLess,
	}

	var solution []*platecalc.Tree
	if *simple {
		solution = platecalc.SimpleSolution(tree, setWeights, opts)
	} else {
		solution = platecalc.BestSolution(tree, setWeights, *maxDistance, opts)
	}
	if solution == nil {
		log.Fatalf("no solution found")
		return
	}

	for _, node := range solution {
		fmt.Printf("%3v: %v\n", node.TotalWeight(), node)
	}
}

func parsePlates() ([]float32, error) {
	plates := []float32{}
	for _, s := range strings.Split(*platesFlag, ",") {
		n, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		plates = append(plates, float32(n))
	}
	return plates, nil
}

func parseWeights() ([]float32, error) {
	weights := []float32{}
	for _, s := range flag.Args() {
		n, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		weights = append(weights, float32(n))
	}
	return weights, nil
}
