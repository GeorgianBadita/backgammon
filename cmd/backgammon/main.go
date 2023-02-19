package main

import (
	"fmt"

	"github.com/GeorgianBadita/backgammon/pkg/board"
)

func main() {
	backgammonBoard := board.NewBoard(board.COLOR_WHITE)
	fmt.Println(backgammonBoard)
}
