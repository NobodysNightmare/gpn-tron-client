package main

import "fmt"

// The value of any cell that is not yet occupied by any player.
const EmptyCell = -1

type Direction string

const DirectionLeft = Direction("left")
const DirectionRight = Direction("right")
const DirectionUp = Direction("up")
const DirectionDown = Direction("down")

// World represents the current playing field and all information around it
// Its state will be updated through dedicated methods from the game loop.
// A decider can read from it to inform its decisions.
type World struct {
	Width, Height int
	Players       []Player // A list of all players connected to the game
	MyId          int      // The ID of the player that this client can control

	// The current state of the game grid. Each cell is either the
	// ID of the player that occupied it or EmptyCell if it is not yet occupied
	// by any player.
	Cells [][]int
}

// Player represents the actors that are part of the current playing round.
type Player struct {
	Id   int
	Name string
	Pos  Position // The position of the player's head
}

// A position on the playing field
type Position struct {
	X, Y int
}

// NewWorld creates a new world with the given width and height, as well as
// the given ID for the current player.
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

// Me returns the player that this client can control.
func (w World) Me() Player {
	return w.Players[w.MyId]
}

// NextCell calculates the position that's next to the given position
// in the given direction.
// Wrapping of world coordinates is taken into account.
func (w World) NextCell(pos Position, direction Direction) Position {
	if direction == DirectionLeft {
		return Position{w.wrapX(pos.X - 1), pos.Y}
	} else if direction == DirectionRight {
		return Position{w.wrapX(pos.X + 1), pos.Y}
	} else if direction == DirectionUp {
		return Position{pos.X, w.wrapY(pos.Y - 1)}
	}

	return Position{pos.X, w.wrapY(pos.Y + 1)}
}

// Neighbours returns the positions that are reachable from the given position
// within a single move.
func (w World) Neighbours(pos Position) []Position {
	return []Position{
		w.NextCell(pos, DirectionLeft),
		w.NextCell(pos, DirectionUp),
		w.NextCell(pos, DirectionRight),
		w.NextCell(pos, DirectionDown),
	}
}

func (w World) wrapX(x int) int {
	if x < 0 {
		x += w.Width
	}

	return x % w.Width
}

func (w World) wrapY(y int) int {
	if y < 0 {
		y += w.Height
	}

	return y % w.Height
}

// RegisterPlayer can be used by the game loop to add a new player to the world.
func (w *World) RegisterPlayer(playerId int, name string) {
	for playerId >= len(w.Players) {
		w.Players = append(w.Players, Player{})
	}

	w.Players[playerId].Id = playerId
	w.Players[playerId].Name = name
}

// UpdatePosition can be used by the game loop to update the known position of a given player.
func (w *World) UpdatePosition(playerId int, pos Position) {
	w.Cells[pos.X][pos.Y] = playerId
	w.Players[playerId].Pos = pos
}

// KillPlayer can be used by the game loop to remove a player from the world.
func (w *World) KillPlayer(playerId int) {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Cells[x][y] == playerId {
				w.Cells[x][y] = EmptyCell
			}
		}
	}
}

// PrettyPrint will print a text representation of the game world to the console.
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
