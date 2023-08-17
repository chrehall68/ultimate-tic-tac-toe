package board

// board related constants
const (
	ROWS  = 3
	COLS  = 3
	CELLS = ROWS * COLS
)

// protobuf related constants
const (
	STATE_PORT   = "8000"
	ACTION_PORT  = "8001"
	RETURN_PORT  = "8002"
	MAX_MSG_SIZE = 1024
)
