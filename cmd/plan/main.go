package main

import (
	"encoding/csv"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kdeloach/platecalc"
	"github.com/kdeloach/platecalc/plans"
	"gopkg.in/yaml.v3"
)

var barWeight = flag.Int("bar", 45, "bar weight")
var platesFlag = flag.String("plates", "45,35,25,10,10,5,5,2.5", "available plates")
var file = flag.String("file", "", "workout plan settings file")
var maxDistance = flag.Int("maxdistance", 5, "maximum distance to search tree")
var delim = flag.String("delim", ",", "output delimiter")
var debug = flag.Bool("debug", false, "display debug output")

func main() {
	flag.Parse()

	plates, err := parsePlates()
	if err != nil {
		log.Fatalf(err.Error())
	}

	tree := platecalc.NewTree(nil, float32(*barWeight))
	for _, perm := range platecalc.Permutations(plates...) {
		tree.Add(perm...)
	}

	buf, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf(err.Error())
	}

	settings := &plans.WorkoutPlanSettings{}
	err = yaml.Unmarshal(buf, settings)
	if err != nil {
		log.Fatalf(err.Error())
	}

	settings.PlateCalcFn = func(setWeights []int) []*platecalc.Tree {
		return platecalc.BestSolution(tree, setWeights, *maxDistance, *debug)
	}

	var plan plans.WorkoutPlan
	if settings.Plan == "Wendler531BBB" {
		plan = plans.NewWendler531BBB(settings)
	} else if settings.Plan == "Stronglifts" {
		plan = plans.NewStrongliftsPlan(settings)
	} else {
		log.Fatalf("unknown plan: %s (must be Wendler531BBB)\n", settings.Plan)
	}

	w := csv.NewWriter(os.Stdout)
	w.Comma = []rune(*delim)[0]
	plan.Write(w)
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
