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
	pw.writeLift(SQUAT, week, day)
	pw.writeLift(BENCH, week, day)
	pw.writeLift(DEADLIFT, week, day)
}

func (pw *strongliftsPlanWriter) writeWorkoutB(week, day int) {
	pw.writeLift(SQUAT, week, day)
	pw.writeLift(PRESS, week, day)
	pw.writeLift(DEADLIFT, week, day)
}

func (pw *strongliftsPlanWriter) writeLift(liftName string, week, day int) {
	repMax := float32(pw.plan.settings.repMax(liftName))
	tmPerc := float32(pw.plan.settings.TrainingMaxPercent) / 100

	// increase weight for bench and press half as much since they
	// only appear on alternating days
	var rate float32 = 1.0
	if liftName == BENCH || liftName == PRESS {
		rate = 0.5
	}

	// increase weight by 5 pounds per day
	inc := float32(float32((week-1)*3+(day-1)) * (5 * rate))

	setWeights := []int{
		platecalc.RoundUpToNearest(repMax*tmPerc+inc, 5),
	}

	plates := pw.plan.settings.PlateCalcFn(setWeights)
	if plates == nil {
		log.Fatalf("no solution found for: %v setWeights=%v", liftName, setWeights)
	}

	if liftName == DEADLIFT {
		pw.writeRow(liftName, week, day, tmPerc, setWeights[0], plates[0], 1, 5)
	} else {
		pw.writeRow(liftName, week, day, tmPerc, setWeights[0], plates[0], 5, 5)
	}
}

func (pw *strongliftsPlanWriter) writeHeader() {
	pw.Write([]string{
		"Week", "Day", "Lift", "Weight", "Plates", "Sets", "Reps", "TM %",
	})
	pw.Flush()
}

func (pw *strongliftsPlanWriter) writeRow(liftName string, week int, day int, tmPerc float32, weight int, plates *platecalc.Tree, sets int, reps int) {
	pw.Write([]string{
		fmt.Sprintf("%v", week),
		fmt.Sprintf("%v", day),
		liftName,
		fmt.Sprintf("%v", weight),
		plates.String(),
		fmt.Sprintf("%v", sets),
		fmt.Sprintf("%v", reps),
		fmt.Sprintf("%v%%", int(tmPerc*100)),
	})
	pw.Flush()
}
