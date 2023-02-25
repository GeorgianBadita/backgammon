package board

import "sort"

type MoveType int

const (
	NORMAL_MOVE         MoveType = 0
	CHECKER_ON_BAR_MOVE MoveType = 1
	BEARING_OFF_MOVE    MoveType = 2
)

type Move struct {
	From PointIndex
	To   PointIndex
	Type MoveType
}

// Funcion to apply a move to a given board
// NOTE: this function assumes the move is legal
func (m Move) MakeMove(b Board) Board {
	if b.ColorToMove != b.Points[m.From].Checker.Color {
		panic("Checker color to move is different than board player turn")
	}

	boardForMove := b.CopyBoard()
	if m.Type == NORMAL_MOVE || m.Type == CHECKER_ON_BAR_MOVE {
		checker := boardForMove.Points[m.From].Checker
		boardForMove.Points[m.From].CheckerCount -= 1
		// If the move leads to barring opponent's checkers
		if boardForMove.Points[m.To].CheckerCount == 1 && boardForMove.Points[m.To].Checker.Color != b.ColorToMove {
			// Increase checkers on bar index for color
			if boardForMove.Points[m.To].Checker.Color == COLOR_BLACK {
				boardForMove.Points[BLACK_PIECES_BAR_POINT_INDEX].CheckerCount += 1
			} else {
				boardForMove.Points[WHITE_PIECES_BAR_POINT_INDEX].CheckerCount += 1
			}
			boardForMove.Points[m.To].CheckerCount = 1
			boardForMove.Points[m.To].Checker = checker
		} else {
			boardForMove.Points[m.To].CheckerCount += 1
			boardForMove.Points[m.To].Checker = checker
		}
	} else if m.Type == BEARING_OFF_MOVE {
		if m.To != TO_INDEX_FOR_BEARING_OFF {
			panic("Bearing off move with destination is not valid!")
		}
		boardForMove.Points[m.From].CheckerCount -= 1
	}
	return boardForMove
}

type MoveRoll []Move

func (mvRoll MoveRoll) MakeMoveRoll(b Board) Board {
	boardForRoll := b.CopyBoard()
	for idx := 0; idx < len(mvRoll); idx++ {
		boardForRoll = mvRoll[idx].MakeMove(boardForRoll)
	}
	return boardForRoll
}

// This compares two move rolls, VERY naively
// NOTE: DO NOT USE, ONLY IN TESTS
// TODO: make this more efficient, use maps/something else, I don't want to
// do that now
func (mvRoll MoveRoll) isEqual(ot MoveRoll) bool {
	if len(mvRoll) != len(ot) {
		return false
	}

	sort.Slice(mvRoll, func(i, j int) bool {
		if mvRoll[i].From < mvRoll[j].From {
			return true
		} else if mvRoll[i].From > mvRoll[j].From {
			return false
		}

		if mvRoll[i].To < mvRoll[j].To {
			return true
		} else if mvRoll[i].To > mvRoll[j].To {
			return false
		}

		return mvRoll[i].Type == mvRoll[j].Type
	})

	sort.Slice(ot, func(i, j int) bool {
		if ot[i].From < ot[j].From {
			return true
		} else if ot[i].From > ot[j].From {
			return false
		}

		if ot[i].To < ot[j].To {
			return true
		} else if ot[i].To > ot[j].To {
			return false
		}

		return ot[i].Type == ot[j].Type
	})

	for idx := 0; idx < len(mvRoll); idx++ {
		if mvRoll[idx].From != ot[idx].From ||
			mvRoll[idx].To != ot[idx].To ||
			mvRoll[idx].Type != ot[idx].Type {
			return false
		}
	}
	return true
}

func areMoveRollListsEqual(curr []MoveRoll, ot []MoveRoll) bool {
	if len(curr) != len(ot) {
		return false
	}

	foundMatch := []bool{}
	for idx := 0; idx < len(curr); idx++ {
		foundMatch = append(foundMatch, false)
	}

	for idx := 0; idx < len(curr); idx++ {
		foundMatchForCurr := false
		for jdx := 0; jdx < len(ot); jdx++ {
			if curr[idx].isEqual(ot[jdx]) && !foundMatch[jdx] {
				foundMatchForCurr = true
				foundMatch[jdx] = true
				break
			}
		}
		if !foundMatchForCurr {
			return false
		}
	}

	for idx := 0; idx < len(curr); idx++ {
		if !foundMatch[idx] {
			return false
		}
	}

	return true
}
