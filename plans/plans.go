package plans

import (
	"encoding/csv"
	"log"

	"github.com/kdeloach/platecalc"
)

type WorkoutPlan interface {
	Write(w *csv.Writer)
}

type WorkoutPlanSettings struct {
	Plan               string `yaml:"Plan"`
	SquatRepMax        int    `yaml:"SquatRepMax"`
	DeadliftRepMax     int    `yaml:"DeadliftRepMax"`
	PressRepMax        int    `yaml:"PressRepMax"`
	BenchRepMax        int    `yaml:"BenchRepMax"`
	TrainingMaxPercent int    `yaml:"TrainingMaxPercent"`
	Progression5s      bool   `yaml:"Progression5s"`
	PlateCalcFn        PlateCalcFunction
}

type PlateCalcFunction func(setWeights []int) []*platecalc.Tree

func (settings *WorkoutPlanSettings) repMax(liftName string) int {
	switch liftName {
	case "Squat":
		return settings.SquatRepMax
	case "Deadlift":
		return settings.DeadliftRepMax
	case "Press":
		return settings.PressRepMax
	case "Bench":
		return settings.BenchRepMax
	}
	log.Fatalf("unknown lift name: %v", liftName)
	return 0
}