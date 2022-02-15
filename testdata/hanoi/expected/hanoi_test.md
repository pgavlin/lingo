# Tests

This file contains some tests for our solver. These tests use a common driver that accepts
the same parameters as `move` plus the expected set of moves and validates that the
returned moves match what is expected.

```go
package main

import "testing"

func testMove(t *testing.T, n int, source, target, spare string, expected ...string) {
	var moves []string
	move(&moves, n, source, target, spare)

	if len(moves) != len(expected) {
		t.Fatalf("len(moves) != len(expected): %v != %v", len(moves), len(expected))
	}

	for i := range expected {
		if moves[i] != expected[i] {
			t.Fatalf("moves[%v] != expected[%v]: '%v' != '%v'", i, i, moves[i], expected[i])
		}
	}
} 
```

We'll add a few tests for various values of `n`, `source`, `target`, and `spare`.

```go
func Test3ABC(t *testing.T) {
	testMove(t, 3, "A", "B", "C",
		"A -> B",
		"A -> C",
		"B -> C",
		"A -> B",
		"C -> A",
		"C -> B",
		"A -> B")
}

func Test5ABC(t *testing.T) {
	testMove(t, 5, "A", "B", "C",
		"A -> B",
		"A -> C",
		"B -> C",
		"A -> B",
		"C -> A",
		"C -> B",
		"A -> B",
		"A -> C",
		"B -> C",
		"B -> A",
		"C -> A",
		"B -> C",
		"A -> B",
		"A -> C",
		"B -> C",
		"A -> B",
		"C -> A",
		"C -> B",
		"A -> B",
		"C -> A",
		"B -> C",
		"B -> A",
		"C -> A",
		"C -> B",
		"A -> B",
		"A -> C",
		"B -> C",
		"A -> B",
		"C -> A",
		"C -> B",
		"A -> B")
}

func Test4XYZ(t *testing.T) {
	testMove(t, 4, "X", "Y", "Z",
		"X -> Z",
		"X -> Y",
		"Z -> Y",
		"X -> Z",
		"Y -> X",
		"Y -> Z",
		"X -> Z",
		"X -> Y",
		"Z -> Y",
		"Z -> X",
		"Y -> X",
		"Z -> Y",
		"X -> Z",
		"X -> Y",
		"Z -> Y")
}
```
