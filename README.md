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
program in CSV format. The plan is based on your 1RM (one-rep max) defined in a setting file.

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
$ go run ./cmd/plan/ -file settings.yaml
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

Format of `settings.yaml`:

```yaml
SquatRepMax: 300
DeadliftRepMax: 310
PressRepMax: 145
BenchRepMax: 205
```
