package main

import (
	"encoding/csv"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/kdeloach/platecalc"
	"github.com/kdeloach/platecalc/plans"
	"gopkg.in/yaml.v3"
)

var barWeight = flag.Int("bar", 45, "bar weight")
var file = flag.String("file", "", "workout plan settings file")
var maxDistance = flag.Int("maxdistance", 5, "maximum distance to search tree")
var delim = flag.String("delim", ",", "output delimiter")
var debug = flag.Bool("debug", false, "display debug output")

func main() {
	flag.Parse()

	buf, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf(err.Error())
	}

	settings := &plans.WorkoutPlanSettings{}
	err = yaml.Unmarshal(buf, settings)
	if err != nil {
		log.Fatalf(err.Error())
	}

	plates, err := plans.ParsePlates(settings.Plates)
	if err != nil {
		log.Fatalf(err.Error())
	}

	tree := platecalc.NewTree(nil, float32(*barWeight))
	for _, perm := range platecalc.Permutations(plates...) {
		tree.Add(perm...)
	}

	opts := &platecalc.SolutionOpts{
		Debug:            *debug,
		PreferLessPlates: settings.PreferLessPlates,
	}

	settings.PlateCalcFn = func(setWeights []int) []*platecalc.Tree {
		return platecalc.BestSolution(tree, setWeights, *maxDistance, opts)
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
