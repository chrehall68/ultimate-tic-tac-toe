package game

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
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
type NetResources struct {
	stateConn, actionConn, returnConn net.Conn
}
type AIPlayer struct {
	player board.Owner
	nr     *NetResources
}

func NewNetResources() *NetResources {
	sListener, err := net.Listen("tcp", "localhost:"+board.STATE_PORT)
	if err != nil {
		log.Fatalln("failed to listen on state port")
	}
	aListener, err := net.Listen("tcp", "localhost:"+board.ACTION_PORT)
	if err != nil {
		log.Fatalln("failed to listen on action port")
	}
	rListener, err := net.Listen("tcp", "localhost:"+board.RETURN_PORT)
	if err != nil {
		log.Fatalln("failed to listen on return port")
	}

	sConn, err := sListener.Accept()
	if err != nil {
		log.Fatalln("failed to accept state connection")
	}
	aConn, err := aListener.Accept()
	if err != nil {
		log.Fatalln("failed to accept action connection")
	}
	rConn, err := rListener.Accept()
	if err != nil {
		log.Fatalln("failed to accept return connection")
	}

	return &NetResources{stateConn: sConn, actionConn: aConn, returnConn: rConn}
}
func NewAIPlayer(player_num board.Owner, nr *NetResources) *AIPlayer {
	return &AIPlayer{player: player_num, nr: nr}
}
func write(m protoreflect.ProtoMessage, con net.Conn) {
	bytes, err := proto.Marshal(m)
	if err != nil {
		log.Fatalln("Failed to encode State Message")
	}

	// write the bytes
	if _, err := con.Write(bytes); err != nil {
		log.Fatalln("failed to write bytes to socket")
	}
}
func (a *AIPlayer) getStateMessage(b *board.Board, player *board.Owner) *board.StateMessage {
	owners := make([]board.Owner, 9)
	for i := 0; i < board.CELLS; i++ {
		owners[i] = b.Get(board.ToCoord(uint32(i))).Owner()
	}
	winner := b.Owner()
	done := winner != board.Owner_NONE || b.Full()
	return &board.StateMessage{Board: b, Cellowners: owners, Turn: *player, Winner: winner, Done: done}
}

func (a *AIPlayer) displayBoard(b *board.Board, player *board.Owner) {
	write(a.getStateMessage(b, player), a.nr.stateConn)
}
func (a *AIPlayer) afterMove(b *board.Board, prevValid bool) {

	ret := board.ReturnMessage{State: a.getStateMessage(b, &a.player), Valid: prevValid}
	write(&ret, a.nr.returnConn)
}
func (a *AIPlayer) getMove() (*board.Move, bool) {
	// open content
	bytes := make([]byte, board.MAX_MSG_SIZE)
	n, err := a.nr.actionConn.Read(bytes)
	if err != nil {
		log.Fatalln("failed to read in action with error: ", err.Error())
	}

	message := &board.ActionMessage{}
	err = proto.Unmarshal(bytes[:n], message)
	if err != nil {
		log.Fatalln("failed to decode action")
	}

	// assume that the ai doesn't quit
	return message.Move, false
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
	//fmt.Println("playing Ultimate Tic-Tac-Toe")

	var curPlayer Player
	for runner.gameboard.Owner() == board.Owner_NONE && !runner.gameboard.Full() {
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

	// check if either player was a terminal player
	// if so, print out final message
	_, valid1 := player1.(*TerminalPlayer)
	_, valid2 := player2.(*TerminalPlayer)
	if valid1 || valid2 {
		fmt.Println(runner.gameboard.TerminalString())
		fmt.Printf("%v won\n", runner.gameboard.Owner())
	}
}

func (runner *Runner) RunPVP() {
	runner.run(NewTerminalPlayer(runner), NewTerminalPlayer(runner))
}
func (runner *Runner) RunPVAI() {
	nr := NewNetResources()
	runner.run(NewTerminalPlayer(runner), NewAIPlayer(board.Owner_PLAYER2, nr))

	time.Sleep(1 * time.Second)
	nr.actionConn.Close()
	nr.stateConn.Close()
	nr.returnConn.Close()
}
func (runner *Runner) RunAIVP() {
	nr := NewNetResources()
	runner.run(NewAIPlayer(board.Owner_PLAYER2, nr), NewTerminalPlayer(runner))

	time.Sleep(1 * time.Second)
	nr.actionConn.Close()
	nr.stateConn.Close()
	nr.returnConn.Close()
}
func (runner *Runner) RunAIs() {
	nr := NewNetResources()
	runner.run(NewAIPlayer(board.Owner_PLAYER1, nr), NewAIPlayer(board.Owner_PLAYER2, nr))

	time.Sleep(1 * time.Second)
	nr.actionConn.Close()
	nr.stateConn.Close()
	nr.returnConn.Close()
}
