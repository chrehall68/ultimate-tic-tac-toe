package game

import (
	"fmt"
	"strconv"
)

type Runner struct {
	turn      bool
	gameboard board
}

func NewRunner() Runner {
	runner := Runner{turn: false, gameboard: newBoard()}
	return runner
}

func (runner *Runner) Run() {
	fmt.Println("playing Ultimate Tic-Tac-Toe")
	for runner.gameboard.owner() == 0 {
		fmt.Println(runner.gameboard.String())

		// take input
		fmt.Println("Where r u going (0 - 8)")
		var inp string
		fmt.Scanln(&inp)

		// process input
		if inp == "q" {
			break
		}
		num, _ := strconv.ParseInt(inp, 10, 8)
		row := uint8(num) / 3
		col := uint8(num) % 3

		// get the turn number
		var playerNum uint8
		if runner.turn {
			playerNum = 1
		} else {
			playerNum = 2
		}

		// validate move
		if runner.gameboard.spaces[runner.gameboard.curCell.row][runner.gameboard.curCell.col].spaces[row][col].val == 0 {
			runner.gameboard.spaces[runner.gameboard.curCell.row][runner.gameboard.curCell.col].spaces[row][col].val = playerNum
			runner.gameboard.curCell.row, runner.gameboard.curCell.col = row, col
			runner.turn = !runner.turn
		} else {
			fmt.Println("invalid move")
		}
	}
	fmt.Println(runner.gameboard.String())
}
