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

func (m Move) makeMove(b Board) Board {
	return b
}

type MoveRoll []Move
