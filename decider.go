package main

type Decider interface {
	DecideMove(World) Direction
}

type CollisionAvoidDecider struct{}

func (CollisionAvoidDecider) DecideMove(w World) Direction {
	for _, currentDir := range []Direction{DirectionLeft, DirectionUp, DirectionRight, DirectionDown} {
		pos := w.NextCell(w.Me().Pos, currentDir)
		if w.Cells[pos.X][pos.Y] == EmptyCell {
			return currentDir
		}
	}

	return DirectionLeft
}

type LongPathDecider struct{}

func (LongPathDecider) DecideMove(w World) Direction {
	choice := DirectionLeft
	maxLength := PathLength(w, w.Me().Pos, "left")
	for _, currentDir := range []Direction{DirectionUp, DirectionRight, DirectionDown} {
		currentLength := PathLength(w, w.Me().Pos, currentDir)
		if currentLength > maxLength {
			choice = currentDir
		}
	}

	return choice
}

type HighScoreDecider struct{}

const playerHeadPenalty = 5

func (HighScoreDecider) DecideMove(w World) Direction {
	choice := DirectionLeft
	maxScore := PathScore(w, w.Me(), "left")
	for _, currentDir := range []Direction{DirectionUp, DirectionRight, DirectionDown} {
		currentScore := PathScore(w, w.Me(), currentDir)
		if currentScore > maxScore {
			choice = currentDir
		}
	}

	return choice
}

func PathLength(w World, pos Position, direction Direction) int {
	length := 0
	pos = w.NextCell(pos, direction)
	for w.Cells[pos.X][pos.Y] == EmptyCell {
		length += 1
		pos = w.NextCell(pos, direction)
	}

	return length
}

func PathScore(w World, p Player, direction Direction) int {
	score := 0
	pos := w.NextCell(p.Pos, direction)

	for w.Cells[pos.X][pos.Y] == EmptyCell {
		score += 1
		// TODO: score free space positively
		pos = w.NextCell(pos, direction)
	}

	// deduct points for player heads next to the target (unpredictable collision risk)
	for _, neighbourPos := range w.Neighbours(pos) {
		for _, player := range w.Players {
			if player.Pos == neighbourPos && player.Id != w.MyId {
				// ensure that steering next to player head is always better than immediate collision
				penalty := max(min(score-1, playerHeadPenalty), 0)
				score -= penalty
			}
		}
	}

	return score
}
