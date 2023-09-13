package main

import (
	"expert_systems/pkg/models/stack"
	"log"
)

func main() {
	stck := stack.NewStack[string]()
	stck.Push("1")
	stck2 := stack.NewStack[[]int]()
	stck2.Push([]int{1, 2, 3, 4, 5})
	element, _ := stck2.Peek()

	log.Println(element)
}
