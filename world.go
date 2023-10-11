package main

import "fmt"

const EmptyCell = -1

type World struct {
	Width, Height int
	Cells         [][]int // TODO: associate to players and remove upon player death
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
