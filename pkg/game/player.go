package game

import (
	"github.com/GeorgianBadita/backgammon-move-generator/pkg/ai"
	"github.com/GeorgianBadita/backgammon-move-generator/pkg/board"
)

type IPlayer interface {
	GetColor() board.Color
}

type IAIPlayer interface {
	GetColor() board.Color
	GetMove(b board.Board, d board.DieRoll) board.MoveRoll
}

type HumanPlayer struct {
	Color board.Color
	Name  string
}

func (hp HumanPlayer) GetColor() board.Color {
	return hp.Color
}

type AIPlayer struct {
	Color board.Color
	AI    ai.AI
}

func (ai AIPlayer) GetColor() board.Color {
	return ai.Color
}

func (ai AIPlayer) GetMove(b board.Board, d board.DieRoll) board.MoveRoll {
	return ai.AI.ChooseMove(b, d)
}
