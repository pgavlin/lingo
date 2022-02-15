//line hanoi.md:28
package main

import "fmt"
//line hanoi.md:38
func move(moves *[]string, m int, source, target, spare string) {
//line hanoi.md:43
	if m == 0 {
		return
	}
//line hanoi.md:51
	move(moves, m - 1, source, spare, target)

	*moves = append(*moves, fmt.Sprintf("%v -> %v", source, target))

	move(moves, m - 1, spare, target, source)
}
//line hanoi.md:62
func main() {
	var moves []string
	move(&moves, 10, "A", "B", "C")
	for _, m := range moves {
		fmt.Println(m)
	}
}
