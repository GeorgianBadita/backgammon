package game

import (
	"math/rand"

	"github.com/GeorgianBadita/backgammon/pkg/board"
)

type BackgammonGame struct {
	board             board.Board
	player1           IPlayer
	player2           IPlayer
	playerToMoveColor board.Color
}

func NewBackgammonGame(startingPlayerColor board.Color, pl1 IPlayer, pl2 IPlayer) BackgammonGame {
	if pl1.GetColor() == pl2.GetColor() {
		panic("Both players cannot have the same color")
	}
	board := board.NewBoard(startingPlayerColor)
	return BackgammonGame{board, pl1, pl2, startingPlayerColor}
}

func (bg BackgammonGame) MakeMoveForPlayer(mv board.MoveRoll, playerColor board.Color) BackgammonGame {
	if playerColor != bg.playerToMoveColor {
		panic("Move must be made by the player to move")
	}
	// Apply move
	newBoard := mv.MakeMoveRoll(bg.board)
	// Set the new color to move for board
	newBoard.ColorToMove = board.Color(1 - bg.playerToMoveColor)
	// Set the board to the new board
	bg.board = newBoard
	// Change current player to move
	bg.playerToMoveColor = board.Color(1 - bg.playerToMoveColor)
	return bg
}

func (bg BackgammonGame) IsGameOver() bool {
	return bg.board.ComputeGameState() == board.GAME_OVER
}

func (bg BackgammonGame) GetPossibleMoves(d board.DieRoll) []board.MoveRoll {
	return bg.board.GetValidMovesForDie(d)
}

func (bg BackgammonGame) GetDieRoll() board.DieRoll {
	return board.DieRoll{Die1: rand.Intn(6) + 1, Die2: rand.Intn(6) + 1}
}

func (bg BackgammonGame) GetPlayerToMove() board.Color {
	return bg.playerToMoveColor
}
