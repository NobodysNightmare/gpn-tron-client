package main

import "fmt"

const EmptyCell = -1

const DirectionLeft = "left"
const DirectionRight = "right"
const DirectionUp = "up"
const DirectionDown = "down"

type World struct {
	Width, Height int
	Cells         [][]int
	Players       []Player
	MyId          int
}

type Player struct {
	Name string
	X, Y int
}

func NewWorld(w, h, playerId int) World {
	world := World{Width: w, Height: h, Players: []Player{}, MyId: playerId}

	world.Cells = make([][]int, w)
	for x := 0; x < w; x++ {
		world.Cells[x] = make([]int, h)
		for y := 0; y < h; y++ {
			world.Cells[x][y] = EmptyCell
		}
	}

	return world
}

func (w World) Me() Player {
	return w.Players[w.MyId]
}

func (w World) NextCell(x, y int, direction string) (int, int) {
	if direction == DirectionLeft {
		return w.WrapX(x - 1), y
	} else if direction == DirectionRight {
		return w.WrapX(x + 1), y
	} else if direction == DirectionUp {
		return x, w.WrapY(y - 1)
	}

	return x, w.WrapY(y + 1)
}

func (w World) WrapX(x int) int {
	if x < 0 {
		x += w.Width
	}

	return x % w.Width
}

func (w World) WrapY(y int) int {
	if y < 0 {
		y += w.Height
	}

	return y % w.Height
}

func (w *World) RegisterPlayer(playerId int, name string) {
	for playerId >= len(w.Players) {
		w.Players = append(w.Players, Player{})
	}

	w.Players[playerId].Name = name
}

func (w *World) UpdatePosition(playerId, x, y int) {
	w.Cells[x][y] = playerId
	w.Players[playerId].X = x
	w.Players[playerId].Y = y
}

func (w *World) KillPlayer(playerId int) {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Cells[x][y] == playerId {
				w.Cells[x][y] = EmptyCell
			}
		}
	}
}

func (w World) PrettyPrint() {
	for x := 0; x < w.Width; x++ {
		fmt.Print("=")
	}
	fmt.Println()

	myX := w.Me().X
	myY := w.Me().Y

	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Cells[x][y] == EmptyCell {
				fmt.Print(" ")
			} else if myX == x && myY == y {
				fmt.Print("👑")
			} else {
				fmt.Print("X")
			}
		}
		fmt.Println()
	}

	for x := 0; x < w.Width; x++ {
		fmt.Print("=")
	}
	fmt.Println()
}
