package ai

import (
	"math"

	"github.com/GeorgianBadita/backgammon-move-generator/pkg/board"
)

func EvaluateBoard(b board.Board) float32 {
	if b.ColorToMove == board.COLOR_WHITE {
		return evaluateBoard(b)
	}
	return -evaluateBoard(b)
}

func evaluateBoard(b board.Board) float32 {
	blackBoard := b.CopyBoard()
	blackBoard.ColorToMove = board.COLOR_BLACK
	whiteBoard := b.CopyBoard()
	whiteBoard.ColorToMove = board.COLOR_WHITE

	distanceScoreDiff := checkerDistanceScore(blackBoard) - checkerDistanceScore(whiteBoard)
	coverageScoreDiff := goodCoverageCheckerScore(blackBoard) - goodCoverageCheckerScore(whiteBoard)
	unsafeScoreDiff := unsafeCheckersScore(blackBoard) - unsafeCheckersScore(whiteBoard)

	return 100*distanceScoreDiff + float32(50*coverageScoreDiff) + float32(50*unsafeScoreDiff)
}

func checkerDistanceScore(b board.Board) float32 {
	var sum float32 = 0.0
	for idx := 0; idx < board.NUM_PLAYABLE_POINTS; idx++ {
		if b.ColorToMove == board.COLOR_BLACK {
			if b.Points[idx].Checker.Color == b.ColorToMove {
				sum += float32(b.Points[idx].CheckerCount * (board.NUM_PLAYABLE_POINTS - idx))
			}
		} else {
			if b.Points[idx].Checker.Color == b.ColorToMove {
				sum += float32(b.Points[idx].CheckerCount * (idx + 1))
			}
		}
	}

	if (b.ColorToMove == board.COLOR_WHITE && b.Points[board.WHITE_PIECES_BAR_POINT_INDEX].CheckerCount > 0) ||
		(b.ColorToMove == board.COLOR_BLACK && b.Points[board.BLACK_PIECES_BAR_POINT_INDEX].CheckerCount > 0) {
		spots := freeSpots(b, board.Color(1-b.ColorToMove))
		expectedIdle := expectedDelay(float32(spots))

		cntBarred := b.Points[board.WHITE_PIECES_BAR_POINT_INDEX].CheckerCount
		if b.ColorToMove == board.COLOR_BLACK {
			cntBarred = b.Points[board.BLACK_PIECES_BAR_POINT_INDEX].CheckerCount
		}
		sum += ((expectedIdle-1)*7.0 + 24) * float32(cntBarred)
	}
	return sum
}

func unsafeCheckersScore(b board.Board) int {
	unsafeCheckers := 0
	for idx := 0; idx < board.NUM_PLAYABLE_POINTS; idx++ {
		if b.Points[idx].CheckerCount == 0 || b.Points[idx].Checker.Color != b.ColorToMove {
			continue
		}

		if b.Points[idx].CheckerCount == 1 {
			if b.ColorToMove == board.COLOR_WHITE && idx >= 18 {
				unsafeCheckers++
			} else if b.ColorToMove == board.COLOR_BLACK && idx <= 5 {
				unsafeCheckers++
			} else if canBeBarredNextMove(b, idx) {
				unsafeCheckers += 5
			}
		}
	}

	return unsafeCheckers
}

func goodCoverageCheckerScore(b board.Board) float32 {
	var score float32
	score = 0.0

	if b.ColorToMove == board.COLOR_WHITE {
		for idx := board.NUM_PLAYABLE_POINTS - 1; idx >= 6; idx-- {
			if isSafe(b, board.PointIndex(idx)) {
				score++
			}
		}
		if isSafe(b, 6) {
			score += 3.5
		}
		if isSafe(b, 5) {
			score += 5
		}
		if isSafe(b, 4) {
			score += 4
		}
		if isSafe(b, 3) {
			score += 3
		}
		if isSafe(b, 2) {
			score += 2
		}
		if isSafe(b, 1) {
			score += 1.5
		}
		if isSafe(b, 0) {
			score += 1
		}
	} else {
		for idx := 0; idx < 18; idx++ {
			if isSafe(b, board.PointIndex(idx)) {
				score++
			}
		}
		if isSafe(b, 17) {
			score += 3.5
		}
		if isSafe(b, 18) {
			score += 5
		}
		if isSafe(b, 19) {
			score += 4
		}
		if isSafe(b, 20) {
			score += 3
		}
		if isSafe(b, 21) {
			score += 2
		}
		if isSafe(b, 22) {
			score += 1.5
		}
		if isSafe(b, 23) {
			score += 1
		}
	}

	return score
}

func canBeBarredNextMove(b board.Board, idx int) bool {
	if b.ColorToMove == board.COLOR_WHITE {
		for jdx := idx - 1; jdx >= max(idx-6, 0); jdx-- {
			if b.Points[jdx].CheckerCount > 0 && b.Points[jdx].Checker.Color == board.COLOR_BLACK {
				return true
			}
		}
	} else {
		for jdx := idx + 1; jdx < min(idx+6, board.NUM_PLAYABLE_POINTS); jdx++ {
			if b.Points[jdx].CheckerCount > 0 && b.Points[jdx].Checker.Color == board.COLOR_WHITE {
				return true
			}
		}
	}
	return false
}

func freeSpots(b board.Board, forColor board.Color) int {
	cnt := 0
	if forColor == board.COLOR_BLACK {
		for idx := 18; idx < board.NUM_PLAYABLE_POINTS; idx++ {
			if b.Points[idx].Checker.Color == board.COLOR_BLACK && b.Points[idx].CheckerCount > 1 {
				continue
			}
			cnt++
		}
	} else {
		for idx := 0; idx < 6; idx++ {
			if b.Points[idx].Checker.Color == board.COLOR_WHITE && b.Points[idx].CheckerCount > 1 {
				continue
			}
			cnt++
		}
	}
	return cnt
}

func expectedDelay(spots float32) float32 {
	if spots == 0 {
		return 10
	}

	missProbabilty := float32(math.Pow(float64((6.0-spots)/6.0), 2))
	probabiltyToEnter := 1 - missProbabilty
	return 1.0 / probabiltyToEnter
}

func isSafe(b board.Board, idx board.PointIndex) bool {
	return b.Points[idx].Checker.Color == b.ColorToMove && b.Points[idx].CheckerCount > 1
}

func min(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a >= b {
		return a
	}
	return b
}

func maxFloat32(a float32, b float32) float32 {
	if a >= b {
		return a
	}
	return b
}

func minFloat32(a float32, b float32) float32 {
	if a <= b {
		return a
	}
	return b
}
