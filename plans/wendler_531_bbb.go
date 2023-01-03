package plans

import (
	"encoding/csv"
	"fmt"
	"log"

	"github.com/kdeloach/platecalc"
)

type wendler531BBB struct {
	settings *WorkoutPlanSettings
}

type wendler531BBBPlanWriter struct {
	*csv.Writer
	plan *wendler531BBB
}

func NewWendler531BBB(settings *WorkoutPlanSettings) *wendler531BBB {
	return &wendler531BBB{
		settings: settings,
	}
}

func (plan *wendler531BBB) Write(w *csv.Writer) {
	pw := &wendler531BBBPlanWriter{
		Writer: w,
		plan:   plan,
	}

	pw.writeHeader()
	pw.writeWeek(1, []float32{0.65, 0.75, 0.85, 0.60})
	pw.writeWeek(2, []float32{0.70, 0.80, 0.90, 0.60})
	pw.writeWeek(3, []float32{0.75, 0.85, 0.95, 0.60})
	pw.writeWeek(4, []float32{0.50, 0.60, 0.70, 0.60})
}

func (pw *wendler531BBBPlanWriter) writeWeek(week int, tmPercs []float32) {
	pw.writeDay(SQUAT, week, 1, tmPercs)
	pw.writeDay(BENCH, week, 2, tmPercs)
	pw.writeDay(DEADLIFT, week, 3, tmPercs)
	pw.writeDay(PRESS, week, 4, tmPercs)
}

func (pw *wendler531BBBPlanWriter) writeDay(liftName string, week, day int, tmPercs []float32) {
	repMax := float32(pw.plan.settings.repMax(liftName))
	tmPerc := float32(pw.plan.settings.TrainingMaxPercent) / 100

	setWeights := []int{
		platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[0], 5),
		platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[1], 5),
		platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[2], 5),
		platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[3], 5),
		platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[4], 5),
	}

	plates := pw.plan.settings.PlateCalcFn(setWeights)
	if plates == nil {
		log.Fatalf("no solution found for: %v setWeights=%v", liftName, setWeights)
	}

	var reps []int

	if pw.plan.settings.Progression5s {
		reps = []int{5, 5, 5}
	} else {
		if week == 2 {
			reps = []int{3, 3, 3}
		} else if week == 3 {
			reps = []int{5, 3, 1}
		} else {
			reps = []int{5, 5, 5}
		}
	}

	// Wendler 531 main lifts
	pw.writeRow(liftName, week, day, tmPercs[0], setWeights[0], plates[0], 1, reps[0])
	pw.writeRow(liftName, week, day, tmPercs[1], setWeights[1], plates[1], 1, reps[1])
	pw.writeRow(liftName, week, day, tmPercs[2], setWeights[2], plates[2], 1, reps[2])

	// Wendler BBB 5x10 supplemental lift
	pw.writeRow(liftName, week, day, tmPercs[3], setWeights[3], plates[3], 5, 10)
}

func (pw *wendler531BBBPlanWriter) writeHeader() {
	pw.Write([]string{
		"Lift", "Week", "Day", "TM %", "Weight", "Plates", "Sets", "Reps",
	})
	pw.Flush()
}

func (pw *wendler531BBBPlanWriter) writeRow(liftName string, week int, day int, tmPerc float32, weight int, plates *platecalc.Tree, sets int, reps int) {
	pw.Write([]string{
		liftName,
		fmt.Sprintf("%v", week),
		fmt.Sprintf("%v", day),
		fmt.Sprintf("%v%%", int(tmPerc*100)),
		fmt.Sprintf("%v", weight),
		plates.String(),
		fmt.Sprintf("%v", sets),
		fmt.Sprintf("%v", reps),
	})
	pw.Flush()
}
