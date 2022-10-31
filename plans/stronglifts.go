package plans

import (
	"encoding/csv"
	"fmt"
	"log"

	"github.com/kdeloach/platecalc"
)

type strongliftsPlan struct {
	settings *WorkoutPlanSettings
}

type strongliftsPlanWriter struct {
	*csv.Writer
	plan *strongliftsPlan
}

func NewStrongliftsPlan(settings *WorkoutPlanSettings) *strongliftsPlan {
	return &strongliftsPlan{
		settings: settings,
	}
}

func (plan *strongliftsPlan) Write(w *csv.Writer) {
	pw := &strongliftsPlanWriter{
		Writer: w,
		plan:   plan,
	}

	pw.writeHeader()
	pw.writeWeek(1)
	pw.writeWeek(2)
	pw.writeWeek(3)
	pw.writeWeek(4)
}

func (pw *strongliftsPlanWriter) writeWeek(week int) {
	if week%2 == 0 {
		pw.writeWorkoutB(week, 1)
		pw.writeWorkoutA(week, 2)
		pw.writeWorkoutB(week, 3)
	} else {
		pw.writeWorkoutA(week, 1)
		pw.writeWorkoutB(week, 2)
		pw.writeWorkoutA(week, 3)
	}
}

func (pw *strongliftsPlanWriter) writeWorkoutA(week, day int) {
	pw.writeLift("Squat", week, day)
	pw.writeLift("Bench", week, day)
	pw.writeLift("Deadlift", week, day)
}

func (pw *strongliftsPlanWriter) writeWorkoutB(week, day int) {
	pw.writeLift("Squat", week, day)
	pw.writeLift("Press", week, day)
	pw.writeLift("Deadlift", week, day)
}

func (pw *strongliftsPlanWriter) writeLift(liftName string, week, day int) {
	repMax := float32(pw.plan.settings.repMax(liftName))
	tmPercs := []float32{0.70, 0.80, 0.90, 1.00}
	tmPerc := float32(pw.plan.settings.TrainingMaxPercent) / 100

	load := 5
	if liftName == "Squat" || liftName == "Deadlift" {
		load = 10
	}

	// Squat Week 1, Day 1 = +0
	// Squat Week 1, Day 2 = +10
	// Squat Week 1, Day 3 = +20
	// Squat Week 2, Day 1 = +30
	// ...
	// multiply week by 3 because there are 3 workout days per week
	inc := float32(((week-1)*3 + (day - 1)) * load)

	setWeights := []int{
		platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[0]+inc, 5),
		platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[1]+inc, 5),
		platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[2]+inc, 5),
		platecalc.RoundUpToNearest(repMax*tmPerc*tmPercs[3]+inc, 5),
	}

	plates := pw.plan.settings.PlateCalcFn(setWeights)
	if plates == nil {
		log.Fatalf("no solution found for: %v setWeights=%v", liftName, setWeights)
	}

	pw.writeRow(liftName, week, day, tmPercs[0], setWeights[0], plates[0], 1, 5)
	pw.writeRow(liftName, week, day, tmPercs[1], setWeights[1], plates[1], 1, 5)
	pw.writeRow(liftName, week, day, tmPercs[2], setWeights[2], plates[2], 1, 5)

	// deadlifts reps should be 1x5 but I decided to leave it as 5x5 for
	// increased volume
	pw.writeRow(liftName, week, day, tmPercs[3], setWeights[3], plates[3], 5, 5)
}

func (pw *strongliftsPlanWriter) writeHeader() {
	pw.Write([]string{
		"Lift", "Week", "Day", "TM %", "Weight", "Plates", "Sets", "Reps",
	})
	pw.Flush()
}

func (pw *strongliftsPlanWriter) writeRow(liftName string, week int, day int, tmPerc float32, weight int, plates *platecalc.Tree, sets int, reps int) {
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
