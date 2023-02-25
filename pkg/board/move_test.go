package board

import (
	"fmt"
	"testing"
)

type moveTest struct {
	initialBoard  Board
	moveRoll      MoveRoll
	expectedBoard Board
}

func TestMove_NORMAL_MOVE(t *testing.T) {
	for _, test := range makeNormalMoveTests() {
		if output := test.moveRoll.MakeMoveRoll(test.initialBoard); !output.IsEqual(test.expectedBoard) {
			fmt.Println("Expected")
			fmt.Println(test.expectedBoard)
			fmt.Println()
			fmt.Println("Actual")
			fmt.Println(output)
			t.Errorf("Output %q not equal to expected %q", output, test.expectedBoard)
		}
	}
}

func TestMove_CHECKER_ON_BAR_MOVE(t *testing.T) {
	for _, test := range makeCheckerOnBarMoveTests() {
		if output := test.moveRoll.MakeMoveRoll(test.initialBoard); !output.IsEqual(test.expectedBoard) {
			fmt.Println("Expected")
			fmt.Println(test.expectedBoard)
			fmt.Println()
			fmt.Println("Actual")
			fmt.Println(output)
			t.Errorf("Output %q not equal to expected %q", output, test.expectedBoard)
		}
	}
}

func TestMove_BEARING_OFF_MOVE(t *testing.T) {
	// ARRANGE
	board := NewBoard(COLOR_BLACK)
	board.Points[16].CheckerCount = 0
	board.Points[11].CheckerCount = 0
	board.Points[0].CheckerCount = 0
	board.Points[18].CheckerCount = 2
	board.Points[22].CheckerCount = 3
	board.Points[22].Checker.Color = COLOR_BLACK
	moveRoll := MoveRoll{Move{
		From: 22, To: TO_INDEX_FOR_BEARING_OFF, Type: BEARING_OFF_MOVE,
	}}
	expectedBoard := board.CopyBoard()
	expectedBoard.Points[22].CheckerCount -= 1

	// ACT
	newBoard := moveRoll.MakeMoveRoll(board)

	// ASSERT
	if !newBoard.IsEqual(expectedBoard) {
		t.Errorf("Output %q not equal to expected %q", newBoard, expectedBoard)
	}
}

func makeNormalMoveTests() []moveTest {
	// test1 - apply a 6 move from white's perspective at the
	// beginning of the game
	startBoard := NewBoard(COLOR_WHITE)
	moveRoll := MoveRoll{Move{From: 23, To: 17, Type: NORMAL_MOVE}}
	endBoard := NewBoard(COLOR_WHITE)
	endBoard.Points[NUM_PLAYABLE_POINTS-1].CheckerCount -= 1
	endBoard.Points[NUM_PLAYABLE_POINTS-1-6].CheckerCount += 1
	endBoard.Points[NUM_PLAYABLE_POINTS-1-6].Checker = endBoard.Points[NUM_PLAYABLE_POINTS-1].Checker

	// test2 - apply a 6 move from black's perspective at the
	// beginning of the game
	startBoard1 := NewBoard(COLOR_BLACK)
	moveRoll1 := MoveRoll{Move{From: 11, To: 17, Type: NORMAL_MOVE}}
	endBoard1 := NewBoard(COLOR_BLACK)
	endBoard1.Points[11].CheckerCount -= 1
	endBoard1.Points[17].CheckerCount += 1
	endBoard1.Points[17].Checker = endBoard.Points[11].Checker

	// test 3 - apply two moves a 3 and a 1 at the beginning of the game for white
	startBoard2 := NewBoard(COLOR_WHITE)
	moveRoll2 := MoveRoll{
		Move{From: 7, To: 4, Type: NORMAL_MOVE},
		Move{From: 5, To: 4, Type: NORMAL_MOVE},
	}
	endBoard2 := NewBoard(COLOR_WHITE)
	endBoard2.Points[7].CheckerCount -= 1
	endBoard2.Points[5].CheckerCount -= 1
	endBoard2.Points[4].CheckerCount = 2
	endBoard2.Points[4].Checker = endBoard.Points[7].Checker

	// test 4 - apply 4 moves of 6es at the beginning of the game for black
	startBoard3 := NewBoard(COLOR_BLACK)
	moveRoll3 := MoveRoll{
		Move{From: 11, To: 17, Type: NORMAL_MOVE},
		Move{From: 11, To: 17, Type: NORMAL_MOVE},
		Move{From: 11, To: 17, Type: NORMAL_MOVE},
		Move{From: 11, To: 17, Type: NORMAL_MOVE}}
	endBoard3 := NewBoard(COLOR_BLACK)
	endBoard3.Points[11].CheckerCount = 1
	endBoard3.Points[17].CheckerCount = 4
	endBoard3.Points[17].Checker = endBoard3.Points[11].Checker

	// test 5 - apply 1 move that bars one of black's pieces
	startBoard4 := NewBoard(COLOR_WHITE)
	startBoard4.Points[17].CheckerCount = 1
	startBoard4.Points[16].CheckerCount -= 1
	startBoard4.Points[17].Checker = startBoard4.Points[16].Checker
	moveRoll4 := MoveRoll{
		Move{From: 23, To: 17, Type: NORMAL_MOVE},
	}
	endBoard4 := startBoard4.CopyBoard()
	endBoard4.Points[BLACK_PIECES_BAR_POINT_INDEX].CheckerCount += 1
	endBoard4.Points[17].CheckerCount = 1
	endBoard4.Points[17].Checker = endBoard4.Points[23].Checker
	endBoard4.Points[23].CheckerCount -= 1

	// test 6 - apply 1 mote that bars one of white's pieces
	startBoard5 := NewBoard(COLOR_BLACK)
	startBoard5.Points[7].CheckerCount = 1
	startBoard5.Points[8].CheckerCount -= 1
	startBoard5.Points[7].Checker = startBoard4.Points[8].Checker
	moveRoll5 := MoveRoll{
		Move{From: 0, To: 7, Type: NORMAL_MOVE},
	}
	endBoard5 := startBoard5.CopyBoard()
	endBoard5.Points[WHITE_PIECES_BAR_POINT_INDEX].CheckerCount += 1
	endBoard5.Points[7].CheckerCount = 1
	endBoard5.Points[7].Checker = endBoard4.Points[0].Checker
	endBoard5.Points[0].CheckerCount -= 1

	return []moveTest{
		{
			startBoard,
			moveRoll,
			endBoard,
		}, {
			startBoard1,
			moveRoll1,
			endBoard1,
		}, {
			startBoard2,
			moveRoll2,
			endBoard2,
		}, {
			startBoard3,
			moveRoll3,
			endBoard3,
		}, {
			startBoard4,
			moveRoll4,
			endBoard4,
		},
		{
			startBoard5,
			moveRoll5,
			endBoard5,
		},
	}
}

func makeCheckerOnBarMoveTests() []moveTest {
	// test 1 - Black has a piece on bar
	startBoard := NewBoard(COLOR_BLACK)
	startBoard.Points[17].CheckerCount = 1
	startBoard.Points[16].CheckerCount -= 1
	startBoard.Points[17].Checker = startBoard.Points[16].Checker
	startBoard.Points[BLACK_PIECES_BAR_POINT_INDEX].CheckerCount += 1
	moveRoll := MoveRoll{
		Move{From: BLACK_PIECES_BAR_POINT_INDEX, To: 0, Type: CHECKER_ON_BAR_MOVE},
	}
	endBoard := startBoard.CopyBoard()
	endBoard.Points[0].CheckerCount += 1
	endBoard.Points[BLACK_PIECES_BAR_POINT_INDEX].CheckerCount = 0

	// test 2 - White has two pieces on bar, when enters black home
	// and lands on JUST one checker, it barrs it
	startBoard1 := NewBoard(COLOR_WHITE)
	startBoard1.Points[12].CheckerCount -= 1
	startBoard1.Points[23].CheckerCount -= 1
	startBoard1.Points[WHITE_PIECES_BAR_POINT_INDEX].CheckerCount += 2
	startBoard1.Points[18].CheckerCount = 1
	moveRoll1 := MoveRoll{
		Move{From: WHITE_PIECES_BAR_POINT_INDEX, To: 18, Type: CHECKER_ON_BAR_MOVE},
	}
	endBoard1 := startBoard1.CopyBoard()
	endBoard1.Points[WHITE_PIECES_BAR_POINT_INDEX].CheckerCount = 1
	endBoard1.Points[BLACK_PIECES_BAR_POINT_INDEX].CheckerCount = 1
	endBoard1.Points[18].Checker.Color = COLOR_WHITE

	return []moveTest{
		{
			startBoard,
			moveRoll,
			endBoard,
		}, {
			startBoard1,
			moveRoll1,
			endBoard1,
		},
	}
}
