package game

import (
	"github.com/GeorgianBadita/backgammon-move-generator/pkg/board"
)

// Function that applies a move to a serialized board string
func MakeMoveOnSerializedBoard(boardString string, mv board.Move) string {
	b := board.DeserializeBoard(boardString)
	newBoard := mv.MakeMove(b)
	newBoard.ColorToMove = board.Color(1 - b.ColorToMove)
	return newBoard.SerializeBoard()
}

// Function that applies a move roll to a serialized board string
func MakeMoveRollOnSerializedBoard(boardString string, mvRoll board.MoveRoll) string {
	b := board.DeserializeBoard(boardString)
	for idx := 0; idx < len(mvRoll); idx++ {
		b = mvRoll[idx].MakeMove(b)
	}
	b.ColorToMove = board.Color(1 - b.ColorToMove)
	return b.SerializeBoard()
}

// Function that gets valid moves for a serialized board and die roll
func GetMovesForSerializedBoard(boardString string, dieRoll board.DieRoll) []board.MoveRoll {
	board := board.DeserializeBoard(boardString)
	return board.GetValidMovesForDie(dieRoll)
}
