package ai

import "github.com/GeorgianBadita/backgammon/pkg/board"

type AI interface {
	ChooseMove(b board.Board, d board.DieRoll) board.MoveRoll
}
