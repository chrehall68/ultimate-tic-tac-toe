package game

import (
	"fmt"
	"strconv"
)

type Runner struct {
	turn      bool
	gameboard board
}

func NewRunner() *Runner {
	return &Runner{turn: true, gameboard: newBoard()}
}

func getCoord(where string) (c coord, quit bool) {
	fmt.Println("Where r u going (0 - 8)", where)
	var inp string
	fmt.Scanln(&inp)

	// process input
	quit = inp == "q"
	num, _ := strconv.ParseInt(inp, 10, 8)
	c = coord{row: uint8(num) / 3, col: uint8(num) % 3}
	return c, quit
}

func (runner *Runner) Run() {
	fmt.Println("playing Ultimate Tic-Tac-Toe")
	for runner.gameboard.owner() == 0 {
		// get the turn number
		var playerNum uint8
		if runner.turn {
			playerNum = 1
		} else {
			playerNum = 2
		}

		// print messages
		fmt.Printf("Player %v's turn:\n", playerNum)
		fmt.Println(runner.gameboard.String())

		if !runner.gameboard.curCell.valid {
			c, q := getCoord("in large cells")
			if q {
				break
			}
			runner.gameboard.curCell.row, runner.gameboard.curCell.col = c.row, c.col
		}

		// take input
		innerCoord, q := getCoord("in small cells")
		if q {
			break
		}

		// validate move
		if validateMove(&runner.gameboard, runner.gameboard.curCell, innerCoord) {
			runner.gameboard.get(runner.gameboard.curCell).get(innerCoord).val = playerNum

			// go to the next space
			if validateCell(&runner.gameboard, innerCoord) {
				runner.gameboard.curCell.row, runner.gameboard.curCell.col = innerCoord.row, innerCoord.col
				runner.gameboard.curCell.valid = true
			} else {
				runner.gameboard.curCell.row, runner.gameboard.curCell.col = 255, 255
				runner.gameboard.curCell.valid = false
			}

			// change turn
			runner.turn = !runner.turn
		} else {
			fmt.Println("invalid move")
		}
	}
	fmt.Println(runner.gameboard.String())
}
