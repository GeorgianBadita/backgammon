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

func makeNormalMoveTests() []moveTest {
	normalMovesTests := []moveTest{}

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

	normalMovesTests = append(normalMovesTests, moveTest{
		startBoard,
		moveRoll,
		endBoard,
	}, moveTest{
		startBoard1,
		moveRoll1,
		endBoard1,
	}, moveTest{
		startBoard2,
		moveRoll2,
		endBoard2,
	}, moveTest{
		startBoard3,
		moveRoll3,
		endBoard3,
	}, moveTest{
		startBoard4,
		moveRoll4,
		endBoard4,
	})

	return normalMovesTests
}

func makeCheckerOnBarMoveTests() []moveTest {
	normalMovesTests := []moveTest{}

	//TODO: come back here to add more tests
	startBoard := NewBoard(COLOR_WHITE)
	startBoard.Points[17].CheckerCount = 1
	startBoard.Points[16].CheckerCount -= 1
	startBoard.Points[17].Checker = startBoard.Points[16].Checker
	moveRoll := MoveRoll{
		Move{From: 23, To: 17, Type: NORMAL_MOVE},
	}
	endBoard := startBoard.CopyBoard()
	endBoard.Points[BLACK_PIECES_BAR_POINT_INDEX].CheckerCount += 1
	endBoard.Points[17].CheckerCount = 1
	endBoard.Points[17].Checker = endBoard.Points[23].Checker
	endBoard.Points[23].CheckerCount -= 1

	normalMovesTests = append(normalMovesTests, moveTest{
		startBoard,
		moveRoll,
		endBoard,
	})

	return normalMovesTests
}
