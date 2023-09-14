package vertex

type Vertex struct {
	Number  int
	Feature int
}

func NewVertex(number int) *Vertex {
	return &Vertex{
		Number:  number,
		Feature: 0,
	}
}
