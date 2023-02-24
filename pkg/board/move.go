package board

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

func (m Move) MakeMove(b Board) Board {
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
