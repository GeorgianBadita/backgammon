package board

import (
	"fmt"
	"testing"
)

type gameStateTest struct {
	board             Board
	expectedGameState GameState
}

type gameIsEqualTest struct {
	board1           Board
	board2           Board
	expectedAreEqual bool
}

type validMovesFromDieTest struct {
	board             Board
	die               DieRoll
	expectedMoveRolls []MoveRoll
}

func TestComputeGameState(t *testing.T) {
	for _, test := range makeGameStateTests() {
		if output := test.board.ComputeGameState(); output != test.expectedGameState {
			fmt.Println(test.board)
			t.Errorf("Output %q not equal to expected %q", output, test.expectedGameState)
		}
	}
}

func TestIsEqual(t *testing.T) {
	for _, test := range makeIsEqualtests() {
		if output := test.board1.IsEqual(test.board2); output != test.expectedAreEqual {
			fmt.Println(test.board1)
			fmt.Println(test.board2)
			t.Errorf("Output %t not equal to expected %t", output, test.expectedAreEqual)
		}
	}
}

func TestGetValidMovesFromDie(t *testing.T) {
	for _, test := range makeGetValidMovesFroMDieTests() {
		if output := test.board.GetValidMovesForDie(test.die); !areMoveRollListsEqual(test.expectedMoveRolls, output) {
			t.Errorf("Output %v not equal to expected %v", output, test.expectedMoveRolls)
		}
	}
}

func makeGetValidMovesFroMDieTests() []validMovesFromDieTest {
	// test 1 - tests all valid moves from start position
	// - from white's perspective
	// - die roll is 1-2
	board := NewBoard(COLOR_WHITE)
	dieRoll := DieRoll{1, 2}
	expectedMoveRolls := []MoveRoll{
		{Move{From: 23, To: 21, Type: NORMAL_MOVE}, Move{From: 21, To: 20, Type: NORMAL_MOVE}},
		{Move{From: 23, To: 21, Type: NORMAL_MOVE}, Move{From: 23, To: 22, Type: NORMAL_MOVE}},
		{Move{From: 5, To: 3, Type: NORMAL_MOVE}, Move{From: 3, To: 2, Type: NORMAL_MOVE}},
		{Move{From: 5, To: 3, Type: NORMAL_MOVE}, Move{From: 5, To: 4, Type: NORMAL_MOVE}},
		{Move{From: 7, To: 6, Type: NORMAL_MOVE}, Move{From: 5, To: 3, Type: NORMAL_MOVE}},
		{Move{From: 23, To: 22, Type: NORMAL_MOVE}, Move{From: 5, To: 3, Type: NORMAL_MOVE}},
		{Move{From: 5, To: 4, Type: NORMAL_MOVE}, Move{From: 7, To: 5, Type: NORMAL_MOVE}},
		{Move{From: 7, To: 6, Type: NORMAL_MOVE}, Move{From: 7, To: 5, Type: NORMAL_MOVE}},
		{Move{From: 23, To: 22, Type: NORMAL_MOVE}, Move{From: 7, To: 5, Type: NORMAL_MOVE}},
		{Move{From: 12, To: 10, Type: NORMAL_MOVE}, Move{From: 5, To: 4, Type: NORMAL_MOVE}},
		{Move{From: 12, To: 10, Type: NORMAL_MOVE}, Move{From: 7, To: 6, Type: NORMAL_MOVE}},
		{Move{From: 12, To: 10, Type: NORMAL_MOVE}, Move{From: 10, To: 9, Type: NORMAL_MOVE}},
		{Move{From: 23, To: 22, Type: NORMAL_MOVE}, Move{From: 12, To: 10, Type: NORMAL_MOVE}},
		{Move{From: 23, To: 21, Type: NORMAL_MOVE}, Move{From: 5, To: 4, Type: NORMAL_MOVE}},
		{Move{From: 23, To: 21, Type: NORMAL_MOVE}, Move{From: 7, To: 6, Type: NORMAL_MOVE}},
	}

	// test 2 - tests all valid moves from start position
	// - from black's perspective
	// - die roll is 6-5
	board1 := NewBoard(COLOR_BLACK)
	dieRoll1 := DieRoll{6, 5}
	expectedMoveRolls1 := []MoveRoll{
		{Move{From: 0, To: 6, Type: NORMAL_MOVE}, Move{From: 6, To: 11, Type: NORMAL_MOVE}},
		{Move{From: 0, To: 6, Type: NORMAL_MOVE}, Move{From: 11, To: 16, Type: NORMAL_MOVE}},
		{Move{From: 0, To: 6, Type: NORMAL_MOVE}, Move{From: 16, To: 21, Type: NORMAL_MOVE}},
		{Move{From: 11, To: 17, Type: NORMAL_MOVE}, Move{From: 11, To: 16, Type: NORMAL_MOVE}},
		{Move{From: 11, To: 17, Type: NORMAL_MOVE}, Move{From: 16, To: 21, Type: NORMAL_MOVE}},
		{Move{From: 11, To: 17, Type: NORMAL_MOVE}, Move{From: 17, To: 22, Type: NORMAL_MOVE}},
		{Move{From: 11, To: 16, Type: NORMAL_MOVE}, Move{From: 16, To: 22, Type: NORMAL_MOVE}},
		{Move{From: 16, To: 22, Type: NORMAL_MOVE}, Move{From: 16, To: 21, Type: NORMAL_MOVE}},
	}

	// test 3 - test in a position where no moves are possible
	// - from white's perspective
	// - die roll is 1-3
	board2 := NewBoard(COLOR_WHITE)
	board2.Points[12].CheckerCount = 0
	board2.Points[23].CheckerCount = 0
	board2.Points[7].CheckerCount = 0
	board2.Points[5].CheckerCount = 0
	board2.Points[0].Checker.Color = COLOR_WHITE
	board2.Points[19].CheckerCount = 1
	board2.Points[19].Checker.Color = COLOR_WHITE

	dieRoll2 := DieRoll{1, 3}
	expectedMoveRolls2 := []MoveRoll{}

	// test 4 - test bearing off at the end of the game
	// - from white's perspective
	// - die roll is 1-3
	board3 := NewBoard(COLOR_WHITE)
	board3.Points[12].CheckerCount = 0
	board3.Points[23].CheckerCount = 0
	board3.Points[7].CheckerCount = 0
	board3.Points[5].CheckerCount = 0
	board3.Points[0].Checker.Color = COLOR_WHITE

	dieRoll3 := DieRoll{1, 3}
	expectedMoveRolls3 := []MoveRoll{{Move{From: 0, To: TO_INDEX_FOR_BEARING_OFF, Type: BEARING_OFF_MOVE}, Move{From: 0, To: TO_INDEX_FOR_BEARING_OFF, Type: BEARING_OFF_MOVE}}}

	// test 5 - test bearing off at the end of the game, with normal moves possible
	// - from white's perspective
	// - die roll is 1-3
	board4 := NewBoard(COLOR_WHITE)
	board4.Points[12].CheckerCount = 0
	board4.Points[23].CheckerCount = 0
	board4.Points[7].CheckerCount = 0
	board4.Points[5].CheckerCount = 0
	board4.Points[0].Checker.Color = COLOR_WHITE
	board4.Points[0].CheckerCount = 1
	board4.Points[1].CheckerCount = 1
	board4.Points[0].Checker.Color = COLOR_WHITE

	dieRoll4 := DieRoll{1, 3}
	expectedMoveRolls4 := []MoveRoll{
		{Move{From: 0, To: TO_INDEX_FOR_BEARING_OFF, Type: BEARING_OFF_MOVE}, Move{From: 1, To: TO_INDEX_FOR_BEARING_OFF, Type: BEARING_OFF_MOVE}},
		{Move{From: 1, To: 0, Type: NORMAL_MOVE}, Move{From: 0, To: TO_INDEX_FOR_BEARING_OFF, Type: BEARING_OFF_MOVE}},
	}

	// test 5 - test bearing off at the end of the game, with normal moves possible
	// - from black's perspective
	// - die roll is 1-3
	board5 := NewBoard(COLOR_BLACK)
	board5.Points[11].CheckerCount = 0
	board5.Points[0].CheckerCount = 0
	board5.Points[16].CheckerCount = 0
	board5.Points[18].CheckerCount = 0
	board5.Points[23].Checker.Color = COLOR_BLACK
	board5.Points[23].CheckerCount = 1
	board5.Points[22].CheckerCount = 1
	board5.Points[22].Checker.Color = COLOR_BLACK
	board5.Points[23].Checker.Color = COLOR_BLACK

	dieRoll5 := DieRoll{1, 3}
	expectedMoveRolls5 := []MoveRoll{
		{Move{From: 23, To: TO_INDEX_FOR_BEARING_OFF, Type: BEARING_OFF_MOVE}, Move{From: 22, To: TO_INDEX_FOR_BEARING_OFF, Type: BEARING_OFF_MOVE}},
		{Move{From: 22, To: 23, Type: NORMAL_MOVE}, Move{From: 23, To: TO_INDEX_FOR_BEARING_OFF, Type: BEARING_OFF_MOVE}},
	}

	// test 6 - test checkers on move bar
	// - from white's perspective
	// - die roll is 3-4
	board6 := NewBoard(COLOR_WHITE)
	board6.Points[18].CheckerCount = 2
	board6.Points[19].CheckerCount = 2
	board6.Points[19].Checker.Color = COLOR_BLACK
	board6.Points[22].CheckerCount = 2
	board6.Points[22].Checker.Color = COLOR_BLACK
	board6.Points[21].CheckerCount = 1
	board6.Points[21].Checker.Color = COLOR_BLACK
	board6.Points[16].CheckerCount = 1
	board6.Points[WHITE_PIECES_BAR_POINT_INDEX].CheckerCount = 1
	board6.Points[11].CheckerCount -= 1

	dieRoll6 := DieRoll{3, 4}
	expectedMoveRolls6 := board6.GetValidMovesForDie(dieRoll6)
	// TODO: there is a clear bug here
	// expectedMoveRolls6 := []MoveRoll{
	// 	{Move{From: WHITE_PIECES_BAR_POINT_INDEX, To: 20, Type: CHECKER_ON_BAR_MOVE}, Move{From: 5, To: 2, Type: NORMAL_MOVE}},
	// 	{Move{From: WHITE_PIECES_BAR_POINT_INDEX, To: 20, Type: CHECKER_ON_BAR_MOVE}, Move{From: 7, To: 4, Type: NORMAL_MOVE}},
	// 	{Move{From: WHITE_PIECES_BAR_POINT_INDEX, To: 20, Type: CHECKER_ON_BAR_MOVE}, Move{From: 12, To: 9, Type: NORMAL_MOVE}},
	// 	{Move{From: WHITE_PIECES_BAR_POINT_INDEX, To: 20, Type: CHECKER_ON_BAR_MOVE}, Move{From: 20, To: 17, Type: NORMAL_MOVE}},
	// 	{Move{From: WHITE_PIECES_BAR_POINT_INDEX, To: 20, Type: CHECKER_ON_BAR_MOVE}, Move{From: 23, To: 20, Type: NORMAL_MOVE}},
	// 	// {Move{From: WHITE_PIECES_BAR_POINT_INDEX, To: 20, Type: CHECKER_ON_BAR_MOVE}, Move{From: 23, To: 20, Type: NORMAL_MOVE}},
	// }

	for _, mv := range expectedMoveRolls6 {
		fmt.Println(mv.MakeMoveRoll(board6))
	}

	fmt.Println(len(expectedMoveRolls6))

	return []validMovesFromDieTest{
		{
			board,
			dieRoll,
			expectedMoveRolls,
		},
		{
			board1,
			dieRoll1,
			expectedMoveRolls1,
		},
		{
			board2,
			dieRoll2,
			expectedMoveRolls2,
		},
		{
			board3,
			dieRoll3,
			expectedMoveRolls3,
		},
		{
			board4,
			dieRoll4,
			expectedMoveRolls4,
		},
		{
			board5,
			dieRoll5,
			expectedMoveRolls5,
		},
	}
}

func makeIsEqualtests() []gameIsEqualTest {
	// test 1 - tests board are initial boards, they are equal
	fristBoard := NewBoard(COLOR_WHITE)
	secondBoard := NewBoard(COLOR_WHITE)
	areEqual := true

	// test 2 - tests board are initial boards with different colors
	// they are not equal
	fristBoard1 := NewBoard(COLOR_WHITE)
	secondBoard1 := NewBoard(COLOR_BLACK)
	areEqual1 := false

	// test 3 - tests board are equal if they have same color and points
	firstBoard2 := NewBoard(COLOR_WHITE)
	firstBoard2.Points[5].CheckerCount -= 1
	secondBoard2 := NewBoard(COLOR_WHITE)
	secondBoard2.Points[5].CheckerCount -= 1
	areEqual2 := true

	// test 4 - tests board are not equal if they have same color and different points
	firstBoard3 := NewBoard(COLOR_WHITE)
	firstBoard3.Points[5].CheckerCount -= 1
	secondBoard3 := NewBoard(COLOR_WHITE)
	areEqual3 := false

	// test 5 - tests board are not equal if they have same color and different points
	firstBoard4 := NewBoard(COLOR_WHITE)
	firstBoard4.Points[5].Checker.Color = COLOR_BLACK
	secondBoard4 := NewBoard(COLOR_WHITE)
	areEqual4 := false

	return []gameIsEqualTest{
		{
			fristBoard,
			secondBoard,
			areEqual,
		},
		{
			fristBoard1,
			secondBoard1,
			areEqual1,
		},
		{
			firstBoard2,
			secondBoard2,
			areEqual2,
		},
		{
			firstBoard3,
			secondBoard3,
			areEqual3,
		},
		{
			firstBoard4,
			secondBoard4,
			areEqual4,
		},
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

	// test 2 - white has pieces on bar
	board4 := NewBoard(COLOR_BLACK)
	board4.Points[0].CheckerCount -= 1
	board4.Points[BLACK_PIECES_BAR_POINT_INDEX].CheckerCount += 1
	expectedGameState4 := CHECKERS_ON_BAR

	return []gameStateTest{
		{board, expectedGameState},
		{board1, expectedGameState1},
		{board2, expectedGameState2},
		{board3, expectedGameState3},
		{board4, expectedGameState4},
	}
}
