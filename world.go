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
	Pos  Position
}

type Position struct {
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

func (w World) NextCell(pos Position, direction string) Position {
	if direction == DirectionLeft {
		return Position{w.WrapX(pos.X - 1), pos.Y}
	} else if direction == DirectionRight {
		return Position{w.WrapX(pos.X + 1), pos.Y}
	} else if direction == DirectionUp {
		return Position{pos.X, w.WrapY(pos.Y - 1)}
	}

	return Position{pos.X, w.WrapY(pos.Y + 1)}
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

func (w *World) UpdatePosition(playerId int, pos Position) {
	w.Cells[pos.X][pos.Y] = playerId
	w.Players[playerId].Pos = pos
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

	pos := w.Me().Pos

	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Cells[x][y] == EmptyCell {
				fmt.Print(" ")
			} else if pos.X == x && pos.Y == y {
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
