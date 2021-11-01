package plans

import (
	"encoding/csv"
	"fmt"
	"log"

	"github.com/kdeloach/platecalc"
)

type WorkoutPlanSettings struct {
	SquatRepMax    int `yaml:"SquatRepMax"`
	DeadliftRepMax int `yaml:"DeadliftRepMax"`
	PressRepMax    int `yaml:"PressRepMax"`
	BenchRepMax    int `yaml:"BenchRepMax"`
}

type wendler531BBB struct {
	tree        *platecalc.Tree
	settings    *WorkoutPlanSettings
	maxDistance int
	debug       bool
}

func NewWendler531BBB(tree *platecalc.Tree, settings *WorkoutPlanSettings, maxDistance int, debug bool) *wendler531BBB {
	return &wendler531BBB{
		tree:        tree,
		settings:    settings,
		maxDistance: maxDistance,
		debug:       debug,
	}
}

func (w *wendler531BBB) Write(writer *csv.Writer) {
	tms := [][]float32{
		{.65, .75, .85},
		{.7, .8, .9},
		{.75, .85, .95},
	}
	reps := [][]int{
		{5, 5, 5},
		{3, 3, 3},
		{5, 3, 1},
	}

	lifts1 := []string{"Press", "Deadlift", "Bench", "Squat"}

	// Calculate training max from 90% of 1 rep max
	weights := map[string]int{
		"Squat":    platecalc.RoundUpToNearest(float32(w.settings.SquatRepMax)*.9, 5),
		"Deadlift": platecalc.RoundUpToNearest(float32(w.settings.DeadliftRepMax)*.9, 5),
		"Press":    platecalc.RoundUpToNearest(float32(w.settings.PressRepMax)*.9, 5),
		"Bench":    platecalc.RoundUpToNearest(float32(w.settings.BenchRepMax)*.9, 5),
	}

	// Output header
	writer.Write([]string{
		"Lift", "Week", "Day", "TM %", "Weight", "Plates", "Sets", "Reps",
	})
	writer.Flush()

	for week := 0; week < 3; week++ {
		for day := 0; day < 4; day++ {
			lift1 := lifts1[day]

			w1 := weights[lift1]

			// Find optimal sequence of plate changes for 5/3/1 lift
			setWeights := make([]int, 0)
			for i := 0; i < 4; i++ {
				tm := tms[week][i]
				wt := platecalc.RoundUpToNearest(float32(w1)*tm, 5)
				setWeights = append(setWeights, wt)
			}

			plates := platecalc.BestSolution(w.tree, setWeights, w.maxDistance, w.debug)
			if plates == nil {
				log.Fatalf("no solution found for: %v weight=%v", lift1, w1)
			}

			// Output rows for 5/3/1 lift
			for i := 0; i < 3; i++ {
				r := reps[week][i]
				tm := tms[week][i]
				ps := plates[i]
				wt := setWeights[i]
				writeRow(writer, lift1, week, day, tm, wt, ps, 1, r)
			}

			// Output rows for 5x10 lift
			tm := float32(0.6)
			wt := platecalc.RoundUpToNearest(float32(w1)*tm, 5)
			ps := platecalc.BestSolution(w.tree, []int{wt}, w.maxDistance, w.debug)
			if plates == nil {
				log.Fatalf("no solution found for: %v weight=%v", lift1, w2)
			}
			writeRow(writer, lift1, week, day, tm, wt, ps[0], 5, 10)
		}
	}
}

func writeRow(writer *csv.Writer, lift string, week int, day int, tm float32, wt int, plates *platecalc.Tree, sets int, reps int) {
	writer.Write([]string{
		lift,
		fmt.Sprintf("%v", week+1),
		fmt.Sprintf("%v", day+1),
		fmt.Sprintf("%v%%", int(tm*100)),
		fmt.Sprintf("%v", wt),
		plates.String(),
		fmt.Sprintf("%v", sets),
		fmt.Sprintf("%v", reps),
	})
	writer.Flush()
}
