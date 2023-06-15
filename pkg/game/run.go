package game

import (
	"fmt"
	"strconv"
	"uttt/pkg/board"
)

type Runner struct {
	turn      bool
	gameboard *board.Board
}

func NewRunner() *Runner {
	return &Runner{turn: true, gameboard: board.NewProtoBoard()}
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

type moveFunc func() (*board.Move, bool)
type printStateFunc func(*board.Board, *board.Owner)
type returnStateFunc func(bool)

func (runner *Runner) run(m moveFunc, p printStateFunc, r returnStateFunc) {
	fmt.Println("playing Ultimate Tic-Tac-Toe")

	for runner.gameboard.Owner() == board.Owner_NONE {
		// get the turn number
		var playerNum board.Owner
		if runner.turn {
			playerNum = board.Owner_PLAYER1
		} else {
			playerNum = board.Owner_PLAYER2
		}

		p(runner.gameboard, &playerNum)

		move, quit := m()
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
			r(true)
		} else {
			r(false)
		}
	}

	fmt.Println(runner.gameboard.TerminalString())
	fmt.Printf("%v won\n", runner.gameboard.Owner())
}

func (runner *Runner) RunTerminal() {
	m := func() (*board.Move, bool) { return runner.getMoveTerminal() }

	p := func(_ *board.Board, player *board.Owner) {
		// print messages
		fmt.Printf("%v's turn:\n", *player)
		fmt.Println(runner.gameboard.TerminalString())
	}

	b := func(prevValid bool) {
		if !prevValid {
			fmt.Println("invalid move!!!")
		} else {
			fmt.Println("valid move")
		}
	}

	runner.run(m, p, b)
}
