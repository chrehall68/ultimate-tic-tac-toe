package game

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"uttt/pkg/board"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Runner struct {
	turn      bool
	gameboard *board.Board
}

func NewRunner() *Runner {
	return &Runner{turn: true, gameboard: board.NewProtoBoard()}
}

// =======================================================
// =========== Player Types ===========
// =======================================================

type Player interface {
	// getMove returns:
	//     - *board.Move - the move to make
	//     - bool - whether or not the user requested to quit;
	//            true if yes, false if no
	getMove() (*board.Move, bool)

	// displayBoard parameters:
	//     - *board.Board - the current board
	//     - *board.Owner - whose turn it is currently (PLAYER1 or PLAYER2)
	displayBoard(*board.Board, *board.Owner)

	// afterMove parameters:
	//     - *board.Board - the (changed) board
	//     - bool - whether or not the previous move was valid
	afterMove(*board.Board, bool)
}

// =========== TerminalPlayer ===========
// TerminalPlayer is a human player
type TerminalPlayer struct {
	runner *Runner
}

func NewTerminalPlayer(runner *Runner) *TerminalPlayer {
	return &TerminalPlayer{runner: runner}
}
func (t *TerminalPlayer) getMove() (*board.Move, bool) {
	return t.runner.getMoveTerminal()
}
func (t *TerminalPlayer) displayBoard(b *board.Board, player *board.Owner) {
	// print messages
	fmt.Printf("%v's turn:\n", *player)
	fmt.Println(b.TerminalString())
}
func (t *TerminalPlayer) afterMove(_ *board.Board, prevValid bool) {
	if !prevValid {
		fmt.Println("invalid move!!!")
	}
}

// =========== AIPlayer ===========
// represents an AI that communicates via protocol buffers
type AIPlayer struct{}

func NewAIPlayer() *AIPlayer {
	return &AIPlayer{}
}
func write(m protoreflect.ProtoMessage, filename string) {
	bytes, err := proto.Marshal(m)
	if err != nil {
		log.Fatalln("Failed to encode State Message")
	}

	// write the bytes
	if err := os.WriteFile(filename, bytes, 0644); err != nil {
		log.Fatalln("failed to write bytes")
	}
	fmt.Printf("sent message to %v\n", filename)
}
func (a *AIPlayer) displayBoard(b *board.Board, player *board.Owner) {
	state := board.StateMessage{Board: b, Turn: *player}
	write(&state, board.STATE_FILE)
}
func (a *AIPlayer) afterMove(b *board.Board, prevValid bool) {
	ret := board.ReturnMessage{Board: b, Valid: prevValid}
	write(&ret, board.RETURN_FILE)
}
func (a *AIPlayer) getMove() (*board.Move, bool) {
	// wait for content
	info, err := os.Stat(board.ACTION_FILE)
	for err != nil || info.Size() == 0 {
		info, err = os.Stat(board.ACTION_FILE)
	}

	// open content
	bytes, err := os.ReadFile(board.ACTION_FILE)
	if err != nil {
		log.Fatalln("failed to read in action")
	}

	message := &board.ActionMessage{}
	err = proto.Unmarshal(bytes, message)
	if err != nil {
		log.Fatalln("failed to decode action")
	}
	os.Truncate(board.ACTION_FILE, 0)

	return message.Move, true
}

// =======================================================
// =========== Terminal Helpers ===========
// =======================================================

func getCoord(where string) (c *board.Coord, quit bool) {
	fmt.Println("Where r u going (0 - 8)", where)
	var inp string
	fmt.Scanln(&inp)

	// process input
	quit = inp == "q"
	num, _ := strconv.ParseInt(inp, 10, 8)
	c = board.ToCoord(uint32(num))
	return
}

func (runner *Runner) getMoveTerminal() (move *board.Move, quit bool) {
	move = &board.Move{}
	if !runner.gameboard.CurCell.Valid() {
		move.Large, quit = getCoord("in large cells")
		if quit {
			return
		}
	} else {
		move.Large = &board.Coord{Row: runner.gameboard.CurCell.Row, Col: runner.gameboard.CurCell.Col}
	}

	move.Small, quit = getCoord("in small cells")
	return
}

// =======================================================
// =========== Run Section ===========
// =======================================================

func (runner *Runner) run(player1, player2 Player) {
	fmt.Println("playing Ultimate Tic-Tac-Toe")

	var curPlayer Player
	for runner.gameboard.Owner() == board.Owner_NONE {
		// get the turn number
		var playerNum board.Owner
		if runner.turn {
			playerNum = board.Owner_PLAYER1
			curPlayer = player1
		} else {
			playerNum = board.Owner_PLAYER2
			curPlayer = player2
		}

		curPlayer.displayBoard(runner.gameboard, &playerNum)

		move, quit := curPlayer.getMove()
		if quit {
			break
		}

		// validate move
		if validateMove(runner.gameboard, move) {
			runner.gameboard.Get(move.Large).(*board.Cell).Get(move.Small).(*board.Space).Val = playerNum

			// go to the next space
			if validateCell(runner.gameboard, move.Small) {
				runner.gameboard.CurCell.Row, runner.gameboard.CurCell.Col = move.Small.Row, move.Small.Col
			} else {
				runner.gameboard.CurCell.Invalidate()
			}

			// change turn
			runner.turn = !runner.turn
			curPlayer.afterMove(runner.gameboard, true)
		} else {
			curPlayer.afterMove(runner.gameboard, false)
		}
	}

	fmt.Println(runner.gameboard.TerminalString())
	fmt.Printf("%v won\n", runner.gameboard.Owner())
}

func (runner *Runner) RunPVP() {
	runner.run(NewTerminalPlayer(runner), NewTerminalPlayer(runner))
}
func (runner *Runner) RunPVAI() {
	runner.run(NewTerminalPlayer(runner), NewAIPlayer())
}
func (runner *Runner) RunAIs() {
	runner.run(NewAIPlayer(), NewAIPlayer())
}
