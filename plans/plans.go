package plans

import (
	"encoding/csv"
	"log"
	"strconv"
	"strings"

	"github.com/kdeloach/platecalc"
)

var (
	SQUAT    = "Squat"
	DEADLIFT = "Deadlift"
	BENCH    = "Bench"
	PRESS    = "Press"
)

type WorkoutPlan interface {
	Write(w *csv.Writer)
}

type WorkoutPlanSettings struct {
	Plan               string `yaml:"Plan"`
	Plates             string `yaml:"Plates"`
	SquatRepMax        int    `yaml:"SquatRepMax"`
	DeadliftRepMax     int    `yaml:"DeadliftRepMax"`
	PressRepMax        int    `yaml:"PressRepMax"`
	BenchRepMax        int    `yaml:"BenchRepMax"`
	TrainingMaxPercent int    `yaml:"TrainingMaxPercent"`
	Progression5s      bool   `yaml:"Progression5s"`
	PreferLessPlates   bool   `yaml:"PreferLessPlates"`
	PlateCalcFn        PlateCalcFunction
}

type PlateCalcFunction func(setWeights []int) []*platecalc.Tree

func (settings *WorkoutPlanSettings) repMax(liftName string) int {
	switch liftName {
	case SQUAT:
		return settings.SquatRepMax
	case DEADLIFT:
		return settings.DeadliftRepMax
	case PRESS:
		return settings.PressRepMax
	case BENCH:
		return settings.BenchRepMax
	}
	log.Fatalf("unknown lift name: %v", liftName)
	return 0
}

func ParsePlates(strPlates string) ([]float32, error) {
	plates := []float32{}
	for _, s := range strings.Split(strPlates, ",") {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		plates = append(plates, float32(f))
	}
	return plates, nil
}
