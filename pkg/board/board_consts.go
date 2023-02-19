package board

type Color int8
type PointIndex int
type GameState string

const (
	COLOR_WHITE Color = 0
	COLOR_BLACK Color = 1
)

const (
	BEARING_OFF     GameState = "BEARING-OFF"
	CHECKERS_ON_BAR GameState = "CHECKERS-ON-BAR"
	NORMAL_PLAY     GameState = "NORMAL-PLAY"
	GAME_OVER       GameState = "GAME-OVER"
)

// There are 26 indexes where pieces can be
// Indexes [0..23] are where pieces can move
// Index 24 is where WHITE pieces are being barred
// Index 25 is where BLACK pieces are being barred
const NUM_POINTS = 26
const NUM_PLAYABLE_POINTS = 24
const BLACK_PIECES_BAR_POINT_INDEX = 24
const WHITE_PIECES_BAR_POINT_INDEX = 25

const INIT_NUM_CHECKERS = 15
