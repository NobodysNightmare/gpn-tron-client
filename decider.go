package main

type Decider interface {
	DecideMove(World) string
}

type CollisionAvoidDecider struct{}

func (CollisionAvoidDecider) DecideMove(w World) string {
	for _, currentDir := range []string{DirectionLeft, DirectionRight, DirectionUp, DirectionDown} {
		pos := w.NextCell(w.Me().Pos, currentDir)
		if w.Cells[pos.X][pos.Y] == EmptyCell {
			return currentDir
		}
	}

	return DirectionLeft
}

type LongPathDecider struct{}

func (LongPathDecider) DecideMove(w World) string {
	choice := DirectionLeft
	maxLength := PathLength(w, w.Me().Pos, "left")
	for _, currentDir := range []string{DirectionRight, DirectionUp, DirectionDown} {
		currentLength := PathLength(w, w.Me().Pos, currentDir)
		if currentLength > maxLength {
			choice = currentDir
		}
	}

	return choice
}

type HighScoreDecider struct{}

func (HighScoreDecider) DecideMove(w World) string {
	choice := DirectionLeft
	maxScore := PathScore(w, w.Me().Pos, "left")
	for _, currentDir := range []string{DirectionRight, DirectionUp, DirectionDown} {
		currentScore := PathScore(w, w.Me().Pos, currentDir)
		if currentScore > maxScore {
			choice = currentDir
		}
	}

	return choice
}

func PathLength(w World, pos Position, direction string) int {
	length := 0
	pos = w.NextCell(pos, direction)
	for w.Cells[pos.X][pos.Y] == EmptyCell {
		length += 1
		pos = w.NextCell(pos, direction)
	}

	return length
}

func PathScore(w World, pos Position, direction string) int {
	score := 0
	pos = w.NextCell(pos, direction)
	for w.Cells[pos.X][pos.Y] == EmptyCell {
		score += 1
		// TODO: score heads of others negatively
		// TODO: score free space positively
		pos = w.NextCell(pos, direction)
	}

	return score
}
