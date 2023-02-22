package main

import (
	"fmt"

	"github.com/GeorgianBadita/backgammon/pkg/board"
)

func main() {
	backgammonBoard := board.NewBoard(board.COLOR_WHITE)
	fmt.Println(backgammonBoard)
	moveRolls := backgammonBoard.GetValidMovesForDie(board.DieRoll{Die1: 6, Die2: 6})

	newBoard := moveRolls[0].MakeMoveRoll(backgammonBoard)
	newBoard.ColorToMove = board.COLOR_BLACK
	newMoves := newBoard.GetValidMovesForDie(board.DieRoll{Die1: 6, Die2: 1})

	for idx := 0; idx < len(newMoves); idx++ {
		moveRoll := newMoves[idx]
		fmt.Println("BOARD BEFORE MOVE ROLL: ")
		fmt.Println(newBoard)
		fmt.Println("BOARD AFTER MOVE ROLL: ")
		fmt.Println(moveRoll.MakeMoveRoll(newBoard))
	}
}
