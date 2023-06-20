package main

import (
	"fmt"
	"os"
	"uttt/pkg/game"
)

func main() {
	msg := "Please provide either `pvp` for Player vs Player, `pvai` for Player vs AI, or `aivai` for AI vs AI"
	if len(os.Args) > 1 {
		runner := game.NewRunner()

		mode := os.Args[1]
		switch mode {
		case "pvp":
			runner.RunPVP()
		case "pvai":
			runner.RunPVAI()
		case "aivai":
			runner.RunAIs()
		default:
			fmt.Println("That is not a valid option.")
			fmt.Println(msg)
		}

	} else {
		fmt.Println(msg)
	}
}
