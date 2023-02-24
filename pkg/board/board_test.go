package board

import (
	"fmt"
	"testing"
)

type gameStateTest struct {
	board             Board
	expectedGameState GameState
}

func TestComputeGameState(t *testing.T) {
	for _, test := range makeGameStateTests() {
		if output := test.board.ComputeGameState(); output != test.expectedGameState {
			fmt.Println(test.board)
			t.Errorf("Output %q not equal to expected %q", output, test.expectedGameState)
		}
	}
}

func makeGameStateTests() []gameStateTest {
	// test 1 - initial board game state is normal play
	board := NewBoard(COLOR_BLACK)
	expectedGameState := NORMAL_PLAY

	// test 2 - white has pieces on bar
	board1 := NewBoard(COLOR_WHITE)
	board1.Points[7].CheckerCount -= 1
	board1.Points[WHITE_PIECES_BAR_POINT_INDEX].CheckerCount += 1
	expectedGameState1 := CHECKERS_ON_BAR

	// test 3 - white doesn't have any piece left outside of home
	board2 := NewBoard(COLOR_WHITE)
	board2.Points[7].CheckerCount = 0
	board2.Points[12].CheckerCount = 0
	board2.Points[23].CheckerCount = 0
	expectedGameState2 := BEARING_OFF

	// test 4 - black doesn't have any piece left
	board3 := NewBoard(COLOR_BLACK)
	board3.Points[18].CheckerCount = 0
	board3.Points[16].CheckerCount = 0
	board3.Points[11].CheckerCount = 0
	board3.Points[0].CheckerCount = 0
	expectedGameState3 := GAME_OVER

	return []gameStateTest{
		{board, expectedGameState},
		{board1, expectedGameState1},
		{board2, expectedGameState2},
		{board3, expectedGameState3},
	}
}
