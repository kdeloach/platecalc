package main

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/kdeloach/platecalc"
)

func calcWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		barWeight := 45
		less := false
		simple := false
		maxDistance := 5
		plates := []float32{45, 35, 25, 10, 10, 5, 5, 2.5}

		// Parse arguments
		setWeights := []int{}
		for _, arg := range args {
			if arg.Type() == js.TypeNumber {
				setWeights = append(setWeights, arg.Int())
			} else if arg.Type() == js.TypeObject {
				if v, err := tryGetBool(arg, "ordered"); err == nil {
					simple = v
				}
				if v, err := tryGetBool(arg, "less"); err == nil {
					less = v
				}
				if v, err := tryGetInt(arg, "barWeight"); err == nil {
					barWeight = v
				}
				if v, err := tryGetFloatArray(arg, "plates"); err == nil {
					plates = v
				}
			}
		}

		tree := platecalc.NewTree(nil, float32(barWeight))
		for _, perm := range platecalc.Permutations(plates...) {
			tree.Add(perm...)
		}

		opts := &platecalc.SolutionOpts{
			PreferLessPlates: less,
		}

		var solution []*platecalc.Tree
		if simple {
			solution = platecalc.SimpleSolution(tree, setWeights, opts)
		} else {
			solution = platecalc.BestSolution(tree, setWeights, maxDistance, opts)
		}

		if solution == nil {
			return map[string]interface{}{
				"error": "no solution",
			}
		}

		result := map[string]interface{}{}
		for i, node := range solution {
			k := fmt.Sprintf("set %d (%d)", i+1, node.TotalWeight())
			result[k] = node.String()
		}
		return result
	})
}

func tryGetBool(obj js.Value, key string) (v bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("key not defined in object")
		}
	}()
	return obj.Get(key).Bool(), nil
}

func tryGetInt(obj js.Value, key string) (v int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("key not defined in object")
		}
	}()
	return obj.Get(key).Int(), nil
}

func tryGetFloatArray(obj js.Value, key string) (ret []float32, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("key not defined in object")
		}
	}()
	arr := obj.Get(key)
	ret = []float32{}
	for i := 0; i < arr.Length(); i++ {
		n := float32(arr.Index(i).Float())
		ret = append(ret, n)
	}
	return
}

func main() {
	js.Global().Set("platecalc", calcWrapper())
	<-make(chan bool)
}
