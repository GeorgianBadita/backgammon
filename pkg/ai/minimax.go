package ai

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/GeorgianBadita/backgammon-move-generator/pkg/board"
)

const NEG_INF float32 = float32(-2 * 1e9)
const PLUS_INF float32 = float32(2 * 1e9)

// TODO: make this more general, currently it will break if maximizing player is not white
func Minimax(b board.Board, depth int, maximizingPlayer board.Color, dieRoll board.DieRoll) board.MoveRoll {
	allPossibleMoves := b.GetValidMovesForDieRoll(dieRoll)
	if len(allPossibleMoves) == 0 {
		return []board.Move{}
	}

	bestMove := allPossibleMoves[0]
	bestScore := NEG_INF

	for idx, mvRoll := range allPossibleMoves {
		score := computeMiniMax(b, depth-1, maximizingPlayer, NEG_INF, PLUS_INF)
		fmt.Printf("Computed move: %d\n", idx)
		if score > bestScore {
			bestScore = score
			bestMove = mvRoll
		}
	}

	return bestMove
}

func computeMiniMax(b board.Board, depth int, maximizingPlayer board.Color, alpha, beta float32) float32 {
	if depth == 0 || b.ComputeGameState() == board.GAME_OVER {
		return EvaluateBoard(b)
	}
	possibleMoves := b.GetAllPossibleMoveRolls()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(possibleMoves), func(i, j int) { possibleMoves[i], possibleMoves[j] = possibleMoves[j], possibleMoves[i] })

	if maximizingPlayer == b.ColorToMove {
		value := NEG_INF

		for _, mv := range possibleMoves[:min(len(possibleMoves), 14)] {
			newBoard := mv.MakeMoveRoll(b)
			newColor := board.Color(1 - newBoard.ColorToMove)
			newBoard.ColorToMove = newColor
			value = maxFloat32(value, computeMiniMax(newBoard, depth-1, newColor, alpha, beta))
			alpha = maxFloat32(alpha, value)
			if value >= beta {
				break
			}
		}
		return value
	} else {
		value := PLUS_INF

		for _, mv := range possibleMoves[:min(len(possibleMoves), 14)] {
			newBoard := mv.MakeMoveRoll(b)
			newColor := board.Color(1 - newBoard.ColorToMove)
			newBoard.ColorToMove = newColor
			value = minFloat32(value, computeMiniMax(newBoard, depth-1, newColor, alpha, beta))
			beta = minFloat32(beta, value)
			if value <= alpha {
				break
			}
		}
		return value
	}
}
