# Towers of Hanoi

The [Towers of Hanoi](https://en.wikipedia.org/wiki/Tower_of_Hanoi) is a mathematical
puzzle consisting of three rods and a number of disks of varying diameters, which can
be placed on any rod. The disks begin stacked on a single rod in order of decreasing
diameter. The goal is to move the disks to a different rod following these rules:

1. Only one disk may be moved at a time.
2. Each move consists of taking the topmost disk from one of the rods and moving it to
   a different rod.
3. No disk may be placed on top of a disk with a smaller diameter.

For three pegs, `n` disks, and a valid starting position as described above, the problem
of moving `n` disks from the starting peg to a goal peg lends itself to a simple
recursive solution. Given `m` disks on a source peg:

1. If `m` is zero, do nothing.
2. Otherwise, move `m - 1` disks from a source peg to a spare peg following these rules.
3. Move disk `m` to the target peg.
4. Move the `m - 1` disks from the spare peg to the target peg following these rules.

The full problem is solved with `m = n`, the starting peg as the source peg, the goal peg
as the target peg, and the spare peg as the spare peg.

To implement this in Go, we start with the usual package clause:

```go
package main

import "fmt"
```

To keep things easily testable, we'll define a function that implements the rules given
above. This function will move the disk labeled `m` from the source to the target peg
using the spare peg. The set of moves required to do so will be collected in `moves`.

```go
func move(moves *[]string, m int, source, target, spare string) {
```

If we're at the base case--i.e. no disks to move--just return.
```go
	if m == 0 {
		return
	}
```

Otherwise, move `m - 1` disks from the source to the spare, move disk `m`, to the target,
move `m - 1` disks from the spare to the target, and return.
```go
	move(moves, m - 1, source, spare, target)

	*moves = append(*moves, fmt.Sprintf("%v -> %v", source, target))

	move(moves, m - 1, spare, target, source)
}
```

Now that we have our solver, we can add a `main` function to drive it:

```go
func main() {
	var moves []string
	move(&moves, 10, "A", "B", "C")
	for _, m := range moves {
		fmt.Println(m)
	}
}
```
