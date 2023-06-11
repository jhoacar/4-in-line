package entities

const (
	DOWN byte = iota
	LEFT
	RIGHT
)

const (
	Empty byte = iota
	Player1
	Player2
)

type Position struct {
	Row    int
	Column int
}

type Game struct {
	Rows    byte
	Columns byte

	RowMovement []byte
	Board       [][]byte

	ActualPlayer   byte
	ActualPosition Position

	Move      func(dir byte)
	MoveLeft  func()
	MoveRight func()
	MoveDown  func()

	IsComingDown bool
	IsGameOver   bool

	ResetMovement func()
	TogglePlayer  func()

	IsValidRow    func(row int) bool
	IsValidColumn func(column int) bool

	CheckBoard                  func() bool
	CheckBoardHorizontal        func() bool
	CheckBoardVertical          func() bool
	CheckBoardPrimaryDiagonal   func() bool
	CheckBoardSecondaryDiagonal func() bool
}
