package board

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/hashstructure/v2"
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

// Function to hash a board
func (b Board) Hash() uint64 {
	hash, err := hashstructure.Hash(b, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}
	return hash
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
	currPlayerCheckerCount := numCheckersOfColor(b, currentPlayerColor)
	otherPlayerCheckerCount := numCheckersOfColor(b, Color(1-currentPlayerColor))

	if currPlayerCheckerCount == 0 || otherPlayerCheckerCount == 0 {
		return GAME_OVER
	}

	// If all checkers are in home
	if currPlayerCheckerCount == numCheckersInHome(b, currentPlayerColor) {
		return BEARING_OFF
	}

	return NORMAL_PLAY
}

func (b Board) GetAllPossibleMoveRolls() []MoveRoll {
	mvRolls := []MoveRoll{}
	for idx := 1; idx <= 6; idx++ {
		for jdx := 1; jdx <= 6; jdx++ {
			moves := b.GetValidMovesForDieRoll(DieRoll{idx, jdx})
			mvRolls = append(mvRolls, moves...)
		}
	}
	return mvRolls
}

func (b Board) GetValidMovesForDieRoll(d DieRoll) []MoveRoll {
	return getPossibleMoves(b, d)
}

func (b Board) GetValidMovesForDie(d int) []Move {
	return getMovesWithOneDie(b, d)
}

/**
 * Function to serialize a board to a string representation
 * serialized string has a form of:
 * 1-3/2-4/3-1:5-1/6-2/11-3 1 0 w
 * meaning: up until : are white's pieces, after are black pieces
 * each group of x-y means there are y checkers at xth point
 * after that the first number is the number of barred checkers for white
 * the other number is the number of barred checkers for black
 * the last number is the current player turn
 */
func (b Board) SerializeBoard() string {
	whiteString := ""
	blackString := ""
	for idx := 0; idx < NUM_PLAYABLE_POINTS; idx++ {
		if b.Points[idx].CheckerCount > 0 {
			fmtString := fmt.Sprintf("%d-%d/", idx+1, b.Points[idx].CheckerCount)
			if b.Points[idx].Checker.Color == COLOR_WHITE {
				whiteString += fmtString
			} else {
				blackString += fmtString
			}
		}
	}
	if len(whiteString) > 0 {
		whiteString = whiteString[:len(whiteString)-1]
	}

	if len(blackString) > 0 {
		blackString = blackString[:len(blackString)-1]
	}
	colorToMove := "w"
	if b.ColorToMove == COLOR_BLACK {
		colorToMove = "b"
	}
	return fmt.Sprintf("%s:%s %d %d %s", whiteString, blackString, b.Points[WHITE_PIECES_BAR_POINT_INDEX].CheckerCount, b.Points[BLACK_PIECES_BAR_POINT_INDEX].CheckerCount, colorToMove)
}

/**
* Function to deserialize a string board to a Board struct
* @param boardStr - serialized string of a backgammon board
* serialized string has a form of:
* 1-3/2-4/3-1:5-1/6-2/11-3 1 0 w
* meaning: up until : are white's pieces, after are black pieces
* each group of x-y means there are y checkers at xth point
* after that the first number is the number of barred checkers for white
* the other number is the number of barred checkers for black
* the last number is the current player turn
 */
func DeserializeBoard(boardStr string) Board {
	splitString := strings.Split(boardStr, ":")
	whiteStr := splitString[0]
	restSplit := strings.Split(splitString[1], " ")
	blackStr := restSplit[0]
	barredWhite := restSplit[1]
	barrdeBlack := restSplit[2]
	turn := restSplit[3]

	board := NewBoard(COLOR_BLACK)
	if turn == "w" {
		board.ColorToMove = COLOR_WHITE
	}

	for idx := 0; idx < NUM_PLAYABLE_POINTS; idx++ {
		board.Points[idx].CheckerCount = 0
	}

	for _, split := range strings.Split(whiteStr, "/") {
		pointSplit := strings.Split(split, "-")
		pointIdx, _ := strconv.Atoi(pointSplit[0])
		numCheckers, _ := strconv.Atoi(pointSplit[1])
		board.Points[pointIdx-1].CheckerCount = numCheckers
		board.Points[pointIdx-1].Checker.Color = COLOR_WHITE
	}

	for _, split := range strings.Split(blackStr, "/") {
		pointSplit := strings.Split(split, "-")
		pointIdx, _ := strconv.Atoi(pointSplit[0])
		numCheckers, _ := strconv.Atoi(pointSplit[1])
		board.Points[pointIdx-1].CheckerCount = numCheckers
		board.Points[pointIdx-1].Checker.Color = COLOR_BLACK
	}
	numCheckers, _ := strconv.Atoi(barredWhite)
	board.Points[WHITE_PIECES_BAR_POINT_INDEX].CheckerCount = numCheckers

	numCheckers, _ = strconv.Atoi(barrdeBlack)
	board.Points[BLACK_PIECES_BAR_POINT_INDEX].CheckerCount = numCheckers
	return board
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
	if p.CheckerCount == 0 {
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

func (b Board) IsEqual(ot Board) bool {
	return b.Hash() == ot.Hash()
}

func numCheckersOfColor(b Board, color Color) int {
	s := 0
	for idx := 0; idx < NUM_PLAYABLE_POINTS; idx++ {
		if b.Points[idx].Checker.Color == color {
			s += b.Points[idx].CheckerCount
		}
	}
	if color == COLOR_WHITE {
		s += b.Points[WHITE_PIECES_BAR_POINT_INDEX].CheckerCount
	} else {
		s += b.Points[BLACK_PIECES_BAR_POINT_INDEX].CheckerCount
	}
	return s
}

func numCheckersInHome(b Board, color Color) int {
	s := 0
	if color == COLOR_WHITE {
		for idx := 0; idx < 6; idx++ {
			if b.Points[idx].Checker.Color == COLOR_WHITE {
				s += b.Points[idx].CheckerCount
			}
		}
	} else {
		for idx := 18; idx < NUM_PLAYABLE_POINTS; idx++ {
			if b.Points[idx].Checker.Color == COLOR_BLACK {
				s += b.Points[idx].CheckerCount
			}
		}
	}
	return s
}

func getPossibleMoves(b Board, d DieRoll) []MoveRoll {
	moveRolls := []MoveRoll{}
	seenBoards := map[uint64]bool{}

	if d.Die1 != d.Die2 {
		// Try the bigger die first
		if d.Die1 < d.Die2 {
			d.Die1, d.Die2 = d.Die2, d.Die1
		}

		d1Moves := getMovesWithOneDie(b, d.Die1)
		for idx := 0; idx < len(d1Moves); idx++ {
			currMove := d1Moves[idx]
			mv1Board := currMove.MakeMove(b)
			d2Moves := getMovesWithOneDie(mv1Board, d.Die2)
			for jdx := 0; jdx < len(d2Moves); jdx++ {
				mvRollToAdd := MoveRoll{currMove, d2Moves[jdx]}
				mv2Board := d2Moves[jdx].MakeMove(mv1Board)
				mv2BoardHash := mv2Board.Hash()
				if _, ok := seenBoards[mv2BoardHash]; !ok {
					moveRolls = append(moveRolls, mvRollToAdd)
					seenBoards[mv2BoardHash] = true
				}
			}
		}

		d1MovesRev := getMovesWithOneDie(b, d.Die2)
		for idx := 0; idx < len(d1MovesRev); idx++ {
			currMove := d1MovesRev[idx]
			mv1Board := currMove.MakeMove(b)
			d2MovesRev := getMovesWithOneDie(mv1Board, d.Die1)
			for jdx := 0; jdx < len(d2MovesRev); jdx++ {
				mvRollToAdd := MoveRoll{currMove, d2MovesRev[jdx]}
				mv2Board := d2MovesRev[jdx].MakeMove(mv1Board)
				mv2BoardHash := mv2Board.Hash()
				if _, ok := seenBoards[mv2BoardHash]; !ok {
					moveRolls = append(moveRolls, mvRollToAdd)
					seenBoards[mv2BoardHash] = true
				}
			}
		}

		// If thereare no possible moves with 2 die, take only possible moves with 1 dice
		if len(moveRolls) == 0 {
			for idx := 0; idx < len(d1Moves); idx++ {
				currMove := d1Moves[idx]
				boardHash := currMove.MakeMove(b).Hash()
				if _, ok := seenBoards[boardHash]; !ok {
					moveRolls = append(moveRolls, MoveRoll{currMove})
					seenBoards[boardHash] = true
				}
			}
		}

		if len(moveRolls) == 0 {
			for idx := 0; idx < len(d1MovesRev); idx++ {
				currMove := d1MovesRev[idx]
				boardHash := currMove.MakeMove(b).Hash()
				if _, ok := seenBoards[boardHash]; !ok {
					moveRolls = append(moveRolls, MoveRoll{currMove})
					seenBoards[boardHash] = true
				}
			}
		}
	} else {
		d1Moves := getMovesWithOneDie(b, d.Die1)
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
						moveRollToAdd := MoveRoll{currd1Move, currd2Move, currd3Move, d4Moves[zdx]}
						finalMoveBaordHash := d4Moves[zdx].MakeMove(move3Board).Hash()
						if _, ok := seenBoards[finalMoveBaordHash]; !ok {
							moveRolls = append(moveRolls, moveRollToAdd)
							seenBoards[finalMoveBaordHash] = true
						}
					}
				}
				// If there are no moves with 4 die, trey with 3 die
				if len(moveRolls) == 0 {
					for tdx := 0; tdx < len(d3Moves); tdx++ {
						moveRollToAdd := MoveRoll{currd1Move, currd2Move, d3Moves[tdx]}
						finalMoveBaordHash := d3Moves[tdx].MakeMove(move2Board).Hash()
						if _, ok := seenBoards[finalMoveBaordHash]; !ok {
							moveRolls = append(moveRolls, moveRollToAdd)
							seenBoards[finalMoveBaordHash] = true
						}
					}
				}
			}
			// If there are no moves with 3 die, try with 2 die
			if len(moveRolls) == 0 {
				for jdx := 0; jdx < len(d2Moves); jdx++ {
					moveRollToAdd := MoveRoll{currd1Move, d2Moves[jdx]}
					finalMoveBaordHash := d2Moves[jdx].MakeMove(move1Board).Hash()
					if _, ok := seenBoards[finalMoveBaordHash]; !ok {
						moveRolls = append(moveRolls, moveRollToAdd)
						seenBoards[finalMoveBaordHash] = true
					}
				}
			}
		}
		// If there are no possible moves with 2 die, take only possible moves with 1 dice
		if len(moveRolls) == 0 {
			for idx := 0; idx < len(d1Moves); idx++ {
				moveRollToAdd := MoveRoll{d1Moves[idx]}
				finalMoveBaordHash := d1Moves[idx].MakeMove(b).Hash()
				if _, ok := seenBoards[finalMoveBaordHash]; !ok {
					moveRolls = append(moveRolls, moveRollToAdd)
					seenBoards[finalMoveBaordHash] = true
				}
			}
		}
	}
	return moveRolls
}

func getMovesWithOneDie(b Board, dValue int) []Move {
	switch b.ComputeGameState() {
	case NORMAL_PLAY:
		return getMovesForNormalGameState(b, dValue)
	case CHECKERS_ON_BAR:
		return getMovesForCheckersOnBarState(b, dValue)
	case BEARING_OFF:
		return getMovesForBearingOffState(b, dValue)
	default:
		return []Move{}
	}
}

// Get all the moves for one die on a normal game state
// Assumes the function is only called if the board state is NORMAL_PLAY
func getMovesForNormalGameState(b Board, dValue int) []Move {
	moves := []Move{}
	direction := 1
	if b.ColorToMove == COLOR_WHITE {
		direction = -1
	}

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

// Get all the moves for one die on a checker on bar state
// Assumes the function is only called if the board state CHECKERS_ON_BAR
func getMovesForCheckersOnBarState(b Board, dValue int) []Move {
	if b.ColorToMove == COLOR_WHITE {
		pointsAtIdx := b.Points[NUM_PLAYABLE_POINTS-dValue]
		if pointsAtIdx.Checker.Color == b.ColorToMove ||
			pointsAtIdx.CheckerCount < 2 {
			return []Move{{WHITE_PIECES_BAR_POINT_INDEX, PointIndex(NUM_PLAYABLE_POINTS - dValue), CHECKER_ON_BAR_MOVE}}
		}

	} else {
		pointsAtIdx := b.Points[dValue-1]
		if pointsAtIdx.Checker.Color == b.ColorToMove ||
			pointsAtIdx.CheckerCount < 2 {
			return []Move{{BLACK_PIECES_BAR_POINT_INDEX, PointIndex(dValue - 1), CHECKER_ON_BAR_MOVE}}
		}
	}
	return []Move{}
}

// Get all the moves for one die on a bearing off state
// Assumes the function is only called if the board state BEARING_OFF
func getMovesForBearingOffState(b Board, dValue int) []Move {
	movesMap := map[Move]bool{}
	moves := []Move{}
	// First go for normal moves that can be done during bear off
	if b.ColorToMove == COLOR_WHITE {
		for idx := 0; idx < 6; idx++ {
			if b.ColorToMove == b.Points[idx].Checker.Color && b.Points[idx].CheckerCount > 0 {
				destPointIndex := PointIndex(idx - dValue)
				if isValidDestinationForChecker(b, destPointIndex) {
					movesMap[Move{PointIndex(idx), destPointIndex, NORMAL_MOVE}] = true
				}
			}
		}
	} else {
		for idx := 18; idx < NUM_PLAYABLE_POINTS; idx++ {
			if b.ColorToMove == b.Points[idx].Checker.Color && b.Points[idx].CheckerCount > 0 {
				destPointIndex := PointIndex(idx + dValue)
				if isValidDestinationForChecker(b, destPointIndex) {
					movesMap[Move{PointIndex(idx), destPointIndex, NORMAL_MOVE}] = true
				}
			}
		}
	}

	// Go for removal moves in bearing off when die value is greater than first index with
	// checker of the same color as the
	if b.ColorToMove == COLOR_WHITE {
		indexOfLastChecker := 5
		for idx := 5; idx >= 0; idx-- {
			if b.Points[idx].CheckerCount > 0 && b.Points[idx].Checker.Color == b.ColorToMove {
				indexOfLastChecker = idx
				break
			}
		}
		if dValue >= indexOfLastChecker+1 {
			movesMap[Move{PointIndex(indexOfLastChecker), TO_INDEX_FOR_BEARING_OFF, BEARING_OFF_MOVE}] = true
		}
		// Base case when player has checkers on the die position
		// This can lead to duplicate moves, thus using a map
		if b.Points[dValue-1].CheckerCount > 0 && b.Points[dValue-1].Checker.Color == b.ColorToMove {
			movesMap[Move{PointIndex(dValue - 1), TO_INDEX_FOR_BEARING_OFF, BEARING_OFF_MOVE}] = true
		}
	} else {
		indexOfLastChecker := 18
		for idx := 18; idx < NUM_PLAYABLE_POINTS; idx++ {
			if b.Points[idx].CheckerCount > 0 && b.Points[idx].Checker.Color == b.ColorToMove {
				indexOfLastChecker = idx
				break
			}
		}
		if dValue >= NUM_PLAYABLE_POINTS-indexOfLastChecker {
			movesMap[Move{PointIndex(indexOfLastChecker), TO_INDEX_FOR_BEARING_OFF, BEARING_OFF_MOVE}] = true
		}

		// Base case when player has checkers on the die position
		// This can lead to duplicate moves, thus using a map
		if b.Points[NUM_PLAYABLE_POINTS-dValue].CheckerCount > 0 && b.Points[NUM_PLAYABLE_POINTS-dValue].Checker.Color == b.ColorToMove {
			movesMap[Move{PointIndex(NUM_PLAYABLE_POINTS - dValue), TO_INDEX_FOR_BEARING_OFF, BEARING_OFF_MOVE}] = true
		}
	}

	for move := range movesMap {
		moves = append(moves, move)
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
