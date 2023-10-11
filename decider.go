package main

type Decider interface {
	DecideMove(World) string
}

type LeftDecider struct{}

func (LeftDecider) DecideMove(w World) string {
	return "left"
}
