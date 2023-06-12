package entities

const (
	DOWN int = iota
	LEFT
	RIGHT
)

const (
	Empty int = iota
	Player1
	Player2
)

type Position struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

type GameAttributes struct {
	Rows    int `json:"rows"`
	Columns int `json:"columns"`

	RowMovement []int   `json:"movement"`
	Board       [][]int `json:"board"`

	ActualPlayer   int      `json:"actual_player"`
	ActualPosition Position `json:"actual_position"`

	IsComingDown bool `json:"is_coming_down"`
	IsGameOver   bool `json:"is_game_over"`
}

type GameMethods struct {
	Move      func(dir int)
	MoveLeft  func()
	MoveRight func()
	MoveDown  func()

	ResetMovement func()
	TogglePlayer  func()

	IsValidRow    func(row int) bool
	IsValidColumn func(column int) bool

	CheckBoard                  func() bool
	CheckBoardHorizontal        func() bool
	CheckBoardVertical          func() bool
	CheckBoardPrimaryDiagonal   func() bool
	CheckBoardSecondaryDiagonal func() bool

	RestartGame        func()
	RestartBoard       func()
	RestartRowMovement func()
}
type Game struct {
	GameAttributes
	GameMethods
}
