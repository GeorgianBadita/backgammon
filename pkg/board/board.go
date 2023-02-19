package board

import "strconv"

type Color int8
type PointIndex int
type GameState string

const (
	COLOR_WHITE Color = 1
	COLOR_BLACK Color = 2
	// NO_COLOR used to represent an empty point and/or the beginning of the game
	// where there is none turn to move
	NO_COLOR Color = 3
)

const (
	BEARING_OFF     GameState = "BEARING-OFF"
	CHECKERS_ON_BAR GameState = "CHECKERS-ON-BAR"
	NORMAL_PLAY     GameState = "NORMAL-PLAY"
)

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
	GameState   GameState
	ColorToMove Color
}

func NewBoard() Board {
	points := make([]Point, 24)

	checkersMap := make(map[int]int)
	checkersMap[5] = 5
	checkersMap[7] = 3
	checkersMap[12] = 5
	checkersMap[23] = 2

	for idx := 0; idx < 24; idx++ {
		if count, ok := checkersMap[idx]; ok {
			points[idx] = NewPoint(count, PointIndex(idx), NewChecker(COLOR_WHITE))
			points[23-idx] = NewPoint(count, PointIndex(idx), NewChecker(COLOR_BLACK))
		} else {
			points[idx] = NewPoint(0, PointIndex(idx), NewChecker(NO_COLOR))
		}
	}
	return Board{points, NORMAL_PLAY, NO_COLOR}
}

func (b Board) String() string {
	boardString := ""
	for idx := 12; idx < 24; idx++ {
		par := strconv.Itoa(b.Points[idx].CheckerCount)
		boardString += par + " "
	}
	boardString += "\n"

	for idx := 11; idx >= 0; idx-- {
		par := strconv.Itoa(b.Points[idx].CheckerCount)
		boardString += par + " "
	}
	boardString += "\n"
	return boardString
}
