package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/GeorgianBadita/backgammon-move-generator/pkg/ai"
	"github.com/GeorgianBadita/backgammon-move-generator/pkg/board"
)

func main() {
	f, err := os.Create("backgammon.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	backgammonBoard := board.NewBoard(board.COLOR_WHITE)
	dieRoll := board.DieRoll{Die1: 6, Die2: 6}

	bestMovRoll := ai.Minimax(backgammonBoard, 4, board.COLOR_WHITE, dieRoll)

	fmt.Println(backgammonBoard)
	newBoard := bestMovRoll.MakeMoveRoll(backgammonBoard)
	fmt.Println(newBoard)
	// fmt.Println(ai.EvaluateBoard(backgammonBoard))

	// fmt.Println(backgammonBoard)
	// moveRolls := backgammonBoard.GetValidMovesForDieRoll(board.DieRoll{Die1: 6, Die2: 6})

	// newBoard := moveRolls[0].MakeMoveRoll(backgammonBoard)
	// newBoard.ColorToMove = board.COLOR_BLACK
	// fmt.Println(ai.EvaluateBoard(newBoard))
	// newMoves := newBoard.GetValidMovesForDieRoll(board.DieRoll{Die1: 6, Die2: 1})

	// for idx := 0; idx < len(newMoves); idx++ {
	// 	moveRoll := newMoves[idx]
	// 	fmt.Println("BOARD BEFORE MOVE ROLL: ")
	// 	fmt.Println(newBoard)
	// 	fmt.Println("BOARD AFTER MOVE ROLL: ")
	// 	fmt.Println(moveRoll.MakeMoveRoll(newBoard))
	// }
}
