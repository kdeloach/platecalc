package plans

import (
	"encoding/csv"
	"fmt"
	"log"

	"github.com/kdeloach/platecalc"
)

// Custom 531 is a modified Wendler 531 program with pyramid sets and increased
// volume. There are 3 heavy/low-rep sets, 3 light/high-rep sets, and 9
// moderate sets.
// Based on The '“Hypertrophy Rep Range” – Fact Or Fiction?' by Stronger By Science:
// https://www.strongerbyscience.com/hypertrophy-range-fact-fiction/
type custom531 struct {
	settings *WorkoutPlanSettings
}

type custom531PlanWriter struct {
	*csv.Writer
	plan *custom531
}

func NewCustom531(settings *WorkoutPlanSettings) *custom531 {
	return &custom531{
		settings: settings,
	}
}

func (plan *custom531) Write(w *csv.Writer) {
	pw := &custom531PlanWriter{
		Writer: w,
		plan:   plan,
	}

	pw.writeHeader()
	pw.writeWeek(1, []float32{0.50, 0.60, 0.80, 0.85, 0.40})
	pw.writeWeek(2, []float32{0.55, 0.65, 0.85, 0.90, 0.40})
	pw.writeWeek(3, []float32{0.60, 0.70, 0.90, 0.95, 0.40})
	pw.writeWeek(4, []float32{0.40, 0.50, 0.60, 0.70, 0.40})
}

func (pw *custom531PlanWriter) writeWeek(week int, tmPercs []float32) {
	pw.writeDay(SQUAT, week, 1, tmPercs)
	pw.writeDay(BENCH, week, 2, tmPercs)
	pw.writeDay(DEADLIFT, week, 3, tmPercs)
	pw.writeDay(PRESS, week, 4, tmPercs)
}

func (pw *custom531PlanWriter) writeDay(liftName string, week, day int, tmPercs []float32) {
	repMax := float32(pw.plan.settings.repMax(liftName))
	tmPerc := float32(pw.plan.settings.TrainingMaxPercent) / 100

	setWeights := []int{
		platecalc.FloorLimit(platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[0], 5), 45),
		platecalc.FloorLimit(platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[1], 5), 45),
		platecalc.FloorLimit(platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[2], 5), 45),
		platecalc.FloorLimit(platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[3], 5), 45),
		platecalc.FloorLimit(platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[4], 5), 45),
	}

	plates := pw.plan.settings.PlateCalcFn(setWeights)
	if plates == nil {
		log.Fatalf("no solution found for: %v setWeights=%v", liftName, setWeights)
	}

	sets := []int{5, 4, 2, 1, 3}
	reps := []int{8, 6, 5, 5, 15}

	pw.writeRow(liftName, week, day, tmPercs[0], setWeights[0], plates[0], sets[0], reps[0])
	pw.writeRow(liftName, week, day, tmPercs[1], setWeights[1], plates[1], sets[1], reps[1])
	pw.writeRow(liftName, week, day, tmPercs[2], setWeights[2], plates[2], sets[2], reps[2])
	pw.writeRow(liftName, week, day, tmPercs[3], setWeights[3], plates[3], sets[3], reps[3])
	pw.writeRow(liftName, week, day, tmPercs[4], setWeights[4], plates[4], sets[4], reps[4])
}

func (pw *custom531PlanWriter) writeHeader() {
	pw.Write([]string{
		"Lift", "Week", "Day", "TM %", "Weight", "Plates", "Sets", "Reps",
	})
	pw.Flush()
}

func (pw *custom531PlanWriter) writeRow(liftName string, week int, day int, tmPerc float32, weight int, plates *platecalc.Tree, sets int, reps int) {
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
