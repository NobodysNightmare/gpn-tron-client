package main

type Decider interface {
	DecideMove(World) string
}

type LeftDecider struct{}

func (LeftDecider) DecideMove(w World) string {
	return "left"
}

type LongPathDecider struct{}

func (LongPathDecider) DecideMove(w World) string {
	choice := DirectionLeft
	maxLength := PathLength(w, w.Me().X, w.Me().Y, "left")
	for _, currentDir := range []string{DirectionRight, DirectionUp, DirectionDown} {
		currentLength := PathLength(w, w.Me().X, w.Me().Y, currentDir)
		if currentLength > maxLength {
			choice = currentDir
		}
	}

	return choice
}

type HighScoreDecider struct{}

func (HighScoreDecider) DecideMove(w World) string {
	choice := DirectionLeft
	maxScore := PathScore(w, w.Me().X, w.Me().Y, "left")
	for _, currentDir := range []string{DirectionRight, DirectionUp, DirectionDown} {
		currentScore := PathScore(w, w.Me().X, w.Me().Y, currentDir)
		if currentScore > maxScore {
			choice = currentDir
		}
	}

	return choice
}

func PathLength(w World, x, y int, direction string) int {
	length := 0
	x, y = w.NextCell(x, y, direction)
	for w.Cells[x][y] == EmptyCell {
		length += 1
		x, y = w.NextCell(x, y, direction)
	}

	return length
}

func PathScore(w World, x, y int, direction string) int {
	score := 0
	x, y = w.NextCell(x, y, direction)
	for w.Cells[x][y] == EmptyCell {
		score += 1
		// TODO: score heads of others negatively
		// TODO: score free space positively
		x, y = w.NextCell(x, y, direction)
	}

	return score
}
