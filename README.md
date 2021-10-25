# platecalc

Command line utility to generate optimal workout plans for weightlifting
barbells.

The `calc` command outputs which plates to use and the optimal order to load
them on the bar to minimize effort across a series of target weights.

The `plan` command outputs a workout plan based on Jim Wendler's 5/3/1 BBB
program in CSV format.

## Commands

### calc

Calculate which plates to use and the order to load them on the bar. The
order is optimized to minimize effort to change plates during a workout.

Usage:

```sh
$ go run ./cmd/calc/ -h
Usage: calc [weight:int]+
  -bar int
        bar weight (default 45)
  -plates string
        available plates (default "45,35,25,10,10,5,5,2.5")
```

Example to calculate plates needed to lift `110`:

```sh
$ go run ./cmd/calc/ 110
110: 25, 5, 2.5  # 110 = (25+5+2.5)*2 + 45(bar)
```

Notice how the plate order for `110` shifts if it's preceded by `100` so you
only need to add 1 plate instead of removing 2 then adding 1:

```sh
$ go run ./cmd/calc/ 100 110
100: 25, 2.5
110: 25, 2.5, 5
```

Example FSL (first-set-last) workout sets:

```sh
$ go run ./cmd/calc/ 210 240 265 210
210: 45, 35, 2.5
240: 45, 35, 2.5, 10, 5
265: 45, 35, 25, 5
210: 45, 35, 2.5
```

### plan

Generate workout plan based on Jim Wendler's 5/3/1 BBB program in CSV format.
The plan is based on your 1RM (one-rep max) defined in a setting file.

Usage:

```sh
$ go run ./cmd/plan/ -h
Usage of /tmp/go-build3738305972/b001/exe/plan:
  -bar int
        bar weight (default 45)
  -file string
        workout plan settings file
  -plates string
        available plates (default "45,35,25,10,10,5,5,2.5")
```

Example:

```sh
$ go run ./cmd/plan/ -file settings.yaml
Lift;Week;Day;TM %;Weight;Plates;Sets;Reps
Press;1;1;65%;90;10, 10, 2.5;1;5
Press;1;1;75%;105;25, 5;1;5
Press;1;1;85%;115;35;1;5
Press;1;1;65%;90;10, 10, 2.5;5;5
Bench;1;1;60%;115;35;5;10
...
```

Format of `settings.yaml`:

```yaml
SquatRepMax: 300
DeadliftRepMax: 310
PressRepMax: 145
BenchRepMax: 205
```
