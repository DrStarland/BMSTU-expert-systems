package node

type Node struct {
	Number  int
	Feature int
}

func NewNode(number int) *Node {
	return &Node{
		Number:  number,
		Feature: 0,
	}
}
