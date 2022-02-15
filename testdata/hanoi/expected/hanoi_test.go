//line hanoi_test.md:8
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
//line hanoi_test.md:31
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
