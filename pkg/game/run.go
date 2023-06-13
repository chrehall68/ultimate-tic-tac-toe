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

func getCoord(where string) (c board.Coord, quit bool) {
	fmt.Println("Where r u going (0 - 8)", where)
	var inp string
	fmt.Scanln(&inp)

	// process input
	quit = inp == "q"
	num, _ := strconv.ParseInt(inp, 10, 8)
	c = *board.ToCoord(uint32(num))
	return
}

func (runner *Runner) Run() {
	fmt.Println("playing Ultimate Tic-Tac-Toe")
	for runner.gameboard.Owner() == board.Owner_NONE {
		// get the turn number
		var playerNum board.Owner
		if runner.turn {
			playerNum = board.Owner_PLAYER1
		} else {
			playerNum = board.Owner_PLAYER2
		}

		// print messages
		fmt.Printf("%v's turn:\n", playerNum)
		fmt.Println(runner.gameboard.TerminalString())

		if !runner.gameboard.CurCell.Valid() {
			c, q := getCoord("in large cells")
			if q {
				break
			}
			runner.gameboard.CurCell.Row, runner.gameboard.CurCell.Col = c.Row, c.Col
		}

		// take input
		innerCoord, q := getCoord("in small cells")
		if q {
			break
		}

		// validate move
		if validateMove(runner.gameboard, runner.gameboard.CurCell, &innerCoord) {
			runner.gameboard.Get(runner.gameboard.CurCell).(*board.Cell).Get(&innerCoord).(*board.Space).Val = playerNum

			// go to the next space
			if validateCell(runner.gameboard, &innerCoord) {
				runner.gameboard.CurCell.Row, runner.gameboard.CurCell.Col = innerCoord.Row, innerCoord.Col
			} else {
				runner.gameboard.CurCell.Row, runner.gameboard.CurCell.Col = -1, -1
			}

			// change turn
			runner.turn = !runner.turn
		} else {
			fmt.Println("invalid move")
		}
	}
	fmt.Println(runner.gameboard.TerminalString())
}
