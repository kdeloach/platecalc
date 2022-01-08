# platecalc

Command line utility to generate optimal workout plans for weightlifting
barbells.

## Commands

### calc

Calculate which plates to use and the optimal order to load them on the bar
to minimize the effort of changing plates during a workout.

Usage:

```sh
$ go run ./cmd/calc/ -h
Usage: calc [weight:int]+
  -bar int
        bar weight (default 45)
  -debug
        display debug output
  -maxdistance int
        maximum distance to search tree (default 5)
  -plates string
        available plates (default "45,35,25,10,10,5,5,2.5")
  -simple
        use simple plate orderings
```

For example, compare the following runs:

```sh
$ go run ./cmd/calc/ 100
100: 10, 10, 5, 2.5

$ go run ./cmd/calc/ 120
120: 25, 10, 2.5

$ go run ./cmd/calc/ 100 120
100: 25, 2.5
120: 25, 2.5, 10
```

Notice how the plates needed for `100` changes when followed by `120`. Also,
notice how the plate order for `120` changes when preceded by `100`.

Compare this to the naive plan which disables optimizations:

```sh
$ go run ./cmd/calc/ -simple 100 120
100: 10, 10, 5, 2.5
120: 25, 10, 2.5
```

Example FSL (first-set-last) workout sets:

```sh
$ go run ./cmd/calc/ 100 125 150 200 100
100: 25, 2.5
125: 25, 10, 5
150: 25, 10, 5, 2.5, 10
200: 25, 10, 5, 2.5, 35
100: 25, 2.5
```

Use the `-simple` flag to generate the simplest plate arrangement for each
target weight instead of calculating the optimal sequence:

```sh
$ go run ./cmd/calc/ -simple 100 125 150 200 100
100: 25, 2.5
125: 35, 5
150: 45, 5, 2.5
200: 45, 25, 5, 2.5
100: 25, 2.5
```

### plan

Generate workout plan based on [Jim Wendler's 5/3/1 BBB](https://www.jimwendler.com/blogs/jimwendler-com/101077382-boring-but-big)
program in CSV format. The plan is based on your 1RM (one-rep max) defined in a
setting file.

Plan details:
- 4 day split: Press, Deadlift, Bench, Squat
- 4 week cycle (3 weeks + deload week)
- training max (TM) is a percent of 1RM (`TrainingMaxPercent` in settings)
- week 1: 65/75/85/60 (percent of TM)
- week 2: 70/80/90/60
- week 3: 75/85/95/60
- week 4: 40/50/60/60
- reps are 5/3/1 based on week or 5/5/5 if `Progression5s` is enabled

Usage:

```sh
$ go run ./cmd/plan/ -h
Usage of /tmp/go-build1991805482/b001/exe/plan:
  -bar int
        bar weight (default 45)
  -debug
        display debug output
  -file string
        workout plan settings file
  -maxdistance int
        maximum distance to search tree (default 5)
  -plates string
        available plates (default "45,35,25,10,10,5,5,2.5")
```

Example:

```sh
$ go run ./cmd/plan/ -file profile.yaml
Lift;Week;Day;TM %;Weight;Plates;Sets;Reps
Press;1;1;65%;90;10, 10, 2.5;1;5
Press;1;1;75%;105;10, 10, 5, 5;1;5
Press;1;1;85%;115;10, 25;1;5
Press;1;1;60%;85;10, 10;5;10
Deadlift;1;2;65%;190;45, 10, 5, 10, 2.5;1;5
Deadlift;1;2;75%;220;45, 10, 5, 25, 2.5;1;5
Deadlift;1;2;85%;250;45, 10, 5, 25, 2.5, 10, 5;1;5
Deadlift;1;2;60%;175;45, 10, 5, 5;5;10
...
```

Format of `profile.yaml`:

```yaml
SquatRepMax: 300
DeadliftRepMax: 310
PressRepMax: 145
BenchRepMax: 205
TrainingMaxPercent: 90
Progression5s: true
```
