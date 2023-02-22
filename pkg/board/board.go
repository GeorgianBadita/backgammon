package board

import (
	"fmt"
	"strconv"
)

type DieRoll struct {
	Die1 int
	Die2 int
}

type Checker struct {
	Color
}

func NewChecker(color Color) Checker {
	return Checker{color}
}

type Point struct {
	CheckerCount int
	PointIndex   PointIndex
	Checker      Checker
}

func (p Point) IsEmpty() bool {
	return p.CheckerCount == 0
}

func NewPoint(count int, index PointIndex, checker Checker) Point {
	return Point{count, index, checker}
}

type Board struct {
	Points      []Point
	ColorToMove Color
}

func NewBoard(color Color) Board {
	points := make([]Point, NUM_POINTS)

	checkersMap := make(map[int]struct {
		Checker
		int
	})
	checkersMap[5] = struct {
		Checker
		int
	}{NewChecker(COLOR_WHITE), 5}
	checkersMap[18] = struct {
		Checker
		int
	}{NewChecker(COLOR_BLACK), 5}
	checkersMap[7] = struct {
		Checker
		int
	}{NewChecker(COLOR_WHITE), 3}
	checkersMap[16] = struct {
		Checker
		int
	}{NewChecker(COLOR_BLACK), 3}
	checkersMap[12] = struct {
		Checker
		int
	}{NewChecker(COLOR_WHITE), 5}
	checkersMap[11] = struct {
		Checker
		int
	}{NewChecker(COLOR_BLACK), 5}
	checkersMap[23] = struct {
		Checker
		int
	}{NewChecker(COLOR_WHITE), 2}
	checkersMap[0] = struct {
		Checker
		int
	}{NewChecker(COLOR_BLACK), 2}

	for idx := 0; idx < NUM_PLAYABLE_POINTS; idx++ {
		if checkersConfig, ok := checkersMap[idx]; ok {
			points[idx] = NewPoint(checkersConfig.int, PointIndex(idx), checkersConfig.Checker)
		} else {
			points[idx] = NewPoint(0, PointIndex(idx), Checker{})
		}
	}

	points[WHITE_PIECES_BAR_POINT_INDEX] = NewPoint(0, PointIndex(WHITE_PIECES_BAR_POINT_INDEX), NewChecker(COLOR_WHITE))
	points[BLACK_PIECES_BAR_POINT_INDEX] = NewPoint(0, PointIndex(BLACK_PIECES_BAR_POINT_INDEX), NewChecker(COLOR_BLACK))

	return Board{points, color}
}

// Function to copy a board
func (b Board) CopyBoard() Board {
	initPoints := make([]Point, len(b.Points))
	copy(initPoints, b.Points)
	newBoard := NewBoard(b.ColorToMove)
	newBoard.Points = initPoints
	return newBoard
}

// Function computing the current board game state
// State for the current player can be:
//   - CHECKERS_ON_BAR - if the current plauyer has any checker on bar
//   - GAME_OVER - if either of the players finished the game (i.e. got rid of all of their checkers)
//   - BEARING_OFF - if the current player has no checker on bar and has ALL of its checkers on home board
//   - NORMAL_PLAY - otherwise, this means no checker is on bar for the current player,  the current player didn't start the bearing off phase and no player finished the game - this is the  general case
func (b Board) ComputeGameState() GameState {
	currentPlayerColor := b.ColorToMove
	// If current player has Barred pieces state is CHECKERS_ON_BAR
	if currentPlayerColor == COLOR_BLACK && b.Points[BLACK_PIECES_BAR_POINT_INDEX].CheckerCount > 0 {
		return CHECKERS_ON_BAR
	}
	if currentPlayerColor == COLOR_WHITE && b.Points[WHITE_PIECES_BAR_POINT_INDEX].CheckerCount > 0 {
		return CHECKERS_ON_BAR
	}

	// If any player has no checker on the board, then it's game over
	//TODO: consider making this a bool
	currPlayerHomePieces := numCheckersOfColor(b, currentPlayerColor)
	otherPlayerHomePieces := numCheckersOfColor(b, Color(1-currentPlayerColor))

	if currPlayerHomePieces == 0 || otherPlayerHomePieces == 0 {
		return GAME_OVER
	}

	if numCheckersInHome(b, currentPlayerColor) <= 15 {
		return BEARING_OFF
	}

	return NORMAL_PLAY
}

func (b Board) GetValidMovesForDie(d DieRoll) []MoveRoll {
	return getPossibleMoves(b, d)
}

// Function that prints a pretty string of the current board
func (b Board) String() string {
	const GREEN_COLOR = string("\033[32m")
	const BLUE_COLOR = string("\033[34m")
	const RED_COLOR = string("\033[31m")

	boardString := GREEN_COLOR + "Backgammon Board:\n"
	for idx := 12; idx < NUM_PLAYABLE_POINTS; idx++ {
		boardString += pointColor(b.Points[idx], idx)
		if idx == 17 {
			boardString += "  "
		}
	}
	boardString += " " + RED_COLOR + strconv.Itoa(b.Points[BLACK_PIECES_BAR_POINT_INDEX].CheckerCount)
	boardString += "\n\n\n\n\n"

	for idx := 11; idx >= 0; idx-- {
		boardString += pointColor(b.Points[idx], idx)
		if idx == 6 {
			boardString += "  "
		}
	}
	boardString += " " + BLUE_COLOR + strconv.Itoa(b.Points[WHITE_PIECES_BAR_POINT_INDEX].CheckerCount)
	boardString += "\n"
	return fmt.Sprint(boardString)
}

func pointColor(p Point, idx int) string {
	const BLUE_COLOR = string("\033[34m")
	const RED_COLOR = string("\033[31m")
	const YELLOW_COLOR = string("\033[33m")
	const PURPLE_COLOR = string("\033[35m")

	boardString := ""
	par := strconv.Itoa(p.CheckerCount)
	if p.IsEmpty() {
		if idx%2 == 0 {
			boardString += YELLOW_COLOR + "-"
		} else {
			boardString += PURPLE_COLOR + "-"
		}
	} else {
		if p.Checker.Color == COLOR_BLACK {
			boardString += RED_COLOR + par
		} else {
			boardString += BLUE_COLOR + par
		}
	}
	return boardString
}

func numCheckersOfColor(b Board, color Color) int {
	s := 0
	for idx := 0; idx < 24; idx++ {
		if b.Points[idx].Checker.Color == color {
			s += b.Points[idx].CheckerCount
		}
	}
	return s
}

func numCheckersInHome(b Board, color Color) int {
	s := 0
	if color == COLOR_WHITE {
		for idx := 0; idx < 6; idx++ {
			s += b.Points[idx].CheckerCount
		}
	} else {
		for idx := 18; idx < NUM_PLAYABLE_POINTS; idx++ {
			s += b.Points[idx].CheckerCount
		}
	}
	return s
}

func getPossibleMoves(b Board, d DieRoll) []MoveRoll {
	d1Moves := getMovesWithOneDie(b, d.Die1)
	moveRolls := []MoveRoll{}
	if d.Die1 != d.Die2 {
		if d.Die1 < d.Die2 {
			// Make sure we first make moves with the bigger die
			d.Die1, d.Die2 = d.Die2, d.Die1
		}

		for idx := 0; idx < len(d1Moves); idx++ {
			currMove := d1Moves[idx]
			d2Moves := getMovesWithOneDie(currMove.MakeMove(b), d.Die2)
			for jdx := 0; jdx < len(d2Moves); jdx++ {
				moveRolls = append(moveRolls, MoveRoll{currMove, d2Moves[jdx]})
			}

		}
		// If thereare no possible moves with 2 die, take only possible moves with 1 dice
		if len(moveRolls) == 0 {
			for idx := 0; idx < len(d1Moves); idx++ {
				moveRolls = append(moveRolls, MoveRoll{d1Moves[idx]})
			}
		}
	} else {
		for idx := 0; idx < len(d1Moves); idx++ {
			currd1Move := d1Moves[idx]
			move1Board := currd1Move.MakeMove(b)
			d2Moves := getMovesWithOneDie(move1Board, d.Die1)
			for jdx := 0; jdx < len(d2Moves); jdx++ {
				currd2Move := d2Moves[jdx]
				move2Board := currd2Move.MakeMove(move1Board)
				d3Moves := getMovesWithOneDie(move2Board, d.Die1)
				for tdx := 0; tdx < len(d3Moves); tdx++ {
					currd3Move := d3Moves[tdx]
					move3Board := currd3Move.MakeMove(move2Board)
					d4Moves := getMovesWithOneDie(move3Board, d.Die1)
					for zdx := 0; zdx < len(d4Moves); zdx++ {
						moveRolls = append(moveRolls, MoveRoll{currd1Move, currd2Move, currd3Move, d4Moves[zdx]})
					}
				}
				// If there are no moves with 4 die, trey with 3 die
				if len(moveRolls) == 0 {
					for tdx := 0; tdx < len(d3Moves); tdx++ {
						moveRolls = append(moveRolls, MoveRoll{currd1Move, currd2Move, d3Moves[tdx]})
					}
				}
			}
			// If there are no moves with 3 die, try with 2 die
			if len(moveRolls) == 0 {
				for jdx := 0; jdx < len(d2Moves); jdx++ {
					moveRolls = append(moveRolls, MoveRoll{currd1Move, d2Moves[jdx]})
				}
			}
		}
		// If thereare no possible moves with 2 die, take only possible moves with 1 dice
		if len(moveRolls) == 0 {
			for idx := 0; idx < len(d1Moves); idx++ {
				moveRolls = append(moveRolls, MoveRoll{d1Moves[idx]})
			}
		}
	}
	return moveRolls
}

func getMovesWithOneDie(b Board, dValue int) []Move {
	direction := 1
	if b.ColorToMove == COLOR_WHITE {
		direction = -1
	}

	moves := []Move{}

	for idx := 0; idx < NUM_PLAYABLE_POINTS; idx++ {
		if b.ColorToMove == b.Points[idx].Checker.Color && b.Points[idx].CheckerCount > 0 {
			destPointIndex := PointIndex(idx + direction*dValue)
			if isValidDestinationForChecker(b, destPointIndex) {
				moves = append(moves, Move{PointIndex(idx), destPointIndex, NORMAL_MOVE})
			}
		}
	}
	return moves
}

// Function that checks if a destination for a checker is correct
// It verifies:
// 1. If the position is within the 24 points board
// 2. If the destination contains more than 1 checkers of the opposition color
func isValidDestinationForChecker(b Board, destinationPoint PointIndex) bool {
	isCoordLegal := destinationPoint >= 0 && destinationPoint < NUM_PLAYABLE_POINTS
	if !isCoordLegal {
		return false
	}

	checkersAtDest := b.Points[destinationPoint]
	if checkersAtDest.Checker.Color != b.ColorToMove && checkersAtDest.CheckerCount > 1 {
		return false
	}
	return true
}
